package pgwire

import (
	"bufio"
	"io"
	"net"
)

type Conn struct {
	tcp net.Conn
	w   io.Writer
	r   io.Reader
}

func NewConn(dbUrl string) (*Conn, error) {
	tcp, err := net.Dial("tcp", dbUrl)
	if err != nil {
		return nil, err
	}

	return &Conn{
		tcp: tcp,
		w:   tcp,
		r:   bufio.NewReader(tcp),
	}, nil
}

func (c *Conn) Close() error {
	return c.tcp.Close()
}
