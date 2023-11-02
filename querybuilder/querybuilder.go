package querybuilder

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/oleksiivelychko/go-mysql-connection/connection"
	"strings"
)

type Builder struct{ Connection *connection.Connection }

func New(c *connection.Connection) *Builder { return &Builder{Connection: c} }

func (b *Builder) db() *sql.DB { return b.Connection.DB }

func (b *Builder) FindAll(table string) (*sql.Rows, error) {
	return b.db().Query(fmt.Sprintf(`SELECT * FROM %s`, table))
}

func (b *Builder) FindOne(table string, id uint) (*sql.Row, error) {
	ctx, cancel := b.Connection.ContextTimeout()
	defer cancel()

	stmt, err := b.db().PrepareContext(ctx, fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, table))
	if err != nil {
		return nil, err
	}

	defer func(s *sql.Stmt) { _ = s.Close() }(stmt)

	return stmt.QueryRowContext(ctx, id), nil
}

func (b *Builder) Insert(table string, data map[string]string) (int64, error) {
	columns, values := b.prepare(data, 0)

	res, err := b.db().ExecContext(
		context.Background(),
		fmt.Sprintf(
			`INSERT INTO %s (%s) VALUES (%s)`,
			table,
			columns,
			strings.TrimSuffix(strings.Repeat("?,", len(values)), ","),
		),
		values...,
	)

	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func (b *Builder) Update(table string, id uint, data map[string]string) (int64, error) {
	columns, values := b.prepare(data, id)

	res, err := b.db().ExecContext(
		context.Background(),
		fmt.Sprintf(`UPDATE %s SET %s WHERE id=?`, table, columns),
		values...,
	)

	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

func (b *Builder) Delete(table string, id uint) (int64, error) {
	res, err := b.db().Exec(fmt.Sprintf(`DELETE FROM %s WHERE id=?`, table), id)
	if err == nil {
		return res.RowsAffected()
	}

	return -1, fmt.Errorf("unable to delete from %s by id %d", table, id)
}

func (b *Builder) Truncate(table string) error {
	if _, err := b.Connection.DB.Exec(fmt.Sprintf(`TRUNCATE TABLE %s`, table)); err != nil {
		return err
	}

	return nil
}

func (b *Builder) AutoIncrement(table string) error {
	if _, err := b.Connection.DB.Exec(fmt.Sprintf(`ALTER TABLE %s AUTO_INCREMENT = 1`, table)); err != nil {
		return err
	}

	return nil
}

func (b *Builder) prepare(data map[string]string, id uint) (string, []any) {
	var columns, suffix string

	if id > 0 {
		suffix = "=?"
	}

	values := make([]any, 0, len(data))
	for k, v := range data {
		columns += fmt.Sprintf("%s%s,", k, suffix)
		values = append(values, v)
	}

	if id > 0 {
		values = append(values, id)
	}

	columns = strings.TrimSuffix(columns, ",")

	return columns, values
}
