package main

import (
	"github.com/oleksiivelychko/go-mysql-connection/connection"
	"github.com/oleksiivelychko/go-mysql-connection/querybuilder"
)

type Product struct {
	ID        int
	Name      string
	Price     float64
	SKU       string
	UpdatedAt string
}

func makeQueryBuilderMySQL() *querybuilder.Builder {
	conn, err := connection.New(&connection.Params{
		Username:     "gouser",
		Password:     "secret",
		DatabaseName: "go_mysql_connection",
		DriverName:   "mysql",
	})

	if err != nil {
		panic(err)
	}

	return querybuilder.New(conn)
}
