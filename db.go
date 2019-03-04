package main

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Domain is a domain for receiving email.
type Domain struct {
	ID     int
	Domain string
}

// Account is a mailbox.
type Account struct {
	ID       int
	Username string
	Domain   string
	Password string
	Quota    int
	Enabled  bool
	sendonly bool
}

// Alias forwards email to another destination.
type Alias struct {
	ID                  int
	SourceUsername      sql.NullString
	SourceDomain        string
	DestinationUsername string
	DestinationDomain   string
	enabled             bool
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
func (db *DB) CreateDomain(d Domain) error {
	res, err := db.Exec("INSERT INTO domains (domain) VALUES (?)", d.Domain)
	if err != nil {
		return err
	}

	fmt.Printf("result: %v\n", res)

	return nil
}

// UpdateDomain saves new values for d.
func (db *DB) UpdateDomain(d Domain) error {
	return nil
}

// FindDomain looks for a domain with the given name in the database.
func (db *DB) FindDomain(name string) (Domain, error) {
	var d Domain
	err := db.Get(&d, "SELECT * from domains WHERE domain = ?", name)
	if err != nil {
		return Domain{}, err
	}

	return d, nil
}
