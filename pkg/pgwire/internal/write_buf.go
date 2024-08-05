package internal

import "encoding/binary"

type WriteBuf []byte

func (b *WriteBuf) Int32(in int32) {
	*b = binary.BigEndian.AppendUint32(*b, uint32(in))
}

func (b *WriteBuf) String(in string) {
	nullTerm := append([]byte(in), '\000')
	*b = append(*b, nullTerm...)
}

// todo: is this too concrete?
func (b *WriteBuf) PrependHeader(t byte) {
	len := len(*b) + 4

	slice := make([]byte, 5)
	slice[0] = t
	binary.BigEndian.PutUint32(slice[1:], uint32(len))
	*b = append(slice, *b...)
}

// StartupMessage exclusive
func (b *WriteBuf) PrependLength() {
	len := len(*b) + 4

	slice := make([]byte, 4)
	binary.BigEndian.PutUint32(slice, uint32(len))
	*b = append(slice, *b...)
}
