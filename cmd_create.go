package main

import (
	"bytes"
	"database/sql"
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

	user, domain := data[0], data[1]
	if user == "" {
		return "", "", fmt.Errorf("invalid email address (user part is empty)")
	}
	if domain == "" {
		return "", "", fmt.Errorf("invalid email address (domain part is empty)")
	}

	return user, domain, nil
}

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create domains, accounts, and aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("the 'create' command needs to know what to create: domain, mailbox or alias?")
	},
}

var cmdCreateDomain = &cobra.Command{
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
}

var createMailboxOpts = struct {
	Quota           uint64
	PasswordHash    string
	Password        string
	RawPasswordHash bool
	SendOnly        bool
}{}

func init() {
	cmdCreateMailbox.Flags().Uint64Var(&createMailboxOpts.Quota, "quota", 0, "grant this mailbox `bytes` ")
	cmdCreateMailbox.Flags().StringVar(&createMailboxOpts.Password, "password", "", "use `pwd` as the password")
	cmdCreateMailbox.Flags().StringVar(&createMailboxOpts.PasswordHash, "password-hash", "", "use `hash` as the password (already hashed)")
	cmdCreateMailbox.Flags().BoolVar(&createMailboxOpts.RawPasswordHash, "raw-password-hash", false, "do not check password hash")
	cmdCreateMailbox.Flags().BoolVar(&createMailboxOpts.SendOnly, "send-only", false, "do not receive mail for this account")
}

var cmdCreateMailbox = &cobra.Command{
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

		fmt.Printf("options: %+v\n", createMailboxOpts.RawPasswordHash)

		if !createMailboxOpts.RawPasswordHash {
			fmt.Printf("check hash\n")
			err = checkHash(pwhash)
			if err != nil {
				return err
			}
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

var cmdCreateAlias = &cobra.Command{
	Use:   "alias [flags] SRC DEST [DEST] [DEST...]",
	Short: "Create a new alias",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) < 2 {
			return errors.New("pass source and destinations")
		}

		var (
			srcuser   sql.NullString
			srcdomain string
		)

		srcuser.String, srcdomain, err = splitMailAddress(args[0])
		if err != nil {
			return err
		}

		// handle catchall alias
		if srcuser.String == "*" {
			srcuser.String = ""
			srcuser.Valid = false
		} else {
			srcuser.Valid = true
		}

		for _, dest := range args[1:] {
			dstuser, dstdomain, err := splitMailAddress(dest)
			if err != nil {
				return err
			}

			err = opts.db.CreateAlias(Alias{
				SourceUsername:      srcuser,
				SourceDomain:        srcdomain,
				DestinationUsername: dstuser,
				DestinationDomain:   dstdomain,
				Enabled:             true,
			})

			if err != nil {
				return fmt.Errorf("creating alias %v@%v -> %v@%v failed: %v",
					srcuser, srcdomain, dstuser, dstdomain, err)
			}
		}

		msg("alias created successfully")
		return nil
	},
}

func init() {
	cmdCreate.AddCommand(cmdCreateDomain)
	cmdCreate.AddCommand(cmdCreateMailbox)
	cmdCreate.AddCommand(cmdCreateAlias)
	root.AddCommand(cmdCreate)
}
