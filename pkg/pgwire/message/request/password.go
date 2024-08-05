package request

import "github.com/starnuik/pogoda/pkg/pgwire/internal"

type Password struct {
	Password string
}

func (m *Password) Encode() []byte {
	out := internal.WriteBuf{}

	out.AddString(m.Password)

	out.PrependHeader('p')
	return out
}
