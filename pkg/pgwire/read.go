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
func (c *Wire) Read() (message.Response, error) {
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
