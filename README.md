vmail - command-line tool for managing mailboxes and aliases for mail server
setups based on the [howto (English)](https://thomas-leister.de/en/mailserver-debian-stretch)
([German](https://thomas-leister.de/mailserver-debian-stretch/))
by [Thomas Leister](https://thomas-leister.de).

This one is written in Go, there's also a version in Rust here: https://github.com/awidegreen/vmail-rs

Rationale: I was unable to compile a static binary using Rust locally, and the
version available in Debian stable (stretch) was too old to build it. So I
wrote my own version of the program in Go, which results in a static binary.

Building
========

You need Go >= 1.11, then run the following command inside the checked-out repository:

    $ GO111MODULE=on go build

This will pull the needed dependencies, verify them cryptographically and build
a static binary called `vmail`.

Database Connection
===================

The `vmail` binary will try to connect to the MySQL socket
`/run/mysqld/mysqld.sock` as the current user and tries to use the `vmail`
database. If you need to connect to a different database, you can pass the connection string (format described [here](https://github.com/go-sql-driver/mysql#dsn-data-source-name)) in the environment variable `$VMAIL_DB`. The format is:

    [username[:password]@][protocol[(address)]]/dbname

For connecting a database on localhost via TCP with the user `foo`, the password `bar` and the database name `zzz`, use the following string:

    $ export VMAIL_DB='foo:bar@tcp(localhost:3306)/zzz'

Managing Domains, Mailboxes, and Aliases
========================================

Create new domain:

    $ vmail create domain example.com

Create new alias:

    $ vmail create alias foo@example.com bar@otherhost.example.com

Add a new email address to an existing alias:

    $ vmail create alias foo@example.com baz@otherhost.example.com

Create a mailbox:

    $ vmail create mailbox admin@example.com
    enter password:
    repeat password:
    mailbox admin@example.com created

List a domain:

    $ vmail show example.com
     Mailbox              Quota    Enabled
    ---------------------------------------
     admin@example.com             true
    ---------------------------------------

     Alias              Destinations
    ----------------------------------------------
     foo@example.com    bar@otherhost.example.com
                        baz@otherhost.example.com
    ----------------------------------------------

Add a catch-all alias (mind the quotes):

    $ vmail create alias '*@example.com' admin@example.com

List the domain again:

    $ vmail show example.com
         Mailbox              Quota    Enabled
    ---------------------------------------
     admin@example.com             true
    ---------------------------------------

     Alias              Destinations
    ----------------------------------------------
     *@example.com      admin@example.com
     foo@example.com    bar@otherhost.example.com
                        baz@otherhost.example.com
    ----------------------------------------------

List all domains:

    $ vmail domains
    example.com

Change the password for a mailbox:

    $ vmail password admin@example.com
    enter password:
    repeat password:
    password for admin@example.com updated
