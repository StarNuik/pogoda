package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/starnuik/naive-pgd/pkg/pgwire"
	"github.com/starnuik/naive-pgd/pkg/pgwire/message"
)

// type pgdDriver struct{}

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

func requireNil(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// fmt.Println(conn, err)
	// conn.Close()
	// fmt.Println(conn, err)
	// pgVersion := 196608
	// bytes := make([]byte, 4)
	// binary.BigEndian.PutUint32(bytes, uint32(pgVersion))
	pgUser := os.Getenv("PG_USER")
	pgUrl := os.Getenv("PG_URL")
	pgPass := os.Getenv("PG_PASSWORD")

	// setup tcp
	conn, err := pgwire.NewConn(pgUrl)
	requireNil(err)
	defer conn.Close()

	// send startup
	stup := &message.Startup{
		User: pgUser,
	}
	fmt.Println("--> REQUEST")
	fmt.Printf("%#v\n%s\n", stup, hex.Dump(stup.Bytes()))

	err = conn.Write(stup)
	requireNil(err)

	// AuthenticationCleartextPassword
	printResponse(conn)

	// send auth
	pass := &message.Password{
		Password: pgPass,
	}

	err = conn.Write(pass)
	requireNil(err)

	// AuthenticationOk
	printResponse(conn)
	// ParameterStatus-es
	printResponse(conn)
	printResponse(conn)
	printResponse(conn)
	printResponse(conn)
	printResponse(conn)
	printResponse(conn)
	printResponse(conn)
	printResponse(conn)
	printResponse(conn)
	printResponse(conn)
	printResponse(conn)
	printResponse(conn)
	printResponse(conn)
	printResponse(conn)

	// BackendKeyData
	printResponse(conn)

	// ReadyForQuery
	printResponse(conn)

	exec(conn, "select * from users limit 1;")
	exec(conn, "select usr_name, usr_id from users;")
	// exec(conn, reader, "insert into users (usr_name) values ('big dog');")
	// exec(conn, reader, "select usr_name, usr_id from users;")
	// exec(conn, reader, "delete from users where usr_name = 'big dog';")
	// exec(conn, reader, "select usr_name, usr_id from users;")
}

func exec(conn *pgwire.Conn, query string) {
	req := &message.Query{
		Query: query,
	}

	err := conn.Write(req)
	requireNil(err)

	fmt.Printf("--> %#v\n", req)

	complete := false
	for !complete {
		res, err := conn.Read()
		if err != nil {
			fmt.Println(err)
			break
		}
		print(res)
		_, complete = res.(*message.ReadyForQuery)
	}
	// fmt.Printf("%#v\n%s\n", query, hex.Dump(query.Bytes()))
}

// AuthenticationCleartextPassword
// AuthenticationOk
// ParameterStatus
// BackendKeyData
// ReadyForQuery
// RowDescription
// DataRow
// CommandComplete

func printResponse(conn *pgwire.Conn) {
	msg, err := conn.Read()
	if err != nil {
		fmt.Printf("<--- %s\n", err.Error())
		return
	}
	print(msg)
}

func print(msg message.Response) {
	switch msg.(type) {
	// case *pgp.RowDescription:
	case *message.ErrorResponse:
		str, err := json.MarshalIndent(msg, "", "  ")
		requireNil(err)
		fmt.Printf("<---\n%s\n", str)
	default:
		fmt.Printf("<--- %#v\n", msg)
	}
}
