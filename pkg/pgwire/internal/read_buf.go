package internal

import (
	"bytes"
	"encoding/binary"
)

//? https://github.com/lib/pq/blob/3d613208bca2e74f2a20e04126ed30bcb5c4cc27/buf.go
// todo: comb lib/pq for why they don't do errors on readBuf's .int()/.string()

type ReadBuf []byte

func (b *ReadBuf) Byte() byte {
	out := (*b)[0]
	*b = (*b)[1:]
	return out
}

func (b *ReadBuf) String() string {
	end := bytes.IndexByte(*b, '\000')
	if end < 0 {
		panic("string terminator not found")
	}
	out := (*b)[:end]
	*b = (*b)[end+1:]
	return string(out)
}

func (b *ReadBuf) Int32() int32 {
	// if len(*b) < 4 {
	// 	return -1, fmt.Errorf("invalid message: couldn't scan int32")
	// }
	out := int32(binary.BigEndian.Uint32(*b))
	*b = (*b)[4:]
	return out /* , nil */
}

// N.B: this is actually an unsigned 16-bit integer, unlike int32
// todo: ^ why ??
func (b *ReadBuf) Int16() int16 {
	out := int16(binary.BigEndian.Uint16(*b))
	*b = (*b)[2:]
	return out
}

func (b *ReadBuf) Bytes(count int32) []byte {
	// todo: will golang panic by itself?
	if count <= 0 {
		panic("readbuf: count <= 0")
	}
	out := (*b)[:count]
	*b = (*b)[count:]
	return out
}

func (b *ReadBuf) BytesRemainder() []byte {
	out := *b
	*b = ReadBuf{}
	return out
}

func (b *ReadBuf) PeekInt32() int32 {
	return int32(binary.BigEndian.Uint32(*b))
}
