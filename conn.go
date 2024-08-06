package pogoda

import (
	"bufio"
	"net"

	"github.com/starnuik/pogoda/pkg/pgwire"
)

type Conn struct {
	tcp net.Conn
	*pgwire.Wire
}

func NewConn(dbUrl string) (*Conn, error) {
	tcp, err := net.Dial("tcp", dbUrl)
	if err != nil {
		return nil, err
	}

	wire := pgwire.NewWire(tcp, bufio.NewReader(tcp))
	return &Conn{
		Wire: wire,
		tcp:  tcp,
	}, nil
}

func (c *Conn) Close() error {
	return c.tcp.Close()
}

// // ? https://github.com/jbarham/gopgsqldriver/blob/f8287ee9bfe224aa4a7edcd73815ecbe69db7f68/pgdriver.go#L350
// func init() {
// 	sql.Register("pgd", &pgdDriver{})
// }

// // ? https://pkg.go.dev/database/sql/driver#Driver
// // This is the MINIMAL requirement for an sql.driver.Driver
// // Anything else is an "IF a Driver implements"
// func (d *pgdDriver) Open(name string) (driver.Conn, error) {
// 	fmt.Println("pgdDriver.Open:", name)
// 	conn, err := net.Dial("tcp", "localhost:5432")
// 	// c = NewConnector
// 	// c.Dialer(defaultDialer{})
// 	// c.open(ctx.Background)
// 	//

// 	panic("not implemented")
// }
