package config

import (
	"fmt"
	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var appOnce = sync.Once{}
var taskOnce = sync.Once{}
var workerOnce = sync.Once{}

// Worker holds worker config
type Worker struct {
	Name             string `yaml:"name"`
	ConcurrencyCount int    `yaml:"concurrency_count"`
	Rabbit           Rabbit
}

// Rabbit ...
type Rabbit struct {
	URL                 string `yaml:"rabbit_url"`
	QueueName           string `yaml:"queue_name"`
	CatalogExchangeName string `yaml:"catalog_exchange_name"`
	CatalogExchangeType string `yaml:"catalog_exchange_type"`
}

// Application holds application configurations
type Application struct {
	Host             string `yaml:"host"`
	Port             int    `yaml:"port"`
	ElasticSearchURL string `yaml:"elasticsearch_url"`
	GracefulTimeout  int    `yaml:"graceful_timeout"`
	APIKey           string `yaml:"api_key"`
}

// AMQP defines amqp config
type AMQP struct {
	Exchange      string
	ExchangeType  string
	BindingKey    string
	PrefetchCount int
}

// TaskConfig deines task worker and server realted config
type TaskConfig struct {
	Broker          string
	DefaultQueue    string
	ResultBackend   string
	AMQP            AMQP
	Worker          Worker
	ResultsExpireIn int
	Rabbit          Rabbit
}

var appConfig *Application
var taskconfig *TaskConfig
var workerConfig *Worker

// loadApp loads config from path
func loadApp() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not found, that's okay!")
	}

	viper.AutomaticEnv()

	appConfig = &Application{
		Host:             viper.GetString("HOST"),
		GracefulTimeout:  viper.GetInt("GRACEFUL_TIME_OUT"),
		Port:             viper.GetInt("PORT"),
		ElasticSearchURL: viper.GetString("ELASTICSEARCH_URL"),
		APIKey:           viper.GetString("X_API_KEY"),
	}

	return nil
}

// loadTaskConfig loads task config
func loadTaskConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not found, that's okay!")
	}

	viper.AutomaticEnv()

	taskconfig = &TaskConfig{
		Broker:        viper.GetString("RABBIT_URL"),
		DefaultQueue:  "search_product_update_score",
		ResultBackend: viper.GetString("TASKRESULT_BACKEND"),
		AMQP: AMQP{
			Exchange:      "search",
			ExchangeType:  "direct",
			BindingKey:    "search_product_update_score",
			PrefetchCount: 1,
		},
		Worker: Worker{
			Name:             viper.GetString("WORKER_NAME"),
			ConcurrencyCount: viper.GetInt("WORKER_CONCURRENCY_LIMIT"),
		},
		ResultsExpireIn: viper.GetInt("RESULT_EXPIRES_IN"),
	}
}

func defaultifyTask() {

}

func loadWorker() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env not found, that's okay!")
	}

	viper.AutomaticEnv()

	workerConfig = &Worker{
		Name:             viper.GetString("WORKER_NAME"),
		ConcurrencyCount: viper.GetInt("WORKER_CONCURRENCY_LIMIT"),
		Rabbit: Rabbit{
			URL:                 viper.GetString("RABBIT_URL"),
			QueueName:           viper.GetString("QUEUE_NAME"),
			CatalogExchangeName: viper.GetString("CATALOG_EXCHANGE_NAME"),
			CatalogExchangeType: viper.GetString("CATALOG_EXCHANGE_TYPE"),
		},
	}
}

// GetApp returns application config
func GetApp() *Application {
	appOnce.Do(func() {
		loadApp()
	})

	return appConfig
}

// GetWorker returns worker config
func GetWorker() *Worker {
	workerOnce.Do(func() {
		loadWorker()
	})

	return workerConfig
}

// GetTask returns taskconfig
func GetTask() *TaskConfig {
	taskOnce.Do(func() {
		loadTaskConfig()
	})

	defaultifyTask()

	return taskconfig
}
