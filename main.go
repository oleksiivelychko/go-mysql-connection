package main

import "github.com/oleksiivelychko/go-utils/database"

func createQueryBuilder() *database.QueryBuilder {
	connection, err := database.NewConnection(&database.Params{
		Username:     "gouser",
		Password:     "secret",
		DatabaseName: "go_microservice",
		DriverName:   "mysql",
	})

	if err != nil {
		panic(err)
	}

	return &database.QueryBuilder{
		Connection: connection,
	}
}
