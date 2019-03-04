package main

import (
	"os"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/spf13/cobra"
)

var opts struct {
	Database string
}

func init() {
	root.Flags().StringVar(&opts.Database, "database", os.Getenv("VMAIL_DB"), "connect to this database")
}

var root = cobra.Command{}

func main() {
	err := root.Execute()
	if err != nil {
		panic(err)
	}
}
