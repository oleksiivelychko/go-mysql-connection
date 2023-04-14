package connection

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const defaultMaxLifetime = 2
const defaultMaxOpenConnections = 5
const defaultMaxIdleConnections = 5
const defaultContextTimeout = 10

type Params struct {
	Username              string
	Password              string
	DatabaseName          string
	DriverName            string
	MaxLifetimeMinutes    time.Duration
	MaxOpenConnections    int
	MaxIdleConnections    int
	ContextTimeoutSeconds time.Duration
}

type Connection struct {
	Params *Params
	DB     *sql.DB
}

func New(params *Params) (*Connection, error) {
	db, err := driverFactory(params)
	if err != nil {
		return nil, err
	}

	if params.MaxLifetimeMinutes == 0 {
		params.MaxLifetimeMinutes = time.Minute * defaultMaxLifetime
	} else {
		params.MaxLifetimeMinutes = time.Minute * params.MaxLifetimeMinutes
	}

	if params.MaxOpenConnections == 0 {
		params.MaxOpenConnections = defaultMaxOpenConnections
	}

	if params.MaxIdleConnections == 0 {
		params.MaxIdleConnections = defaultMaxIdleConnections
	}

	if params.ContextTimeoutSeconds == 0 {
		params.ContextTimeoutSeconds = time.Second * defaultContextTimeout
	} else {
		params.ContextTimeoutSeconds = time.Second * params.ContextTimeoutSeconds
	}

	/*
		It's required to ensure connections are closed by the driver safely before Connection is closed by MySQL server, OS, or other middlewares.
		Since some middlewares close idle connections by 5 minutes, we recommend timeout shorter than 5 minutes.
		This setting helps load balancing and changing system variables too.
	*/
	db.SetConnMaxLifetime(params.MaxLifetimeMinutes)
	/*
		It's highly recommended to limit the number of Connection used by the application.
		There is no recommended limit number because it depends on application and MySQL server.
	*/
	db.SetMaxOpenConns(params.MaxOpenConnections)
	/*
		It's recommended to be set same to db.SetMaxOpenConns().
		When it is smaller than SetMaxOpenConns(), connections can be opened and closed much more frequently than you expect.
		Idle connections can be closed by the db.SetConnMaxLifetime().
		If you want to close idle connections more rapidly, you can use db.SetConnMaxIdleTime() since Go 1.15.
	*/
	db.SetMaxIdleConns(params.MaxIdleConnections)

	return &Connection{DB: db, Params: params}, nil
}

func (conn *Connection) Close() error {
	return conn.DB.Close()
}

func (conn *Connection) ContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), conn.Params.ContextTimeoutSeconds)
}

func driverFactory(params *Params) (*sql.DB, error) {
	switch params.DriverName {
	case "mysql":
		return sql.Open(params.DriverName, fmt.Sprintf("%s:%s@/%s", params.Username, params.Password, params.DatabaseName))
	}

	return nil, fmt.Errorf("unable to open connection, driver name is not provided")
}
