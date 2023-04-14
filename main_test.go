package main

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFindAll(t *testing.T) {
	queryBuilder := makeQueryBuilderMySQL()

	results, err := queryBuilder.FindAll("products")
	if err != nil {
		t.Fatal(err)
	}

	var products []*Product

	for results.Next() {
		var product Product

		err = results.Scan(&product.ID, &product.Name, &product.Price, &product.SKU, &product.UpdatedAt)
		if err != nil {
			t.Fatal(err)
		}

		products = append(products, &product)
	}

	assert.Equal(t, products[0].ID, 1)
	assert.Equal(t, products[0].Name, "Latte")
	assert.Equal(t, products[0].Price, 1.49)
	assert.Equal(t, products[0].SKU, "123-456-789")
	assert.NotEmpty(t, products[0].UpdatedAt)

	_, err = time.Parse(time.DateTime, products[0].UpdatedAt)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, products[1].ID, 2)
	assert.Equal(t, products[1].Name, "Espresso")
	assert.Equal(t, products[1].Price, 0.99)
	assert.Equal(t, products[1].SKU, "000-000-001")
	assert.NotEmpty(t, products[1].UpdatedAt)

	_, err = time.Parse(time.DateTime, products[1].UpdatedAt)
	if err != nil {
		t.Fatal(err)
	}

	defer func(results *sql.Rows) {
		_ = results.Close()
	}(results)

	err = queryBuilder.Connection.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestFindOne(t *testing.T) {
	queryBuilder := makeQueryBuilderMySQL()

	result, err := queryBuilder.FindOne("products", 1)
	if err != nil {
		t.Fatal(err)
	}

	var product Product

	err = result.Scan(&product.ID, &product.Name, &product.Price, &product.SKU, &product.UpdatedAt)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, product.ID, 1)
	assert.Equal(t, product.Name, "Latte")
	assert.Equal(t, product.Price, 1.49)
	assert.Equal(t, product.SKU, "123-456-789")
	assert.NotEmpty(t, product.UpdatedAt)

	err = queryBuilder.Connection.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestInsert(t *testing.T) {
	queryBuilder := makeQueryBuilderMySQL()

	lastInsertID, err := queryBuilder.Insert("products", map[string]string{
		"name":  "Coffee",
		"price": "0.49",
		"SKU":   "123-000-000",
	})

	if err != nil {
		t.Fatal(err)
	}

	if lastInsertID < 0 {
		t.Errorf("unable to insert data, got negative %d", lastInsertID)
	}
}

func TestUpdate(t *testing.T) {
	queryBuilder := makeQueryBuilderMySQL()

	rowsAffected, err := queryBuilder.Update("products", 3, map[string]string{
		"name":  "Tea",
		"price": "0.09",
		"SKU":   "123-000-789",
	})

	if err != nil {
		t.Fatal(err)
	}

	if rowsAffected < 0 {
		t.Errorf("unable to update data, got negative %d", rowsAffected)
	}
}

func TestDelete(t *testing.T) {
	queryBuilder := makeQueryBuilderMySQL()

	rowsAffected, err := queryBuilder.Delete("products", 3)
	if err != nil {
		t.Fatal(err)
	}

	if rowsAffected < 0 {
		t.Fatalf("unable to delete data, got negative %d", rowsAffected)
	}

	err = queryBuilder.AutoIncrement("products")
	if err != nil {
		t.Error(err)
	}
}
