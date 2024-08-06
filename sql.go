package pogoda

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
)

// ? https://pkg.go.dev/database/sql/driver@go1.22.0

func init() {
	f := factory{}
	sql.Register("pogoda", &f)
}

var (
	_ driver.Driver        = &factory{}
	_ driver.DriverContext = &factory{}
	_ driver.Connector     = &factory{}
)

// "Drivers should implement Connector and DriverContext interfaces."
type factory struct{}

// Driver interface
func (f *factory) Open(name string) (driver.Conn, error) {
	fmt.Println("f.Open")
	c, err := f.OpenConnector(name)
	if err != nil {
		return nil, err
	}
	return c.Connect(context.TODO())
}

// DriverContext interface
func (f *factory) OpenConnector(name string) (driver.Connector, error) {
	fmt.Println("f.OpenConnector")
	return f, nil
}

// Connector interface
func (f *factory) Connect(ctx context.Context) (driver.Conn, error) {
	fmt.Println("f.Connect")
	return nil, nil
}

// Connector interface
func (f *factory) Driver() driver.Driver {
	fmt.Println("f.Driver")
	return f
}
