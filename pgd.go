package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"

	_ "github.com/joho/godotenv/autoload"
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

	stup := StartupMessage{
		user: pgUser,
	}
	fmt.Printf("%#v\n%s\n", stup, hex.Dump(stup.Bytes()))

	// setup tcp
	conn, err := net.Dial("tcp", pgUrl)
	requireNil(err)
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// send startup
	_, err = conn.Write(stup.Bytes())
	requireNil(err)

	read(reader)

	pass := AuthPasswordMessage{
		password: pgPass,
	}
	fmt.Printf("%#v\n%s\n", pass, hex.Dump(pass.Bytes()))

	_, err = conn.Write(pass.Bytes())
	requireNil(err)

	read(reader)
}

func read(reader io.Reader) {
	// response header
	resHead := make([]byte, 5)
	_, err := io.ReadFull(reader, resHead)
	requireNil(err)
	fmt.Printf("response header\n%s\n", hex.Dump(resHead))

	// response body
	resLen := binary.BigEndian.Uint32(resHead[1:]) - 4
	resBody := make([]byte, resLen)
	_, err = io.ReadFull(reader, resBody)
	requireNil(err)
	fmt.Printf("response body\n%s\n", hex.Dump(resBody))
}
