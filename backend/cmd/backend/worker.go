package cmd

import (
	"fmt"
	"log"

	"github.com/BackAged/go-elasticsearch-react/backend/config"
	"github.com/BackAged/go-elasticsearch-react/backend/infra"
	"github.com/BackAged/go-elasticsearch-react/backend/repo"
	"github.com/BackAged/go-elasticsearch-react/backend/search"
	"github.com/BackAged/go-elasticsearch-react/backend/worker"
	"github.com/BackAged/steadyrabbit"
	"github.com/spf13/cobra"
)

var srvWrkr = &cobra.Command{
	Use:   "serve-worker",
	Short: "start a worker server",
	RunE:  serveWorker,
}

func serveWorker(cmd *cobra.Command, args []string) error {
	cnf := config.GetApp()
	wrkrCnf := config.GetWorker()
	fmt.Printf("loaded config => %+v\n", cnf)
	fmt.Printf("loaded worker config => %+v\n", wrkrCnf)

	log.Println("connecting elasticSearch")
	es, err := infra.NewEsClient(cnf.ElasticSearchURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected elasticSearch")

	prdRepo := repo.NewProductRepo(es, repo.RepoNameProduct)
	brndRepo := repo.NewBrandRepo(es, repo.RepoNameBrand)
	shpRepo := repo.NewShopRepo(es, repo.RepoNameShop)

	// initiating services
	svc := search.NewService(prdRepo, brndRepo, shpRepo)
	hndlr := worker.NewHandler(svc)

	srCnf := &steadyrabbit.Config{
		URL: wrkrCnf.Rabbit.URL,
		Consumer: &steadyrabbit.ConsumerConfig{
			AutoAck:          false,
			QosPrefetchCount: wrkrCnf.ConcurrencyCount,
			QueueConfig: &steadyrabbit.QueueConfig{
				QueueName:    wrkrCnf.Rabbit.QueueName,
				QueueDurable: true,
				QueueDeclare: true,
			},
			Bindings: []*steadyrabbit.BindingConfig{
				{
					Exchange: &steadyrabbit.ExchangeConfig{
						ExchangeDeclare: true,
						ExchangeName:    wrkrCnf.Rabbit.CatalogExchangeName,
						ExchangeType:    wrkrCnf.Rabbit.CatalogExchangeType,
						ExchangeDurable: true,
					},
					RoutingKeys: worker.RoutingKeys(),
				},
			},
		},
	}

	cnsmr, err := steadyrabbit.NewConsumer(srCnf)
	if err != nil {
		log.Println("error initializing consumer  ", err)
		return err
	}

	wrkr := worker.NewWorker(wrkrCnf.ConcurrencyCount, cnsmr)

	err = wrkr.SetUpRouter(hndlr)

	err = wrkr.Start()
	log.Println("worker shutdown with ", err)

	return err
}
