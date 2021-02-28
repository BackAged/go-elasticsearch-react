package cmd

import (
	"github.com/BackAged/go-elasticsearch-react/backend/cmd/backend/migration"
	"github.com/spf13/cobra"
)

var mgrtnCmd = &cobra.Command{
	Use:   "migration",
	Short: "migrates database schemas",
}

func init() {
	mgrtnCmd.AddCommand(migration.MgrtnUP, migration.MgrtnDOWN)
}
