package request

import "github.com/starnuik/pogoda/pkg/pgwire/internal"

type Query struct {
	Query string
}

func (m *Query) Encode() []byte {
	out := internal.WriteBuf{}

	out.AddString(m.Query)

	out.PrependHeader('Q')
	return out
}
