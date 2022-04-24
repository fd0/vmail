package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
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

			fmt.Println()

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

	t := newColoredTable()
	t.AddColumn(" Mailbox ", " {{ .Username }}@{{ .Domain }} ")
	t.AddColumn(" Quota ", " {{ if gt .Quota 0 }}{{ .Quota }}{{ end }} ")
	t.AddColumn(" Enabled ", " {{ .Enabled }} ")

	for _, a := range accounts {
		t.AddRow(a)
	}

	return t.Write(os.Stdout)
}

func newColoredTable() *table.Table {
	t := table.New()
	var hlen int
	t.PrintSeparator = func(wr io.Writer, s string) error {
		_, err := os.Stdout.WriteString(s + "\n")
		hlen = len(s)
		return err
	}

	highlight := color.New(color.Bold).PrintlnFunc()
	t.PrintHeader = func(wr io.Writer, s string) error {
		highlight(s)
		return nil
	}

	reverse := color.New(color.BgBlack).PrintfFunc()
	t.PrintData = func(wr io.Writer, i int, s string) error {
		if len(s) < hlen {
			// pad with spaces so that the lines with reverse colors looks nice
			s = s + strings.Repeat(" ", hlen-len(s))
		}
		var err error
		if i%2 == 0 {
			_, err = os.Stdout.WriteString(s + "\n")
		} else {
			reverse(s)
			fmt.Printf("\n")
		}

		return err
	}

	return t
}

func printAliases(db *DB, domain string) error {
	aliases, err := opts.db.FindAllAliases(domain)
	if err != nil {
		return err
	}

	if len(aliases) == 0 {
		return nil
	}

	aliasList := make(map[string][]Alias)
	for _, a := range aliases {
		name := a.SourceUsername.String
		if !a.SourceUsername.Valid {
			name = "*"
		}
		aliasList[name] = append(aliasList[name], a)
	}

	t := newColoredTable()
	t.AddColumn(" Alias ", " {{ .Alias }}@{{ .Domain }} ")
	t.AddColumn(" Destinations ", " {{ .Destinations }} ")
	t.AddColumn(" Blacklisted ", " {{ .Blacklisted }} ")
	t.AddColumn(" Enabled ", " {{ .Enabled }} ")

	type rowData struct {
		Alias        string
		Domain       string
		Destinations string
		Blacklisted  string
		Enabled      string
	}

	names := make([]string, 0, len(aliasList))
	for name := range aliasList {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		var destinations, blacklisted, enabled []string

		for _, a := range aliasList[name] {
			destinations = append(destinations, a.DestinationUsername+"@"+a.DestinationDomain)
			blacklisted = append(blacklisted, fmt.Sprintf("%v", a.Blacklisted))
			enabled = append(enabled, fmt.Sprintf("%v", a.Enabled))
		}

		t.AddRow(rowData{
			Alias:        name,
			Domain:       domain,
			Destinations: strings.Join(destinations, "\n "),
			Blacklisted:  strings.Join(blacklisted, "\n "),
			Enabled:      strings.Join(enabled, "\n "),
		})
	}

	return t.Write(os.Stdout)
}
