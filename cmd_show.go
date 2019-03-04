package main

import (
	"errors"
	"os"
	"strings"

	"github.com/fd0/vmail/table"
	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(&cobra.Command{
		Use:   "show [options] [name]",
		Short: "Display a domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("pass domain to show as parameter")
			}

			name := args[0]

			_, err := opts.db.FindDomain(name)
			if err != nil {
				return err
			}

			err = printAccounts(opts.db, name)
			if err != nil {
				return err
			}

			err = printAliases(opts.db, name)
			if err != nil {
				return err
			}

			return nil
		},
	})
}

func printAccounts(db *DB, name string) error {
	accounts, err := opts.db.FindAllAccounts(name)
	if err != nil {
		return err
	}

	if len(accounts) == 0 {
		return nil
	}

	t := table.New()
	t.AddColumn(" Name ", " {{ .Username }}@{{ .Domain }} ")
	t.AddColumn(" Quota ", " {{ .Quota }} ")
	t.AddColumn(" Enabled ", " {{ .Enabled }} ")

	msg("\nMailboxes:\n\n")

	for _, a := range accounts {
		t.AddRow(a)
	}

	return t.Write(os.Stdout)
}

func printAliases(db *DB, name string) error {
	aliases, err := opts.db.FindAllAliases(name)
	if err != nil {
		return err
	}

	if len(aliases) == 0 {
		return nil
	}

	msg("\nAliases:\n\n")

	aliasList := make(map[string][]Alias)
	for _, a := range aliases {
		name := a.SourceUsername
		aliasList[name] = append(aliasList[name], a)
	}

	t := table.New()
	t.AddColumn(" Name ", " {{ .Name }} ")
	t.AddColumn(" Destinations ", " {{ .Destinations }} ")

	for name, aliases := range aliasList {
		var destinations []string

		for _, a := range aliases {
			destinations = append(destinations, a.DestinationUsername+"@"+a.DestinationDomain)
		}

		data := struct {
			Name         string
			Destinations string
		}{
			Name:         name,
			Destinations: strings.Join(destinations, "\n "),
		}

		t.AddRow(data)
	}

	return t.Write(os.Stdout)
}
