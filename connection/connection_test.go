package connection

import (
	"testing"
	"time"
)

func TestConnection_Successful(t *testing.T) {
	c, err := New(&Params{Username: "test", Password: "test", Database: "test", Driver: "mysql"})
	if c == nil {
		t.Error(err)
	}
}

func TestConnection_Failed(t *testing.T) {
	c, _ := New(&Params{Username: "", Password: "", Database: "", Driver: "mysql"})
	ping := c.DB.Ping()

	if ping.Error() == "" {
		t.Fatal("unable to ping connection")
	}

	t.Log(ping.Error())
}

func TestConnection_Closed(t *testing.T) {
	c, _ := New(&Params{Username: "test", Password: "test", Database: "test", Driver: "mysql"})
	err := c.Close()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestConnection_DriverMissed(t *testing.T) {
	_, err := New(&Params{Username: "test", Password: "test", Database: "test", Driver: ""})
	if err == nil {
		t.Error("connection should not have been open, driver name was not specified")
	}
}

func TestConnection_DefaultParams(t *testing.T) {
	c, _ := New(&Params{Username: "test", Password: "test", Database: "test", Driver: "mysql"})

	if c.Params.MaxLifetimeMinutes != time.Minute*maxLifetime {
		t.Error("unable to assign default value to MaxLifetimeMinutes")
	}

	if c.Params.MaxOpenConnections != maxOpenConnections {
		t.Error("unable to assign default value to MaxOpenConnections")
	}

	if c.Params.MaxIdleConnections != maxIdleConnections {
		t.Error("unable to assign default value to MaxIdleConnections")
	}

	if c.Params.ContextTimeoutSeconds != time.Second*contextTimeout {
		t.Error("unable to assign default value to ContextTimeoutSeconds")
	}
}

func TestConnection_Params(t *testing.T) {
	c, _ := New(&Params{
		Username:              "test",
		Password:              "test",
		Database:              "test",
		Driver:                "mysql",
		MaxLifetimeMinutes:    1,
		MaxOpenConnections:    1,
		MaxIdleConnections:    1,
		ContextTimeoutSeconds: 1,
	})

	// compare to 1 second
	if c.Params.MaxLifetimeMinutes/60 != 1000000000 {
		t.Error("unable to assign value to MaxLifetimeMinutes")
	}

	if c.Params.MaxOpenConnections != 1 {
		t.Error("unable to assign value to MaxOpenConnections")
	}

	if c.Params.MaxIdleConnections != 1 {
		t.Error("unable to assign value to MaxIdleConnections")
	}

	// compare to 1 second
	if c.Params.ContextTimeoutSeconds != 1000000000 {
		t.Error("unable to assign value to ContextTimeoutSeconds")
	}
}
