package repo

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
)

// Repo name const
const (
	RepoNameBrand   = "brands"
	RepoNameShop    = "shops"
	RepoNameProduct = "products"
)

// Repo defines base repo interface
type Repo interface {
	EnsureIndexAndMapping() error
}

// MigrateUp migrates up
func MigrateUp(es *elasticsearch.Client) error {
	prdRepo := NewProductRepo(es, RepoNameProduct)
	brndRepo := NewBrandRepo(es, RepoNameBrand)
	shpRepo := NewShopRepo(es, RepoNameShop)

	repos := []Repo{prdRepo, brndRepo, shpRepo}

	for _, v := range repos {
		if err := v.EnsureIndexAndMapping(); err != nil {
			return err
		}
	}

	return nil
}

// MigrateDown migrates up
func MigrateDown(es *elasticsearch.Client) error {
	_, err := es.Indices.Delete([]string{
		RepoNameBrand, RepoNameProduct, RepoNameShop,
	})
	if err != nil {
		log.Println("error deleting indices, error: ", err)
	}

	return err
}

func getConcurrencyControlledUpdateQuery(dataByte []byte, inlineJSONData string) []byte {
	source, err := getUpdateSourceString(dataByte)
	if err != nil {
		log.Println("getConcurrencyControlledUpdateQuery error:", err)
	}

	return []byte(fmt.Sprintf(`{"script": {"source": %s,"lang": "painless","params": %s},"upsert": %s}%s`, source, inlineJSONData, inlineJSONData, "\n"))
}

func getUpdateSourceString(buf []byte) (string, error) {
	var dataMap map[string]interface{}
	err := json.Unmarshal(buf, &dataMap)
	if err != nil {
		log.Println(err)
		return "", err
	}

	source := ""
	i := 0
	for k := range dataMap {
		if i != 0 {
			source += "; "
		}
		source += fmt.Sprintf("ctx._source.%s=params.%s", k, k)
		i++
	}
	if source != "" {
		source = fmt.Sprintf(`"if (ctx._source.version < params.version) {%s}"`, source)
	}

	return source, nil
}
