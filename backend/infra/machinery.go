package infra

import (
	"fmt"

	configCust "github.com/BackAged/go-elasticsearch-react/backend/config"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

// NewMachineryServer return machinery server
func NewMachineryServer() (*machinery.Server, error) {
	cfg := configCust.GetTask()
	fmt.Println(cfg)
	tskSrvr, err := machinery.NewServer(&config.Config{
		Broker:        cfg.Broker,
		DefaultQueue:  cfg.DefaultQueue,
		ResultBackend: cfg.ResultBackend,
		AMQP: &config.AMQPConfig{
			ExchangeType:  cfg.AMQP.ExchangeType,
			Exchange:      cfg.AMQP.Exchange,
			BindingKey:    cfg.AMQP.BindingKey,
			PrefetchCount: cfg.AMQP.PrefetchCount,
		},
		ResultsExpireIn: cfg.ResultsExpireIn,
	})

	fmt.Println(cfg)

	return tskSrvr, err
}
