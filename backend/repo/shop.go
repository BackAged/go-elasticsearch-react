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
	"github.com/mitchellh/mapstructure"
)

type shopRepo struct {
	index  string
	client *elasticsearch.Client
}

// ShopRepo ...
type ShopRepo interface {
	search.ShopRepo
	Repo
}

// NewShopRepo returns a new NewShopRepo
func NewShopRepo(client *elasticsearch.Client, index string) ShopRepo {
	return &shopRepo{
		client: client,
		index:  index,
	}
}

func (sr *shopRepo) EnsureMapping() error {
	log.Println("creating mappings for ", sr.index)
	res, err := sr.client.Indices.PutMapping(strings.NewReader(ShopMapping), sr.client.Indices.PutMapping.WithIndex(sr.index))

	if err != nil {
		log.Println(err)
		return err
	}

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("shopRepo.EnsureIndexMapping: Error decoding error response: %s\n", err)
		} else {
			log.Printf("shopRepo.EnsureIndexMapping: Error response [%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
		log.Println("couldn't create mappings for ", sr.index, "error: ", e)
		return err
	}

	return nil
}

func (sr *shopRepo) EnsureIndex() error {
	res, err := sr.client.Indices.Get([]string{sr.index})
	if err != nil {
		log.Println(err)
		return err
	}

	if res.IsError() {
		log.Println("creating new index ", sr.index)
		indexCreateResponse, err := sr.client.Indices.Create(sr.index)
		if err != nil {
			log.Println(err)
			return err
		}
		if indexCreateResponse.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				log.Printf("shopRepo.EnsureIndexMapping: Error decoding error response: %s\n", err)
			} else {
				log.Printf("shopRepo.EnsureIndexMapping: Error response [%s] %s: %s",
					res.Status(),
					e["error"].(map[string]interface{})["type"],
					e["error"].(map[string]interface{})["reason"],
				)
			}
			log.Println("couldn't create index ", sr.index, "error: ", err)
			return err
		}
	}

	return nil
}

func (sr *shopRepo) EnsureIndexAndMapping() error {
	if err := sr.EnsureIndex(); err != nil {
		return err
	}

	return sr.EnsureMapping()
}

func (sr *shopRepo) SearchAsType(ctx context.Context, term string, skip int64, limit int64) ([]*search.Shop, int64, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"from": skip,
		"size": limit,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": term,
				"type":  "bool_prefix",
				"fields": []string{
					"shop_name.search_as_type",
					"shop_name.search_as_type._2gram",
					"shop_name.search_as_type._3gram",
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("shopRepo.SearchAsType: Error encoding query: %s\n", err)
		return nil, 0, err
	}

	// Perform the search request.
	res, err := sr.client.Search(
		sr.client.Search.WithContext(context.Background()),
		sr.client.Search.WithIndex(sr.index),
		sr.client.Search.WithBody(&buf),
		sr.client.Search.WithTrackTotalHits(true),
		sr.client.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("shopRepo.SearchAsType: Error sending query to server: %s\n", err)
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("shopRepo.SearchAsType: Error decoding error response: %s\n", err)
		} else {
			log.Printf("shopRepo.SearchAsType: Error response [%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return nil, 0, infra.NewESResponseError("query returned error", e)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("shopRepo.SearchAsType: Error decoding success response: %s\n", err)
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

	shps := []*search.Shop{}
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var b search.Shop
		if err := mapstructure.Decode(hit.(map[string]interface{})["_source"], &b); err != nil {
			log.Printf("shopRepo.SearchAsType: Error decoding response: %s\n", err)
			return nil, 0, err
		}
		shps = append(shps, &b)
	}

	return shps, total, nil
}

func (sr *shopRepo) BulkInsert(ctx context.Context, shops []*search.Shop) ([]*search.Shop, error) {
	var buf bytes.Buffer
	for _, shp := range shops {
		jShp, _ := json.Marshal(shp)
		jShp = append(jShp, "\n"...)
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, shp.ID, "\n"))
		buf.Grow(len(meta) + len(jShp))
		buf.Write(meta)
		buf.Write(jShp)
	}

	req := esapi.BulkRequest{
		Index:   sr.index,
		Refresh: "true",
		Body:    bytes.NewReader(buf.Bytes()),
	}

	res, err := req.Do(context.Background(), sr.client)
	if err != nil {
		log.Printf("shopRepo.BulkInsert: Error sending query: %s\n", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("shopRepo.BulkInsert: Error decoding error response: %s\n", err)
		} else {
			log.Printf("shopRepo.BulkInsert: Error response [%s] %s: %s",
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
		log.Printf("shopRepo.BulkInsert: Error decoding success response: %s\n", err)
	} else {
		log.Printf("[%s] %s; ", res.Status(), r["result"])
	}

	return nil, err
}

func (sr *shopRepo) UpdateMany(ctx context.Context, shops []*search.Shop) ([]*search.Shop, error) {
	var buf bytes.Buffer
	for _, shp := range shops {
		jShp, _ := json.Marshal(shp)
		jiShp := string(jShp)
		meta := []byte(fmt.Sprintf(`{ "update" : { "_id" : "%d" } }%s`, shp.ID, "\n"))
		buf.Grow(len(meta) + len(jShp))
		buf.Write(meta)
		doc := getConcurrencyControlledUpdateQuery(jShp, jiShp)
		buf.Write(doc)
		fmt.Println(buf.String())
	}

	req := esapi.BulkRequest{
		Index:   sr.index,
		Refresh: "true",
		Body:    bytes.NewReader(buf.Bytes()),
	}

	res, err := req.Do(context.Background(), sr.client)
	if err != nil {
		log.Printf("shopRepo.UpdateMany: Error sending query: %s\n", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("shopRepo.UpdateMany: Error decoding error response: %s\n", err)
		} else {
			log.Printf("shopRepo.UpdateMany: Error response [%s] %s: %s",
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
		log.Printf("shopRepo.UpdateMany: Error decoding success response: %s\n", err)
	} else {
		log.Printf("[%s] %s; ", res.Status(), r["result"])
	}

	return nil, err
}

func (sr *shopRepo) DeleteMany(ctx context.Context, shopIDS []int64) error {
	var buf bytes.Buffer
	for _, id := range shopIDS {
		meta := []byte(fmt.Sprintf(`{ "delete" : { "_id" : "%d" } }%s`, id, "\n"))
		buf.Grow(len(meta))
		buf.Write(meta)
	}

	req := esapi.BulkRequest{
		Index:   sr.index,
		Refresh: "true",
		Body:    bytes.NewReader(buf.Bytes()),
	}

	res, err := req.Do(context.Background(), sr.client)
	if err != nil {
		log.Printf("shopRepo.DeleteMany: Error sending query: %s\n", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("shopRepo.DeleteMany: Error decoding error response: %s\n", err)
		} else {
			log.Printf("shopRepo.DeleteMany: Error response [%s] %s: %s",
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
	// 	log.Printf("shopRepo.BulkInsert: Error decoding success response: %s\n", err)
	// } else {
	// 	log.Printf("[%s] %s; ", res.Status(), r["result"])
	// }

	return nil
}
