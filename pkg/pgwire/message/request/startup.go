package request

import "github.com/starnuik/pogoda/pkg/pgwire/internal"

// ? https://www.postgresql.org/docs/16/protocol-message-formats.html#PROTOCOL-MESSAGE-FORMATS-STARTUPMESSAGE
type Startup struct {
	User     string
	Database string
}

func (m *Startup) Encode() []byte {
	out := internal.WriteBuf{}

	out.AddInt32(196608)
	out.AddString("user")
	out.AddString(m.User)
	out.AddString("database")
	out.AddString(m.Database)
	// "A zero byte is required as a terminator after the last name/value pair."
	out.AddString("")

	// https://www.postgresql.org/docs/16/protocol-overview.html#PROTOCOL-MESSAGE-CONCEPTS
	// "For historical reasons, the very first message sent by the client (the startup message) has no initial message-type byte."
	out.PrependLength()
	return out
}
