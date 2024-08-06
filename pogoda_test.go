package pogoda_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/starnuik/pogoda"
	"github.com/starnuik/pogoda/pkg/pgwire/message"
	"github.com/starnuik/pogoda/pkg/pgwire/message/request"
	"github.com/stretchr/testify/require"
)

var (
	pgUser = os.Getenv("PG_USER")
	pgUrl  = os.Getenv("PG_URL")
	pgPass = os.Getenv("PG_PASSWORD")
)

func TestMain(m *testing.M) {
	//
	m.Run()
}

func TestPogoda(t *testing.T) {
	require := require.New(t)

	conn, err := pogoda.NewConn(pgUrl)
	require.Nil(err)
	defer conn.Close()

	err = conn.Auth(pgUser, pgPass, "")
	require.Nil(err)

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
func exec(conn *pogoda.Conn, query string) {
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

func printResponse(conn *pogoda.Conn) {
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
		str, _ := json.MarshalIndent(msg, "", "  ")
		fmt.Printf("<---\n%s\n", str)
	default:
		fmt.Printf("<--- %#v\n", msg)
	}
}
