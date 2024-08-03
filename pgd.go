package pgd

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type pgdDriver struct{}

func (d *pgdDriver) Open(name string) (driver.Conn, error) {
	fmt.Println("pgd.Open:", name)
	panic("not implemented")
}

func init() {
	sql.Register("pgd", &pgdDriver{})
}
