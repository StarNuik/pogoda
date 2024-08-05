package message

import "github.com/starnuik/naive-pgd/pkg/pgwire/internal"

type Startup struct {
	User string
}

type Password struct {
	Password string
}

type Query struct {
	Query string
}

// https://www.postgresql.org/docs/16/protocol-message-formats.html#PROTOCOL-MESSAGE-FORMATS-STARTUPMESSAGE
func (m *Startup) Bytes() []byte {
	out := internal.WriteBuf{}
	// 0x00030000
	out.Int32(196608)
	out.String("user")
	out.String(m.User)
	// "A zero byte is required as a terminator after the last name/value pair."
	out.String("")

	out.PrependLength()
	return out
}

func (m *Password) Bytes() []byte {
	out := internal.WriteBuf{}
	out.String(m.Password)

	out.PrependHeader('p')
	return out
}

func (m *Query) Bytes() []byte {
	out := internal.WriteBuf{}
	out.String(m.Query)

	out.PrependHeader('Q')
	return out
}
