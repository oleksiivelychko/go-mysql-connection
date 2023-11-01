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

func newQueryBuilderMySQL() *querybuilder.Builder {
	c, err := connection.New(&connection.Params{
		Username: "gopher",
		Password: "secret",
		Database: "go_mysql_connection",
		Driver:   "mysql",
	})

	if err != nil {
		panic(err)
	}

	return querybuilder.New(c)
}
