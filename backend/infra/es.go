package infra

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
)

// NewEsClient returns es clients
func NewEsClient(esHost string) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:     []string{esHost},
		RetryOnStatus: []int{502, 503, 504, 429},
		MaxRetries:    5,
	})
	if err != nil {
		log.Fatalf("Error connecting elasticSearch: %s", err)
	}

	return es, err
}

// ESResponseError ...
type ESResponseError struct {
	Message string
	Details map[string]interface{}
}

func (e *ESResponseError) Error() string {
	return fmt.Sprintf(e.Message)
}

// NewESResponseError ...
func NewESResponseError(message string, details map[string]interface{}) *ESResponseError {
	return &ESResponseError{
		Message: message,
		Details: details,
	}
}
