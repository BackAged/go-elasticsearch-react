package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/BackAged/go-elasticsearch-react/backend/infra"
	"github.com/BackAged/go-elasticsearch-react/backend/search"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/mitchellh/mapstructure"
)

type productRepo struct {
	index  string
	client *elasticsearch.Client
}

// ProductRepo ...
type ProductRepo interface {
	search.ProductRepo
	Repo
}

// NewProductRepo returns a new productRepo
func NewProductRepo(client *elasticsearch.Client, index string) ProductRepo {
	return &productRepo{
		client: client,
		index:  index,
	}
}

func (pr *productRepo) EnsureMapping() error {
	log.Println("creating mappings for ", pr.index)
	res, err := pr.client.Indices.PutMapping(strings.NewReader(ProductMapping), pr.client.Indices.PutMapping.WithIndex(pr.index))

	if err != nil {
		log.Println(err)
		return err
	}

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("productRepo.EnsureIndexMapping: Error decoding error response: %s\n", err)
		} else {
			log.Printf("productRepo.EnsureIndexMapping: Error response [%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
		log.Println("couldn't create mappings for ", pr.index, "error: ", e)
		return err
	}

	return nil
}

func (pr *productRepo) EnsureIndex() error {
	res, err := pr.client.Indices.Get([]string{pr.index})
	if err != nil {
		log.Println(err)
		return err
	}

	if res.IsError() {
		log.Println("creating new index ", pr.index)
		indexCreateResponse, err := pr.client.Indices.Create(pr.index)
		if err != nil {
			log.Println(err)
			return err
		}
		if indexCreateResponse.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				log.Printf("productRepo.EnsureIndexMapping: Error decoding error response: %s\n", err)
			} else {
				log.Printf("productRepo.EnsureIndexMapping: Error response [%s] %s: %s",
					res.Status(),
					e["error"].(map[string]interface{})["type"],
					e["error"].(map[string]interface{})["reason"],
				)
			}
			log.Println("couldn't create index ", pr.index, "error: ", e)
			return err
		}
	}

	return nil
}

func (pr *productRepo) EnsureIndexAndMapping() error {
	if err := pr.EnsureIndex(); err != nil {
		return err
	}

	return pr.EnsureMapping()
}

func (pr *productRepo) BulkInsert(ctx context.Context, products []*search.Product) ([]*search.Product, error) {
	var buf bytes.Buffer
	for _, prd := range products {
		jPrd, _ := json.Marshal(prd)
		jPrd = append(jPrd, "\n"...)
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, prd.ShopItemID, "\n"))
		buf.Grow(len(meta) + len(jPrd))
		buf.Write(meta)
		buf.Write(jPrd)
	}

	req := esapi.BulkRequest{
		Index:   pr.index,
		Refresh: "true",
		Body:    bytes.NewReader(buf.Bytes()),
	}

	res, err := req.Do(context.Background(), pr.client)
	if err != nil {
		log.Printf("productRepo.BulkInsert: Error sending query: %s\n", err)
		return nil, err
	}
	defer res.Body.Close()

	fmt.Println(res.String())
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("productRepo.BulkInsert: Error decoding error response: %s\n", err)
		} else {
			log.Printf("productRepo.BulkInsert: Error response [%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return nil, infra.NewESResponseError("query returned error", e)
	}

	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("productRepo.BulkInsert: Error decoding success response: %s\n", err)
	} else {
		// Print the response status and indexed document version.
		log.Printf("[%s] %s; ", res.Status(), r["result"])
	}

	return nil, err
}

func (pr *productRepo) Add(ctx context.Context, product *search.Product) (*search.Product, error) {
	fmt.Println(product)
	fmt.Println(esutil.NewJSONReader(product))
	j, _ := json.Marshal(product)
	req := esapi.IndexRequest{
		Index:   pr.index,
		Body:    strings.NewReader(string(j)),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), pr.client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	fmt.Println(res.String())
	if res.IsError() {
		log.Printf("[%s] Error indexing document", res.Status())
		return nil, err
	}

	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("Error parsing the response body: %s", err)
	} else {
		// Print the response status and indexed document version.
		log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
	}

	return nil, err
}

func (pr *productRepo) Search(ctx context.Context, term string, skip int64, limit int64) ([]*search.Product, int64, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"from": skip,
		"size": limit,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": term,
				"type":  "bool_prefix",
				"fields": []string{
					"name",
					"shop_name",
					"category_name",
					"brand_name",
				},
				//"type": "most_fields",
			},
		},
		"sort": []map[string]interface{}{
			{
				"_score": map[string]interface{}{
					"order": "desc",
				},
				"pop_score": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("productRepo.Search: Error encoding query: %s\n", err)
		return nil, 0, err
	}

	// Perform the search request.
	res, err := pr.client.Search(
		pr.client.Search.WithContext(context.Background()),
		pr.client.Search.WithIndex(pr.index),
		pr.client.Search.WithBody(&buf),
		pr.client.Search.WithTrackTotalHits(true),
		pr.client.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("productRepo.Search: Error sending query to server: %s\n", err)
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("productRepo.Search: Error decoding error response: %s\n", err)
		} else {
			log.Printf("productRepo.Search: Error response [%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return nil, 0, infra.NewESResponseError("query returned error", e)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("productRepo.Search: Error decoding success response: %s\n", err)
		return nil, 0, err
	}

	total := int64(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	timeTaken := int64(r["took"].(float64))
	log.Printf(
		"[%s] %d hits; took: %dms\n",
		res.Status(),
		total,
		timeTaken,
	)

	prds := []*search.Product{}
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var b search.Product
		if err := mapstructure.Decode(hit.(map[string]interface{})["_source"], &b); err != nil {
			log.Printf("productRepo.Search: Error decoding response: %s\n", err)
			return nil, 0, err
		}
		prds = append(prds, &b)
	}

	return prds, total, nil

}

func (pr *productRepo) UpdateMany(ctx context.Context, products []*search.Product) ([]*search.Product, error) {
	var buf bytes.Buffer
	for _, prd := range products {
		jBrnd, _ := json.Marshal(prd)
		jiBrnd := string(jBrnd)
		meta := []byte(fmt.Sprintf(`{ "update" : { "_id" : "%d" } }%s`, prd.ShopItemID, "\n")) //"retry_on_conflict" : 5
		buf.Grow(len(meta) + len(jBrnd))
		buf.Write(meta)
		doc := getConcurrencyControlledUpdateQuery(jBrnd, jiBrnd)
		buf.Write(doc)
		fmt.Println(buf.String())
	}

	req := esapi.BulkRequest{
		Index:   pr.index,
		Refresh: "true",
		Body:    bytes.NewReader(buf.Bytes()),
	}

	res, err := req.Do(context.Background(), pr.client)
	if err != nil {
		log.Printf("productRepo.UpdateMany: Error sending query: %s\n", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("productRepo.UpdateMany: Error decoding error response: %s\n", err)
		} else {
			log.Printf("productRepo.UpdateMany: Error response [%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return nil, infra.NewESResponseError("query returned error", e)
	}

	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("productRepo.UpdateMany: Error decoding success response: %s\n", err)
	} else {
		log.Printf("[%s] %s; ", res.Status(), r["result"])
	}

	return nil, err
}

func (pr *productRepo) DeleteMany(ctx context.Context, shopItemIDS []int64) error {
	var buf bytes.Buffer
	for _, id := range shopItemIDS {
		meta := []byte(fmt.Sprintf(`{ "delete" : { "_id" : "%d" } }%s`, id, "\n"))
		buf.Grow(len(meta))
		buf.Write(meta)
	}

	req := esapi.BulkRequest{
		Index:   pr.index,
		Refresh: "true",
		Body:    bytes.NewReader(buf.Bytes()),
	}

	res, err := req.Do(context.Background(), pr.client)
	if err != nil {
		log.Printf("brandRepo.DeleteMany: Error sending query: %s\n", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("productRepo.DeleteMany: Error decoding error response: %s\n", err)
		} else {
			log.Printf("productRepo.DeleteMany: Error response [%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return infra.NewESResponseError("query returned error", e)
	}

	// Deserialize the response into a map.
	// var r map[string]interface{}
	// if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
	// 	log.Printf("productRepo.BulkInsert: Error decoding success response: %s\n", err)
	// } else {
	// 	log.Printf("[%s] %s; ", res.Status(), r["result"])
	// }

	return nil
}

func (pr *productRepo) UpdateProductScore(ctx context.Context, shopItemID int64) error {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"script": map[string]interface{}{
			"source": "if(ctx._source.pop_score==null){ctx._source.pop_score=0.0001} else {ctx._source.pop_score+=0.0001}",
			"lang":   "painless",
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("productRepo.UpdateProductScore: Error encoding query: %s\n", err)
		return err
	}

	// Perform the update by query request.
	res, err := pr.client.Update(
		pr.index, fmt.Sprintf("%d", shopItemID), &buf,
		pr.client.Update.WithContext(context.Background()),
		pr.client.Update.WithRetryOnConflict(5),
	)
	if err != nil {
		log.Printf("productRepo.UpdateProductScore: Error sending query: %s\n", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("productRepo.UpdateProductScore: Error decoding error response: %s\n", err)
		} else {
			log.Printf("productRepo.UpdateProductScore: Error response [%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return infra.NewESResponseError("query returned error", e)
	}

	return nil
}

func boolTermIndividual(name string, vals []string) map[string]interface{} {
	arg := []map[string]interface{}{}
	for _, f := range vals {
		s := map[string]interface{}{
			"term": map[string]interface{}{
				name: f,
			},
		}
		arg = append(arg, s)
	}

	s := map[string]interface{}{
		"bool": map[string]interface{}{
			"should": arg,
		},
	}

	return s
}

func applyPostFilter(query map[string]interface{}, req search.FacetSearchReq) map[string]interface{} {
	ss := []map[string]interface{}{}

	ss = append(ss, boolTermIndividual("category_name.keyword", req.CategoryFilters))
	ss = append(ss, boolTermIndividual("shop_name.keyword", req.ShopFilters))
	ss = append(ss, boolTermIndividual("brand_name.keyword", req.BrandFilters))
	ss = append(ss, boolTermIndividual("color.keyword", req.ColorFilters))

	query["post_filter"] = map[string]interface{}{
		"bool": map[string]interface{}{
			"filter": ss,
		},
	}

	return query
}

func applySort(query map[string]interface{}, req search.FacetSearchReq) map[string]interface{} {
	arg := []map[string]interface{}{}
	for _, sort := range req.Sort {
		s := map[string]interface{}{
			sort.FieldName: map[string]interface{}{
				"order": sort.Order,
			},
		}
		arg = append(arg, s)
	}
	arg = append(arg, map[string]interface{}{
		"_score": map[string]interface{}{
			"order": "desc",
		},
		"pop_score": map[string]interface{}{
			"order": "desc",
		},
	})

	query["sort"] = arg

	return query
}

func applyFacet(query map[string]interface{}, req search.FacetSearchReq) map[string]interface{} {
	ss := []map[string]interface{}{}

	ca := boolTermIndividual("category_name.keyword", req.CategoryFilters)
	s := boolTermIndividual("shop_name.keyword", req.ShopFilters)
	b := boolTermIndividual("brand_name.keyword", req.BrandFilters)
	c := boolTermIndividual("color.keyword", req.ColorFilters)

	cat := append(append(append(append(ss, s), b)), c)
	shp := append(append(append(append(ss, ca), b)), c)
	brnd := append(append(append(append(ss, ca), s)), c)
	clr := append(append(append(append(ss, ca), s)), b)

	query["aggs"] = map[string]interface{}{
		"categories": map[string]interface{}{
			"filter": map[string]interface{}{
				"bool": map[string]interface{}{
					"filter": cat,
				},
			},
			"aggs": map[string]interface{}{
				"categories_filtered": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": "category_name.keyword",
						"size":  req.BucketSize,
					},
				},
			},
		},
		"shops": map[string]interface{}{
			"filter": map[string]interface{}{
				"bool": map[string]interface{}{
					"filter": shp,
				},
			},
			"aggs": map[string]interface{}{
				"shops_filtered": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": "shop_name.keyword",
						"size":  req.BucketSize,
					},
				},
			},
		},
		"brands": map[string]interface{}{
			"filter": map[string]interface{}{
				"bool": map[string]interface{}{
					"filter": brnd,
				},
			},
			"aggs": map[string]interface{}{
				"brands_filtered": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": "brand_name.keyword",
						"size":  req.BucketSize,
					},
				},
			},
		},
		"colors": map[string]interface{}{
			"filter": map[string]interface{}{
				"bool": map[string]interface{}{
					"filter": clr,
				},
			},
			"aggs": map[string]interface{}{
				"color_filtered": map[string]interface{}{
					"terms": map[string]interface{}{
						"field":   "color.keyword",
						"size":    req.BucketSize,
						"exclude": []string{""},
					},
				},
			},
		},
	}

	return query
}

func (pr *productRepo) SearchFacet(ctx context.Context, req search.FacetSearchReq) ([]*search.Product, *search.FacetRes, int64, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"from": req.From,
		"size": req.Size,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": req.Term,
				"fields": []string{
					"name",
					"shop_name",
					"category_name",
					"brand_name",
				},
				"type": "most_fields",
				//"fuzziness": 1,
			},
		},
	}

	if req.Term == "" {
		query["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}

	applyFacet(query, req)
	applyPostFilter(query, req)
	applySort(query, req)

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("productRepo.SearchFacet: Error encoding query: %s\n", err)
		return nil, nil, 0, err
	}

	//println("=============================== ", buf.String())

	// Perform the search request.
	res, err := pr.client.Search(
		pr.client.Search.WithContext(context.Background()),
		pr.client.Search.WithIndex(pr.index),
		pr.client.Search.WithBody(&buf),
		pr.client.Search.WithTrackTotalHits(true),
		pr.client.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("productRepo.SearchFacet: Error sending query: %s\n", err)
		return nil, nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("productRepo.UpdateProductScore: Error decoding error response: %s\n", err)
		} else {
			log.Printf("productRepo.UpdateProductScore: Error response [%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return nil, nil, 0, infra.NewESResponseError("query returned error", e)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("productRepo.Search: Error decoding success response: %s\n", err)
		return nil, nil, 0, err
	}

	total := int64(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	timeTaken := int64(r["took"].(float64))
	log.Printf(
		"[%s] %d hits; took: %dms\n",
		res.Status(),
		total,
		timeTaken,
	)

	prds := []*search.Product{}
	fcts := search.FacetRes{}
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var b search.Product
		if err := mapstructure.Decode(hit.(map[string]interface{})["_source"], &b); err != nil {
			log.Printf("productRepo.Search: Error decoding response products: %s\n", err)
			return nil, nil, 0, err
		}
		prds = append(prds, &b)
	}

	brnds := r["aggregations"].(map[string]interface{})["brands"].(map[string]interface{})["brands_filtered"].(map[string]interface{})["buckets"].([]interface{})
	shps := r["aggregations"].(map[string]interface{})["shops"].(map[string]interface{})["shops_filtered"].(map[string]interface{})["buckets"].([]interface{})
	ctgrs := r["aggregations"].(map[string]interface{})["categories"].(map[string]interface{})["categories_filtered"].(map[string]interface{})["buckets"].([]interface{})
	clrs := r["aggregations"].(map[string]interface{})["colors"].(map[string]interface{})["color_filtered"].(map[string]interface{})["buckets"].([]interface{})

	brBckts := []search.Bucket{}
	for _, b := range brnds {
		var bct search.Bucket
		if err := mapstructure.Decode(b, &bct); err != nil {
			log.Printf("productRepo.Search: Error decoding response brand bucket: %s\n", err)
			return nil, nil, 0, err
		}
		brBckts = append(brBckts, bct)
	}

	clrBckts := []search.Bucket{}
	for _, b := range clrs {
		var bct search.Bucket
		if err := mapstructure.Decode(b, &bct); err != nil {
			log.Printf("productRepo.Search: Error decoding response color bucket: %s\n", err)
			return nil, nil, 0, err
		}
		clrBckts = append(clrBckts, bct)
	}

	ctrBckts := []search.Bucket{}
	for _, b := range ctgrs {
		var bct search.Bucket
		if err := mapstructure.Decode(b, &bct); err != nil {
			log.Printf("productRepo.Search: Error decoding response category bucket: %s\n", err)
			return nil, nil, 0, err
		}
		ctrBckts = append(ctrBckts, bct)
	}

	shpBckts := []search.Bucket{}
	for _, b := range shps {
		var bct search.Bucket
		if err := mapstructure.Decode(b, &bct); err != nil {
			log.Printf("productRepo.Search: Error decoding response shop bucket: %s\n", err)
			return nil, nil, 0, err
		}
		shpBckts = append(shpBckts, bct)
	}

	fcts.Brands = brBckts
	fcts.Categories = ctrBckts
	fcts.Colors = clrBckts
	fcts.Shops = shpBckts

	return prds, &fcts, total, nil
}
