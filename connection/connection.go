package connection

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	maxLifetime        = 2
	maxOpenConnections = 5
	maxIdleConnections = 5
	contextTimeout     = 5
)

type Params struct {
	Username              string
	Password              string
	Database              string
	Driver                string
	MaxLifetimeMinutes    time.Duration
	MaxOpenConnections    int
	MaxIdleConnections    int
	ContextTimeoutSeconds time.Duration
}

type Connection struct {
	Params *Params
	DB     *sql.DB
}

func New(p *Params) (*Connection, error) {
	db, err := p.db()
	if err != nil {
		return nil, err
	}

	if p.MaxLifetimeMinutes == 0 {
		p.MaxLifetimeMinutes = time.Minute * maxLifetime
	} else {
		p.MaxLifetimeMinutes = time.Minute * p.MaxLifetimeMinutes
	}

	if p.MaxOpenConnections == 0 {
		p.MaxOpenConnections = maxOpenConnections
	}

	if p.MaxIdleConnections == 0 {
		p.MaxIdleConnections = maxIdleConnections
	}

	if p.ContextTimeoutSeconds == 0 {
		p.ContextTimeoutSeconds = time.Second * contextTimeout
	} else {
		p.ContextTimeoutSeconds = time.Second * p.ContextTimeoutSeconds
	}

	/*
		It's required to ensure connections are closed by the driver safely before Connection is closed by MySQL server, OS, or other middlewares.
		Since some middlewares close idle connections by 5 minutes, we recommend timeout shorter than 5 minutes.
		This setting helps load balancing and changing system variables too.
	*/
	db.SetConnMaxLifetime(p.MaxLifetimeMinutes)
	/*
		It's highly recommended to limit the number of Connection used by the application.
		There is no recommended limit number because it depends on application and MySQL server.
	*/
	db.SetMaxOpenConns(p.MaxOpenConnections)
	/*
		It's recommended to be set same to db.SetMaxOpenConns().
		When it is smaller than SetMaxOpenConns(), connections can be opened and closed much more frequently than you expect.
		Idle connections can be closed by the db.SetConnMaxLifetime().
		If you want to close idle connections more rapidly, you can use db.SetConnMaxIdleTime() since Go 1.15.
	*/
	db.SetMaxIdleConns(p.MaxIdleConnections)

	return &Connection{p, db}, nil
}

func (c *Connection) Close() error {
	return c.DB.Close()
}

func (c *Connection) ContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.Params.ContextTimeoutSeconds)
}

func (p Params) db() (*sql.DB, error) {
	switch p.Driver {
	case "mysql":
		return sql.Open(p.Driver, fmt.Sprintf("%s:%s@/%s", p.Username, p.Password, p.Database))
	}

	return nil, fmt.Errorf("driver %s is not supported", p.Driver)
}
