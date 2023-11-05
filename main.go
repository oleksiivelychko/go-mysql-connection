package main

import (
	"github.com/oleksiivelychko/go-mysql-connection/connection"
	"github.com/oleksiivelychko/go-mysql-connection/querybuilder"
)

type product struct {
	ID        int
	Name      string
	Price     float64
	SKU       string
	UpdatedAt string
}

func newQueryBuilderMySQL(username, password, database, driver string) *querybuilder.Builder {
	c, err := connection.New(&connection.Params{
		Username: username,
		Password: password,
		Database: database,
		Driver:   driver,
	})

	if err != nil {
		panic(err)
	}

	return querybuilder.New(c)
}
