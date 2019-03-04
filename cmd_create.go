package main

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	var cmdCreate = &cobra.Command{
		Use: "create",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("the 'create' command needs to know what to create: domain, mailbox or alias?")
		},
	}

	cmdCreate.AddCommand(&cobra.Command{
		Use:   "domain [options] name",
		Short: "Create create domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("pass create domain name as parameter")
			}

			name := args[0]
			err := opts.db.CreateDomain(name)
			if err != nil {
				warn("error creating domain %v: %v", name, err)
			}
			msg("domain %v created", name)
			return nil
		},
	})

	root.AddCommand(cmdCreate)
}
