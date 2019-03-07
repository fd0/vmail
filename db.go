package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Domain is a domain for receiving email.
type Domain struct {
	ID     int
	Domain string
}

// DB stores domains, accounts and aliases.
type DB struct {
	*sqlx.DB
}

// ConnectDB opens a connection to the database.
func ConnectDB(driver, source string) (*DB, error) {
	db, err := sqlx.Connect(driver, source)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// CreateDomain creates a new domain d.
func (db *DB) CreateDomain(name string) error {
	_, err := db.Exec("INSERT INTO domains (domain) VALUES (?)", name)
	if err != nil {
		return err
	}

	return nil
}

// FindDomain looks for a domain with the given name in the database.
func (db *DB) FindDomain(name string) (Domain, error) {
	var d Domain
	err := db.Get(&d, "SELECT * from domains WHERE domain = ?", name)
	if err != nil {
		return Domain{}, fmt.Errorf("domain not found: %v", err)
	}

	return d, nil
}

// FindAllDomains returns a list of all domains which contain name.
func (db *DB) FindAllDomains(name string) ([]Domain, error) {
	var ds []Domain
	err := db.Select(&ds, "SELECT * from domains WHERE domain LIKE ? ORDER BY domain", "%"+name+"%")
	if err != nil {
		return nil, err
	}

	return ds, nil
}

// DeleteDomain removes a domain.
func (db *DB) DeleteDomain(name string) error {
	res, err := db.Exec("DELETE FROM domains WHERE domain = ?", name)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("not found")
	}

	return nil
}

// Account is a mailbox.
type Account struct {
	ID       int
	Username string
	Domain   string
	Password string
	Quota    int
	Enabled  bool
	Sendonly bool
}

// CreateAccount creates a new mailbox for the domain d.
func (db *DB) CreateAccount(a Account) error {
	_, err := db.Exec(`INSERT INTO accounts
		(username, domain, password, quota, enabled, sendonly)
		VALUES (?, ?, ?, ?, ?, ?)`,
		a.Username, a.Domain, a.Password, a.Quota, a.Enabled, a.Sendonly)
	if err != nil {
		return err
	}

	return nil
}

// FindAllAccounts returns a list of all accounts for a domain.
func (db *DB) FindAllAccounts(domain string) ([]Account, error) {
	var accounts []Account
	err := db.Select(&accounts, "SELECT * from accounts WHERE domain = ?", domain)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

// Alias forwards email to another destination.
type Alias struct {
	ID                  int            `db:"id"`
	SourceUsername      sql.NullString `db:"source_username"`
	SourceDomain        string         `db:"source_domain"`
	DestinationUsername string         `db:"destination_username"`
	DestinationDomain   string         `db:"destination_domain"`
	Enabled             bool           `db:"enabled"`
}

// FindAllAliases returns a list of all aliases for a domain.
func (db *DB) FindAllAliases(domain string) ([]Alias, error) {
	var aliases []Alias
	err := db.Select(&aliases, "SELECT * from aliases WHERE source_domain = ? ORDER BY source_username", domain)
	if err != nil {
		return nil, err
	}

	return aliases, nil
}
