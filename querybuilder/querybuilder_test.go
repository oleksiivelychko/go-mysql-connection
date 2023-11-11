package querybuilder

import (
	"github.com/oleksiivelychko/go-mysql-connection/connection"
	"testing"
)

func TestQueryBuilder_New(t *testing.T) {
	c, err := connection.New(&connection.Params{Username: "test", Password: "test", Database: "test", Driver: "mysql"})
	if c == nil {
		t.Fatal(err.Error())
	}

	if qb := New(c); qb.db() == nil {
		t.Error("unable to create")
	}
}
