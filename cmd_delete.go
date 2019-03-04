package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	var cmdNew = &cobra.Command{
		Use: "delete",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("the 'delete' command needs to know what to create: domain, mailbox or alias?")
		},
	}

	cmdNew.AddCommand(&cobra.Command{
		Use:   "domain [options] name",
		Short: "Delete a domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("pass a domain name as parameter")
			}

			name := args[0]
			err := opts.db.DeleteDomain(name)
			if err != nil {
				warn("error deleting domain %v: %v", name, err)
				return nil
			}
			msg("domain %v deleted", name)
			return nil
		},
	})

	root.AddCommand(cmdNew)
}
