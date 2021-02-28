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

type brandRepo struct {
	index  string
	client *elasticsearch.Client
}

// BrandRepo ...
type BrandRepo interface {
	search.BrandRepo
	Repo
}

// NewBrandRepo returns a new NewBrandRepo
func NewBrandRepo(client *elasticsearch.Client, index string) BrandRepo {
	return &brandRepo{
		client: client,
		index:  index,
	}
}

func (br *brandRepo) EnsureMapping() error {
	log.Println("creating mappings for ", br.index)
	res, err := br.client.Indices.PutMapping(strings.NewReader(BrandMapping), br.client.Indices.PutMapping.WithIndex(br.index))

	if err != nil {
		log.Println(err)
		return err
	}

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("brandRepo.EnsureIndexMapping: Error decoding error response: %s\n", err)
		} else {
			log.Printf("brandRepo.EnsureIndexMapping: Error response [%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
		log.Println("couldn't create mappings for ", br.index, "error: ", e)
		return err
	}

	return nil
}

func (br *brandRepo) EnsureIndex() error {
	res, err := br.client.Indices.Get([]string{br.index})
	if err != nil {
		log.Println(err)
		return err
	}

	if res.IsError() {
		log.Println("creating new index ", br.index)
		indexCreateResponse, err := br.client.Indices.Create(br.index)
		if err != nil {
			log.Println(err)
			return err
		}
		if indexCreateResponse.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				log.Printf("brandRepo.EnsureIndexMapping: Error decoding error response: %s\n", err)
			} else {
				log.Printf("brandRepo.EnsureIndexMapping: Error response [%s] %s: %s",
					res.Status(),
					e["error"].(map[string]interface{})["type"],
					e["error"].(map[string]interface{})["reason"],
				)
			}
			log.Println("couldn't create index ", br.index, "error: ", e)
			return err
		}
	}

	return nil
}

func (br *brandRepo) EnsureIndexAndMapping() error {
	if err := br.EnsureIndex(); err != nil {
		return err
	}

	return br.EnsureMapping()
}

func (br *brandRepo) SearchAsType(ctx context.Context, term string, skip int64, limit int64) ([]*search.Brand, int64, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"from": skip,
		"size": limit,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": term,
				"type":  "bool_prefix",
				"fields": []string{
					"name.search_as_type",
					"name.search_as_type._2gram",
					"name.search_as_type._3gram",
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("brandRepo.SearchAsType: Error encoding query: %s\n", err)
		return nil, 0, err
	}

	// Perform the search request.
	res, err := br.client.Search(
		br.client.Search.WithContext(context.Background()),
		br.client.Search.WithIndex(br.index),
		br.client.Search.WithBody(&buf),
		br.client.Search.WithTrackTotalHits(true),
		br.client.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("brandRepo.SearchAsType: Error sending query to server: %s\n", err)
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("brandRepo.SearchAsType: Error decoding error response: %s\n", err)
		} else {
			log.Printf("brandRepo.SearchAsType: Error response [%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		return nil, 0, infra.NewESResponseError("query returned error", e)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("brandRepo.SearchAsType: Error decoding success response: %s\n", err)
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

	brnds := []*search.Brand{}
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var b search.Brand
		if err := mapstructure.Decode(hit.(map[string]interface{})["_source"], &b); err != nil {
			log.Printf("brandRepo.SearchAsType: Error decoding response: %s\n", err)
			return nil, 0, err
		}
		brnds = append(brnds, &b)
	}

	return brnds, total, nil

}

func (br *brandRepo) BulkInsert(ctx context.Context, brands []*search.Brand) ([]*search.Brand, error) {
	var buf bytes.Buffer
	for _, brnd := range brands {
		jBrnd, _ := json.Marshal(brnd)
		jBrnd = append(jBrnd, "\n"...)
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, brnd.ID, "\n"))
		buf.Grow(len(meta) + len(jBrnd))
		buf.Write(meta)
		buf.Write(jBrnd)
	}

	req := esapi.BulkRequest{
		Index:   br.index,
		Refresh: "true",
		Body:    bytes.NewReader(buf.Bytes()),
	}

	res, err := req.Do(context.Background(), br.client)
	if err != nil {
		log.Printf("brandRepo.BulkInsert: Error sending query: %s\n", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("brandRepo.BulkInsert: Error decoding error response: %s\n", err)
		} else {
			log.Printf("brandRepo.BulkInsert: Error response [%s] %s: %s",
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
		log.Printf("brandRepo.BulkInsert: Error decoding success response: %s\n", err)
	} else {
		log.Printf("[%s] %s; ", res.Status(), r["result"])
	}

	return nil, err
}

func (br *brandRepo) UpdateMany(ctx context.Context, brands []*search.Brand) ([]*search.Brand, error) {
	var buf bytes.Buffer
	for _, brnd := range brands {
		jBrnd, _ := json.Marshal(brnd)
		jiBrnd := string(jBrnd)
		meta := []byte(fmt.Sprintf(`{ "update" : { "_id" : "%d" } }%s`, brnd.ID, "\n"))
		buf.Grow(len(meta) + len(jBrnd))
		buf.Write(meta)
		doc := getConcurrencyControlledUpdateQuery(jBrnd, jiBrnd)
		buf.Write(doc)
		fmt.Println(buf.String())
	}

	req := esapi.BulkRequest{
		Index:   br.index,
		Refresh: "true",
		Body:    bytes.NewReader(buf.Bytes()),
	}

	res, err := req.Do(context.Background(), br.client)
	if err != nil {
		log.Printf("brandRepo.UpdateMany: Error sending query: %s\n", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("brandRepo.UpdateMany: Error decoding error response: %s\n", err)
		} else {
			log.Printf("brandRepo.UpdateMany: Error response [%s] %s: %s",
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
		log.Printf("brandRepo.UpdateMany: Error decoding success response: %s\n", err)
	} else {
		log.Printf("[%s] %s; ", res.Status(), r["result"])
	}

	return nil, err
}

func (br *brandRepo) DeleteMany(ctx context.Context, brandIDS []int64) error {
	var buf bytes.Buffer
	for _, id := range brandIDS {
		meta := []byte(fmt.Sprintf(`{ "delete" : { "_id" : "%d" } }%s`, id, "\n"))
		buf.Grow(len(meta))
		buf.Write(meta)
	}

	req := esapi.BulkRequest{
		Index:   br.index,
		Refresh: "true",
		Body:    bytes.NewReader(buf.Bytes()),
	}

	res, err := req.Do(context.Background(), br.client)
	if err != nil {
		log.Printf("brandRepo.DeleteMany: Error sending query: %s\n", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("brandRepo.DeleteMany: Error decoding error response: %s\n", err)
		} else {
			log.Printf("brandRepo.DeleteMany: Error response [%s] %s: %s",
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
	// 	log.Printf("brandRepo.BulkInsert: Error decoding success response: %s\n", err)
	// } else {
	// 	log.Printf("[%s] %s; ", res.Status(), r["result"])
	// }

	return nil
}
