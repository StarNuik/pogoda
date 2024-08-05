package pgwire

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/starnuik/pogoda/pkg/pgwire/internal"
	"github.com/starnuik/pogoda/pkg/pgwire/message"
)

type messageHeader struct {
	id     byte
	length int32
}

type fullMessage struct {
	messageHeader
	body internal.ReadBuf
}

// Read retrieves the first pending message from the tcp connection.
// An error is returned on a db initiated ErrorResponse
func (c *Conn) Read() (message.Response, error) {
	msg, err := read(c.r)
	if err != nil {
		return nil, err
	}

	out, err := deduceResponseType(msg)
	if err != nil {
		return nil, err
	}

	if len(msg.body) != int(msg.length) {
		return nil, fmt.Errorf("invalid message length")
	}

	err = out.Populate(msg.body)
	if err != nil {
		return nil, err
	}

	if err, is := out.(*message.ErrorResponse); is {
		return nil, err
	}

	return out, nil
}

func read(reader io.Reader) (*fullMessage, error) {
	head := messageHeader{}

	buf := make([]byte, 5)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}

	head = toHeader(buf)

	buf = make([]byte, head.length)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}
	return &fullMessage{
		messageHeader: head,
		body:          buf,
	}, nil
}

// WILL panic on len(head) < 5
func toHeader(bytes []byte) messageHeader {
	h := messageHeader{}
	h.id = bytes[0]
	h.length = int32(binary.BigEndian.Uint32(bytes[1:])) - 4
	return h
}

func deduceResponseType(in *fullMessage) (message.Response, error) {
	var out message.Response
	switch in.id {
	case 'R':
		return deduceAuth(in)
	case 'S':
		out = &message.ParameterStatus{}
	case 'Z':
		out = &message.ReadyForQuery{}
	case 'K':
		out = &message.BackendKeyData{}
	case 'T':
		out = &message.RowDescription{}
	case 'D':
		out = &message.DataRow{}
	case 'C':
		out = &message.CommandComplete{}
	case 'E':
		out = &message.ErrorResponse{}
	default:
		out = nil
	}

	var err error
	if out == nil {
		err = errorMsgNotSupported(in)
	}
	return out, err
}

//                                  length+4   peekInt32  data
// AuthenticationOk:                Int32(8),  Int32(0)
// AuthenticationKerberosV5:        Int32(8),  Int32(2)
// AuthenticationCleartextPassword: Int32(8),  Int32(3)
// AuthenticationGSS:               Int32(8),  Int32(7)
// AuthenticationSSPI:              Int32(8),  Int32(9)

// AuthenticationMD5Password:       Int32(12), Int32(5),  Byte4

// AuthenticationGSSContinue:       Int32,     Int32(8),  Byten
// AuthenticationSASL:              Int32,     Int32(10), String
// AuthenticationSASLContinue:      Int32,     Int32(11), Byten
// AuthenticationSASLFinal:         Int32,     Int32(12), Byten

func deduceAuth(in *fullMessage) (message.Response, error) {
	var out message.Response

	spec := in.body.PeekInt32()
	switch in.length {
	case 4:
		switch spec {
		case 0:
			out = &message.AuthOk{}
		case 2:
			out = &message.AuthKerberosV5{}
		case 3:
			out = &message.AuthCleartextPassword{}
		case 7:
			out = &message.AuthGSS{}
		case 9:
			out = &message.AuthSSPI{}
		}
	case 8:
		out = &message.AuthMD5Password{}
	default:
		switch spec {
		case 10:
			out = &message.AuthSASL{}
		}
	}

	var err error
	if out == nil {
		err = errorMsgNotSupported(in)
	}
	return out, err
}

func errorMsgNotSupported(in *fullMessage) error {
	return fmt.Errorf("[%x / '%c'] id is not supported", in.id, in.id)
}
