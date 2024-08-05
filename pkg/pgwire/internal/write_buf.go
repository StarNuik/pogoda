package internal

import "encoding/binary"

//? https://github.com/lib/pq/blob/3d613208bca2e74f2a20e04126ed30bcb5c4cc27/buf.go
//? https://github.com/jackc/pgx/blob/a68e14fe5ad7caed8657816b9883ed418f3324ec/internal/pgio/write.go
type WriteBuf []byte

func (b *WriteBuf) AddInt32(in int32) {
	*b = binary.BigEndian.AppendUint32(*b, uint32(in))
}

func (b *WriteBuf) AddString(in string) {
	*b = append(*b, in...)
	*b = append(*b, '\000')
}

func (b *WriteBuf) PrependHeader(t byte) {
	len := len(*b) + 4

	slice := make([]byte, 5)
	slice[0] = t
	binary.BigEndian.PutUint32(slice[1:], uint32(len))
	*b = append(slice, *b...)
}

// request.Startup exclusive
func (b *WriteBuf) PrependLength() {
	len := len(*b) + 4

	slice := make([]byte, 4)
	binary.BigEndian.PutUint32(slice, uint32(len))
	*b = append(slice, *b...)
}
