package querybuilder

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/oleksiivelychko/go-mysql-connection/connection"
	"strings"
)

type Builder struct {
	Connection *connection.Connection
}

func New(connection *connection.Connection) *Builder {
	return &Builder{Connection: connection}
}

func (builder *Builder) db() *sql.DB {
	return builder.Connection.DB
}

func (builder *Builder) FindAll(table string) (*sql.Rows, error) {
	query := fmt.Sprintf(`SELECT * FROM %s`, table)
	return builder.db().Query(query)
}

func (builder *Builder) FindOne(table string, id uint) (*sql.Row, error) {
	ctx, cancelFunc := builder.Connection.ContextTimeout()
	defer cancelFunc()

	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, table)
	stmt, err := builder.db().PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	return stmt.QueryRowContext(ctx, id), nil
}

func (builder *Builder) Insert(table string, data map[string]string) (int64, error) {
	columns, values := builder.prepareData(data, 0)
	params := strings.TrimSuffix(strings.Repeat("?,", len(values)), ",")

	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, table, columns, params)
	result, err := builder.db().ExecContext(context.Background(), query, values...)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}

func (builder *Builder) Update(table string, id uint, data map[string]string) (int64, error) {
	columns, values := builder.prepareData(data, id)

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id=?`, table, columns)
	result, err := builder.db().ExecContext(context.Background(), query, values...)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func (builder *Builder) Delete(table string, id uint) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=?`, table)
	result, err := builder.db().Exec(query, id)

	if err == nil {
		return result.RowsAffected()
	}

	return -1, fmt.Errorf("unable to delete from %s by id %d", table, id)
}

func (builder *Builder) Truncate(table string) error {
	query := fmt.Sprintf(`TRUNCATE TABLE %s`, table)
	if _, err := builder.Connection.DB.Exec(query); err != nil {
		return err
	}

	return nil
}

func (builder *Builder) AutoIncrement(table string) error {
	query := fmt.Sprintf(`ALTER TABLE %s AUTO_INCREMENT = 1`, table)
	if _, err := builder.Connection.DB.Exec(query); err != nil {
		return err
	}

	return nil
}

func (builder *Builder) prepareData(data map[string]string, id uint) (string, []any) {
	columns := *new(string)
	values := make([]any, 0, len(data))

	suffix := ""
	if id > 0 {
		suffix = "=?"
	}

	for key, value := range data {
		columns += fmt.Sprintf("%s%s,", key, suffix)
		values = append(values, value)
	}

	if id > 0 {
		values = append(values, id)
	}

	columns = strings.TrimSuffix(columns, ",")

	return columns, values
}
