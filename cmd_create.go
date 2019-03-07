package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ncw/pwhash/sha512_crypt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func readPassword() (string, error) {
	fmt.Printf("enter password: ")
	buf, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	fmt.Printf("\n")
	fmt.Printf("repeat password: ")
	buf2, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Printf("\n")

	if !bytes.Equal(buf, buf2) {
		return "", errors.New("passwords do not match")
	}

	if len(buf) < 8 {
		return "", errors.New("password is way too short")
	}

	return string(buf), nil
}

const (
	hashRounds  = sha512_crypt.RoundsDefault * 10
	hashSaltLen = 20
)

func hashPassword(pw string) string {
	salt := sha512_crypt.GenerateSalt(hashSaltLen, hashRounds)
	return "{SHA512-CRYPT}" + sha512_crypt.Crypt(pw, salt)
}

func checkHash(hash string) error {
	if !strings.HasPrefix(hash, "{SHA512-CRYPT}$6$") {
		return errors.New("hash is invalid (does not start with '{SHA512-CRYPT}$6$')")
	}
	return nil
}

func splitMailAddress(s string) (string, string, error) {
	data := strings.SplitN(s, "@", -1)
	if len(data) != 2 {
		return "", "", fmt.Errorf("invalid email address %q", s)
	}

	return data[0], data[1], nil
}

func init() {
	var cmdCreate = &cobra.Command{
		Use:   "create",
		Short: "Create domains, accounts, and aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("the 'create' command needs to know what to create: domain, mailbox or alias?")
		},
	}

	cmdCreate.AddCommand(&cobra.Command{
		Use:   "domain [flags] name",
		Short: "Create a new domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("pass domain name to create as parameter")
			}

			name := args[0]
			err := opts.db.CreateDomain(name)
			if err != nil {
				return fmt.Errorf("creating domain %v failed: %v", name, err)
			}
			msg("domain %v created", name)
			return nil
		},
	})

	var createMailboxOpts = struct {
		Quota        uint64
		PasswordHash string
		Password     string
		SendOnly     bool
	}{}

	cmdCreateMailbox := &cobra.Command{
		Use:   "mailbox [flags] user@domain",
		Short: "Create a new mailbox",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) != 1 {
				return errors.New("pass mailbox to create as parameter (foo@example.com)")
			}

			mailbox := args[0]
			user, domain, err := splitMailAddress(mailbox)
			if err != nil {
				return err
			}

			msg("create %v with %+v", args[0], createMailboxOpts)

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

			err = opts.db.CreateAccount(Account{
				Domain:   domain,
				Username: user,
				Password: pwhash,
				Enabled:  true,
				Quota:    int(createMailboxOpts.Quota),
				Sendonly: createMailboxOpts.SendOnly,
			})

			if err != nil {
				return fmt.Errorf("creating mailbox %v failed: %v", mailbox, err)
			}
			msg("mailbox %v created", mailbox)
			return nil
		},
	}
	cmdCreateMailbox.Flags().Uint64Var(&createMailboxOpts.Quota, "quota", 0, "grant this mailbox `bytes` ")
	cmdCreateMailbox.Flags().StringVar(&createMailboxOpts.Password, "password", "", "use `pwd` as the password")
	cmdCreateMailbox.Flags().StringVar(&createMailboxOpts.PasswordHash, "password-hash", "", "use `hash` as the password (already hashed)")
	cmdCreateMailbox.Flags().BoolVar(&createMailboxOpts.SendOnly, "send-only", false, "do not receive mail for this account")

	cmdCreate.AddCommand(cmdCreateMailbox)

	root.AddCommand(cmdCreate)
}
