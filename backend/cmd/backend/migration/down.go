package migration

import (
	"fmt"
	"log"

	"github.com/BackAged/go-elasticsearch-react/backend/config"
	"github.com/BackAged/go-elasticsearch-react/backend/infra"
	"github.com/BackAged/go-elasticsearch-react/backend/repo"
	"github.com/spf13/cobra"
)

// MgrtnDOWN migrates up the db
var MgrtnDOWN = &cobra.Command{
	Use:   "down",
	Short: "migrates down database schema",
	RunE:  down,
}

func down(cmd *cobra.Command, args []string) error {
	cnf := config.GetApp()
	fmt.Printf("loaded config => %+v\n", cnf)

	log.Println("connecting elasticSearch")
	es, err := infra.NewEsClient(cnf.ElasticSearchURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected elasticSearch")

	return repo.MigrateDown(es)
}
