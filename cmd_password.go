package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var passwordOptions = struct {
	PasswordHash string
	Password     string
}{}

func init() {
	cmdPassword.Flags().StringVar(&createMailboxOpts.Password, "password", "", "use `pwd` as the password")
	cmdPassword.Flags().StringVar(&createMailboxOpts.PasswordHash, "password-hash", "", "use `hash` as the password (already hashed)")
	root.AddCommand(cmdPassword)
}

var cmdPassword = &cobra.Command{
	Use:   "password [flags] user@domain",
	Short: "Reset the password of a mailbox",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 1 {
			return errors.New("pass target mailbox as parameter (foo@example.com)")
		}

		mailbox := args[0]
		user, domain, err := splitMailAddress(mailbox)
		if err != nil {
			return err
		}

		var pwhash = createMailboxOpts.PasswordHash
		if pwhash == "" {
			var pw = createMailboxOpts.Password
			if pw == "" && createMailboxOpts.PasswordHash == "" {
				pw, err = readPassword()
				if err != nil {
					return err
				}
			}

			pwhash = hashPassword(pw)
		}

		err = checkHash(pwhash)
		if err != nil {
			return err
		}

		err = opts.db.UpdateAccountPassword(user, domain, pwhash)
		if err != nil {
			return fmt.Errorf("updating password for %v failed: %v", mailbox, err)
		}
		msg("password for %v updated", mailbox)
		return nil
	},
}
