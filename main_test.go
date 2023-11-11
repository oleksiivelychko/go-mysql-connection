package main

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	username = "gopher"
	password = "secret"
	database = "go_mysql_connection"
	driver   = "mysql"
	table    = "products"
)

func TestFindAll(t *testing.T) {
	var qb = newQueryBuilderMySQL(username, password, database, driver)

	rows, err := qb.FindAll(table)
	if err != nil {
		t.Fatal(err)
	}

	var products []*product

	for rows.Next() {
		var p product

		err = rows.Scan(&p.ID, &p.Name, &p.Price, &p.SKU, &p.UpdatedAt)
		if err != nil {
			t.Fatal(err)
		}

		products = append(products, &p)
	}

	assert.Equal(t, products[0].ID, 1)
	assert.Equal(t, products[0].Name, "Latte")
	assert.Equal(t, products[0].Price, 1.49)
	assert.Equal(t, products[0].SKU, "123-456-789")
	assert.NotEmpty(t, products[0].UpdatedAt)

	if _, err = time.Parse(time.DateTime, products[0].UpdatedAt); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, products[1].ID, 2)
	assert.Equal(t, products[1].Name, "Espresso")
	assert.Equal(t, products[1].Price, 0.99)
	assert.Equal(t, products[1].SKU, "000-000-001")
	assert.NotEmpty(t, products[1].UpdatedAt)

	if _, err = time.Parse(time.DateTime, products[1].UpdatedAt); err != nil {
		t.Fatal(err)
	}

	defer func(r *sql.Rows) { _ = r.Close() }(rows)

	err = qb.Connection.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestFindOne(t *testing.T) {
	var qb = newQueryBuilderMySQL(username, password, database, driver)

	row, err := qb.FindOne(table, 1)
	if err != nil {
		t.Fatal(err)
	}

	var p product

	err = row.Scan(&p.ID, &p.Name, &p.Price, &p.SKU, &p.UpdatedAt)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, p.ID, 1)
	assert.Equal(t, p.Name, "Latte")
	assert.Equal(t, p.Price, 1.49)
	assert.Equal(t, p.SKU, "123-456-789")
	assert.NotEmpty(t, p.UpdatedAt)

	if err = qb.Connection.Close(); err != nil {
		t.Error(err)
	}
}

func TestInsert(t *testing.T) {
	var qb = newQueryBuilderMySQL("gopher", "secret", "go_mysql_connection", "mysql")

	id, err := qb.Insert(table, map[string]string{
		"name":  "Coffee",
		"price": "0.49",
		"SKU":   "123-000-000",
	})

	if err != nil {
		t.Fatal(err)
	}

	if id < 0 {
		t.Errorf("expected positive, got %d", id)
	}
}

func TestUpdate(t *testing.T) {
	var qb = newQueryBuilderMySQL(username, password, database, driver)

	affected, err := qb.Update(table, 3, map[string]string{
		"name":  "Tea",
		"price": "0.09",
		"SKU":   "123-000-789",
	})

	if err != nil {
		t.Fatal(err)
	}

	if affected < 0 {
		t.Errorf("expected positive, got %d", affected)
	}
}

func TestDelete(t *testing.T) {
	var qb = newQueryBuilderMySQL(username, password, database, driver)

	affected, err := qb.Delete(table, 3)
	if err != nil {
		t.Fatal(err)
	}

	if affected < 0 {
		t.Fatalf("expected positive, got %d", affected)
	}

	if err = qb.AutoIncrement(table); err != nil {
		t.Error(err)
	}
}
