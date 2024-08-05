package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/starnuik/naive-pgd/pkg/pgp"
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

type buf []byte

func (b *buf) prependInt32(in int) {
	add := make([]byte, 4)
	binary.BigEndian.PutUint32(add, uint32(in))
	*b = append(add, *b...)
}

func (b *buf) addInt32(in int) {
	*b = binary.BigEndian.AppendUint32(*b, uint32(in))
}

func (b *buf) addString(in string) {
	nullTerm := append([]byte(in), '\000')
	*b = append(*b, nullTerm...)
}

func (b *buf) prependByte(in byte) {
	slice := make([]byte, 1)
	slice[0] = in
	*b = append(slice, *b...)
}

type StartupMessage struct {
	user string
	// database string
}

func (m *StartupMessage) Bytes() []byte {
	b := buf{}
	b.addInt32(196608)
	b.addString("user")
	b.addString(m.user)
	// b.addString("database")
	// b.addString(m.database)
	b.addString("")

	len := len(b) + 4
	b.prependInt32(len)

	return b
}

type AuthPasswordMessage struct {
	password string
}

func (m *AuthPasswordMessage) Bytes() []byte {
	b := buf{}
	b.addString(m.password)

	len := len(b) + 4
	b.prependInt32(len)

	b.prependByte('p')
	return b
}

type QueryMessage struct {
	query string
}

func (m *QueryMessage) Bytes() []byte {
	b := buf{}
	b.addString(m.query)

	len := len(b) + 4
	b.prependInt32(len)

	b.prependByte('Q')
	return b
}

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
	conn, err := net.Dial("tcp", pgUrl)
	requireNil(err)
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// send startup
	stup := StartupMessage{
		user: pgUser,
	}
	fmt.Println("--> REQUEST")
	fmt.Printf("%#v\n%s\n", stup, hex.Dump(stup.Bytes()))

	_, err = conn.Write(stup.Bytes())
	requireNil(err)

	// AuthenticationCleartextPassword
	printResponse(reader)

	// send auth
	pass := AuthPasswordMessage{
		password: pgPass,
	}
	fmt.Println("--> REQUEST")
	fmt.Printf("%#v\n%s\n", pass, hex.Dump(pass.Bytes()))

	_, err = conn.Write(pass.Bytes())
	requireNil(err)

	// AuthenticationOk
	printResponse(reader)
	// ParameterStatus-es
	printResponse(reader)
	printResponse(reader)
	printResponse(reader)
	printResponse(reader)
	printResponse(reader)
	printResponse(reader)
	printResponse(reader)
	printResponse(reader)
	printResponse(reader)
	printResponse(reader)
	printResponse(reader)
	printResponse(reader)
	printResponse(reader)
	printResponse(reader)

	// BackendKeyData
	printResponse(reader)

	// ReadyForQuery
	printResponse(reader)

	exec(conn, reader, "select * from users limit 1;")
	exec(conn, reader, "select usr_name, usr_id from users;")
	// exec(conn, reader, "insert into users (usr_name) values ('big dog');")
	// exec(conn, reader, "select usr_name, usr_id from users;")
	// exec(conn, reader, "delete from users where usr_name = 'big dog';")
	// exec(conn, reader, "select usr_name, usr_id from users;")
}

func exec(conn net.Conn, reader io.Reader, query string) {
	req := QueryMessage{
		query: query,
	}

	_, err := conn.Write(req.Bytes())
	requireNil(err)

	fmt.Printf("--> %#v\n", req)

	complete := false
	for !complete {
		res, err := pgp.Read(reader)
		if err != nil {
			fmt.Println(err)
			break
		}
		print(res)
		_, complete = res.(*pgp.ReadyForQuery)
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

func printResponse(reader io.Reader) {
	msg, err := pgp.Read(reader)
	if err != nil {
		fmt.Printf("<--- %s\n", err.Error())
		return
	}
	print(msg)
}

func print(msg pgp.ResponseMessage) {
	switch msg.(type) {
	// case *pgp.RowDescription:
	case *pgp.ErrorResponse:
		str, err := json.MarshalIndent(msg, "", "  ")
		requireNil(err)
		fmt.Printf("<---\n%s\n", str)
	default:
		fmt.Printf("<--- %#v\n", msg)
	}
}
