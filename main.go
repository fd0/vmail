package main

import (
	"fmt"
	"os"
	"os/user"

	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/spf13/cobra"
)

var opts struct {
	Database string

	db *DB
}

func init() {
	var defaultDatabase = os.Getenv("VMAIL_DB")

	if defaultDatabase == "" {
		cur, err := user.Current()
		if err == nil {
			defaultDatabase = fmt.Sprintf("%s@unix(/run/mysqld/mysqld.sock)/vmail", cur.Username)
		}
	}

	root.Flags().StringVar(&opts.Database, "database", defaultDatabase, "connect to this database")
}

var root = cobra.Command{
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) (err error) {
		opts.db, err = ConnectDB("mysql", opts.Database)
		if err != nil {
			return err
		}
		return nil
	},
	PersistentPostRunE: func(_ *cobra.Command, _ []string) error {
		return opts.db.Close()
	},
}

func main() {
	err := root.Execute()
	if err != nil {
		warn("error: %v", err)
	}
}
