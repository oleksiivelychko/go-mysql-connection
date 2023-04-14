package connection

import (
	"testing"
	"time"
)

func TestConnection_Successful(t *testing.T) {
	conn, err := New(&Params{Username: "test", Password: "test", DatabaseName: "test", DriverName: "mysql"})
	if conn == nil {
		t.Fatal(err.Error())
	}
}

func TestConnection_Failed(t *testing.T) {
	conn, _ := New(&Params{Username: "", Password: "", DatabaseName: "", DriverName: "mysql"})
	ping := conn.DB.Ping()

	if ping.Error() == "" {
		t.Fatal("unable to check connection")
	}

	t.Log(ping.Error())
}

func TestConnection_Close(t *testing.T) {
	conn, _ := New(&Params{Username: "test", Password: "test", DatabaseName: "test", DriverName: "mysql"})
	err := conn.Close()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestConnection_DriverFactoryFailed(t *testing.T) {
	_, err := New(&Params{Username: "test", Password: "test", DatabaseName: "test", DriverName: ""})
	if err == nil {
		t.Fatal("unable to open connection, driver name is not provided")
	}
}

func TestConnection_DefaultParams(t *testing.T) {
	conn, _ := New(&Params{Username: "test", Password: "test", DatabaseName: "test", DriverName: "mysql"})

	if conn.Params.MaxLifetimeMinutes != time.Minute*defaultMaxLifetime {
		t.Error("unable to assign default value to MaxLifetimeMinutes")
	}

	if conn.Params.MaxOpenConnections != defaultMaxOpenConnections {
		t.Error("unable to assign default value to MaxOpenConnections")
	}

	if conn.Params.MaxIdleConnections != defaultMaxIdleConnections {
		t.Error("unable to assign default value to MaxIdleConnections")
	}

	if conn.Params.ContextTimeoutSeconds != time.Second*defaultContextTimeout {
		t.Error("unable to assign default value to ContextTimeoutSeconds")
	}
}

func TestConnection_Params(t *testing.T) {
	conn, _ := New(&Params{
		Username:              "test",
		Password:              "test",
		DatabaseName:          "test",
		DriverName:            "mysql",
		MaxLifetimeMinutes:    1,
		MaxOpenConnections:    1,
		MaxIdleConnections:    1,
		ContextTimeoutSeconds: 1,
	})

	// compare to 1 second
	if conn.Params.MaxLifetimeMinutes/60 != 1000000000 {
		t.Error("unable to assign value to MaxLifetimeMinutes")
	}

	if conn.Params.MaxOpenConnections != 1 {
		t.Error("unable to assign value to MaxOpenConnections")
	}

	if conn.Params.MaxIdleConnections != 1 {
		t.Error("unable to assign value to MaxIdleConnections")
	}

	// compare to 1 second
	if conn.Params.ContextTimeoutSeconds != 1000000000 {
		t.Error("unable to assign value to ContextTimeoutSeconds")
	}
}
