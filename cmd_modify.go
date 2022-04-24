package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var modifyOpts = struct {
	Enable    bool
	Disable   bool
	Blacklist bool
}{}

func init() {
	cmdModify.Flags().BoolVar(&modifyOpts.Enable, "enable", false, "Enable alias")
	cmdModify.Flags().BoolVar(&modifyOpts.Disable, "disable", false, "Disable alias")
	cmdModify.Flags().BoolVar(&modifyOpts.Blacklist, "blacklist", false, "Mark alias as blacklist")

	root.AddCommand(cmdModify)
}

var cmdModify = &cobra.Command{
	Use:   "modify",
	Short: "Modify accounts and aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("the 'modify' command needs exactly one alias")
		}

		localPart, domain, err := splitMailAddress(args[0])
		if err != nil {
			return fmt.Errorf("split domain: %w", err)
		}

		aliases, err := opts.db.FindAliases(localPart, domain)
		if err != nil {
			return fmt.Errorf("find alias: %w", err)
		}

		if len(aliases) == 0 {
			return errors.New("no aliases found")
		}

		for _, alias := range aliases {
			if modifyOpts.Blacklist {
				alias.Blacklisted = true
			}

			if modifyOpts.Enable {
				alias.Enabled = true
			} else if modifyOpts.Disable {
				alias.Enabled = false
			}

			err := opts.db.UpdateAlias(alias)
			if err != nil {
				return fmt.Errorf("error updating alias %v: %w", alias.ID, err)
			}
		}

		fmt.Printf("successfully updated %d aliases\n", len(aliases))

		return nil
	},
}
