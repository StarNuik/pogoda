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

// WILL panic on len(head) < 5
func (h *messageHeader) populate(headBytes []byte) {
	h.id = headBytes[0]
	h.length = int32(binary.BigEndian.Uint32(headBytes[1:])) - 4
}

type Message interface {
	populate(body readBuf) error
}

type AuthenticationOk struct{}

func Read(reader io.Reader) (Message, error) {
	headBytes, err := read(reader, 5)
	if err != nil {
		return nil, err
	}

	head := messageHeader{}
	head.populate(headBytes)

	msg, err := chooseMessage(head.id)
	if err != nil {
		discard(reader, head.length)
		return nil, err
	}

	bodyBytes, err := read(reader, head.length)
	if err != nil {
		return nil, err
	}

	err = msg.populate(bodyBytes)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func read(reader io.Reader, length int32) ([]byte, error) {
	buf := make([]byte, length)
	_, err := io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// todo: this is only valid for receiving messages
func chooseMessage(msgId byte) (Message, error) {
	var msg Message
	switch msgId {
	case 'S':
		msg = &ParameterStatus{}
	case 'Z':
		msg = &ReadyForQuery{}
	case 'K':
		msg = &BackendKeyData{}
	case 'T':
		msg = &RowDescription{}
	case 'D':
		msg = &DataRow{}
	case 'C':
		msg = &CommandComplete{}
	case 'E':
		msg = &ErrorResponse{}
	default:
		return nil, fmt.Errorf("[%x / '%c'] id is not supported", msgId, msgId)
	}
	return msg, nil
}

func discard(reader io.Reader, length int32) {
	io.CopyN(io.Discard, reader, int64(length))
	// bytes, err := read(reader, length)
	// if err != nil {
	// 	return
	// }
	// fmt.Println(hex.Dump(bytes))
}
