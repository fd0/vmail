package main

import (
	"errors"

	"github.com/spf13/cobra"
)

var cmdDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete domains, accounts, and aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("the 'delete' command needs to know what to create: domain, mailbox or alias?")
	},
}

var cmdDeleteDomain = &cobra.Command{
	Use:   "domain [flags] name",
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
}

var cmdDeleteMailbox = &cobra.Command{
	Use:   "mailbox [flags] name",
	Short: "Delete a mailbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("pass a mailbox as parameter (foo@example.com)")
		}

		mailbox := args[0]
		user, domain, err := splitMailAddress(mailbox)
		if err != nil {
			return err
		}

		err = opts.db.DeleteMailbox(user, domain)
		if err != nil {
			warn("error deleting mailbox %v: %v", mailbox, err)
			return nil
		}
		msg("mailbox %v deleted", mailbox)
		return nil
	},
}

func init() {
	cmdDelete.AddCommand(cmdDeleteDomain)
	cmdDelete.AddCommand(cmdDeleteMailbox)
	root.AddCommand(cmdDelete)
}
