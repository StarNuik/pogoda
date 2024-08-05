package main

import (
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/starnuik/pogoda/pkg/pgwire"
	"github.com/starnuik/pogoda/pkg/pgwire/message"
	"github.com/starnuik/pogoda/pkg/pgwire/message/request"
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

	err = conn.Auth(pgUser, pgPass, "")
	requireNil(err)

	exec(conn, "select * from users limit 1;")
	exec(conn, "select usr_name, usr_id from users;")
	exec(conn, "select * from users limit 1; select usr_name, usr_id from users;")

	conn.Write(&request.Terminate{})
	printResponse(conn)
	// exec(conn, reader, "insert into users (usr_name) values ('big dog');")
	// exec(conn, reader, "select usr_name, usr_id from users;")
	// exec(conn, reader, "delete from users where usr_name = 'big dog';")
	// exec(conn, reader, "select usr_name, usr_id from users;")
}

// "The simple Query message is approximately equivalent to the series
// Parse, Bind, portal Describe, Execute, Close, Sync,
// using the unnamed prepared statement and portal objects and no parameters."
func exec(conn *pgwire.Conn, query string) {
	req := &request.Query{
		Query: query,
	}

	err := conn.Write(req)
	if err != nil {
		fmt.Printf("ERR> %v\n", err)
	}

	fmt.Printf("---> %#v\n", req)

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

func printResponse(conn *pgwire.Conn) {
	msg, err := conn.Read()
	if err != nil {
		fmt.Printf("<ERR %s\n", err.Error())
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
