package message

import "github.com/starnuik/naive-pgd/pkg/pgwire/internal"

type Response interface {
	Populate(body internal.ReadBuf) error
}

type Request interface {
	Bytes() []byte
}
