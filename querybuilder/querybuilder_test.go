package querybuilder

import (
	"github.com/oleksiivelychko/go-mysql-connection/connection"
	"testing"
)

func TestQueryBuilder_New(t *testing.T) {
	conn, err := connection.New(&connection.Params{
		Username:     "test",
		Password:     "test",
		DatabaseName: "test",
		DriverName:   "mysql",
	})

	if conn == nil {
		t.Fatal(err.Error())
	}

	queryBuilder := New(conn)
	if queryBuilder.db() == nil {
		t.Fatal("unable to create query builder")
	}
}
