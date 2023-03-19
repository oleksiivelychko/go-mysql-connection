package database

import (
	"database/sql"
	"github.com/oleksiivelychko/go-microservice/api"
	"github.com/oleksiivelychko/go-utils/connection_mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabase_SelectAllProduct(t *testing.T) {
	connectionMySQL, err := connection_mysql.NewConnectionMySQL(
		"gouser",
		"secret",
		"go_microservice",
	)
	if err != nil {
		t.Error(err)
	}

	results, err := connectionMySQL.SelectAll("products")
	if err != nil {
		t.Error(err)
	}

	var products []*api.Product

	for results.Next() {
		var product api.Product

		err = results.Scan(&product.ID, &product.Name, &product.Price, &product.SKU, &product.UpdatedAt)
		if err != nil {
			t.Error(err)
		}

		products = append(products, &product)
	}

	assert.Equal(t, products[0].ID, 1)
	assert.Equal(t, products[0].Name, "Latte")
	assert.Equal(t, products[0].Price, 1.49)
	assert.Equal(t, products[0].SKU, "123-456-789")
	assert.NotEmpty(t, products[0].UpdatedAt)
	assert.Equal(t, products[1].ID, 2)
	assert.Equal(t, products[1].Name, "Espresso")
	assert.Equal(t, products[1].Price, 0.99)
	assert.Equal(t, products[1].SKU, "000-000-001")
	assert.NotEmpty(t, products[1].UpdatedAt)

	defer func(results *sql.Rows) {
		_ = results.Close()
	}(results)

	defer func(connectionMySQL *connection_mysql.ConnectionMySQL) {
		_ = connectionMySQL.Close()
	}(connectionMySQL)
}
