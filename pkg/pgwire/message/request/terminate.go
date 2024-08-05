package request

import "github.com/starnuik/pogoda/pkg/pgwire/internal"

type Terminate struct{}

func (m *Terminate) Encode() []byte {
	out := internal.WriteBuf{}

	out.PrependHeader('X')
	return out
}
