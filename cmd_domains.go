package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(&cobra.Command{
		Short: "domains",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := ConnectDB("mysql", opts.Database)
			if err != nil {
				return err
			}

			domains, err := db.FindAllDomains()
			if err != nil {
				return err
			}

			for _, d := range domains {
				fmt.Printf("%v\n", d)
			}

			return db.Close()
		},
	})
}
