package message

import "github.com/starnuik/pogoda/pkg/pgwire/internal"

type Response interface {
	Populate(body internal.ReadBuf) error
}

type Request interface {
	Bytes() []byte
}
