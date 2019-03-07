package main

import (
	"github.com/spf13/cobra"
)

func init() {
	root.AddCommand(&cobra.Command{
		Use:   "domains [flags] [filter]",
		Short: "List domains (with optional filter)",
		RunE: func(cmd *cobra.Command, args []string) error {
			var name string
			if len(args) > 0 {
				name = args[0]
			}

			domains, err := opts.db.FindAllDomains(name)
			if err != nil {
				return err
			}

			for _, d := range domains {
				msg("%v\n", d.Domain)
			}

			return nil
		},
	})
}
