package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/amacneil/dbmate/pkg/dbmate"
	_ "github.com/amacneil/dbmate/pkg/driver/postgres"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"

	"github.com/goforbroke1006/onix/common"
)

func NewUtilMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use: "migrate",
		Run: func(cmd *cobra.Command, args []string) {
			connString := common.GetDBConnString()
			dbUrl, err := url.Parse(connString)
			if err != nil {
				fmt.Println("ERROR:", err.Error())
				os.Exit(1)
			}

			db := dbmate.New(dbUrl)
			db.MigrationsDir = "/db/migrations"

			err = db.CreateAndMigrate()
			if err != nil {
				fmt.Println("ERROR:", err.Error())
				os.Exit(1)
			}

			fmt.Println("ok")
		},
	}
}
