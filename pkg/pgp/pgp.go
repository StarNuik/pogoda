package pgp

import (
	"encoding/binary"
	"fmt"
	"io"
)

type messageHeader struct {
	id     byte
	length int32
}

type fullMessage struct {
	messageHeader
	body readBuf
}

// WILL panic on len(head) < 5
func toHeader(bytes []byte) messageHeader {
	h := messageHeader{}
	h.id = bytes[0]
	h.length = int32(binary.BigEndian.Uint32(bytes[1:])) - 4
	return h
}

type ResponseMessage interface {
	populate(body readBuf) error
}

func Read(reader io.Reader) (ResponseMessage, error) {
	raw, err := read(reader)
	if err != nil {
		return nil, err
	}

	msg, err := deduceResponseType(raw)
	if err != nil {
		return nil, err
	}

	if len(raw.body) != int(raw.length) {
		return nil, fmt.Errorf("invalid message length")
	}

	err = msg.populate(raw.body)
	if err != nil {
		return nil, err
	}
	return msg, nil
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

func errorMsgNotSupported(in *fullMessage) error {
	return fmt.Errorf("[%x / '%c'] id is not supported", in.id, in.id)
}

// todo: this is only valid for receiving messages
func deduceResponseType(in *fullMessage) (ResponseMessage, error) {
	var out ResponseMessage
	switch in.id {
	case 'R':
		return deduceAuth(in)
	case 'S':
		out = &ParameterStatus{}
	case 'Z':
		out = &ReadyForQuery{}
	case 'K':
		out = &BackendKeyData{}
	case 'T':
		out = &RowDescription{}
	case 'D':
		out = &DataRow{}
	case 'C':
		out = &CommandComplete{}
	case 'E':
		out = &ErrorResponse{}
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

func deduceAuth(in *fullMessage) (ResponseMessage, error) {
	var out ResponseMessage

	spec := in.body.peekInt32()
	switch in.length {
	case 4:
		switch spec {
		case 0:
			out = &AuthOk{}
		case 2:
			out = &AuthKerberosV5{}
		case 3:
			out = &AuthCleartextPassword{}
		case 7:
			out = &AuthGSS{}
		case 9:
			out = &AuthSSPI{}
		}
	case 8:
		out = &AuthMD5Password{}
	default:
		switch spec {
		case 10:
			out = &AuthSASL{}
		}
	}

	var err error
	if out == nil {
		err = errorMsgNotSupported(in)
	}
	return out, err
}
