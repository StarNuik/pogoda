package pgwire

import (
	"fmt"

	"github.com/starnuik/pogoda/pkg/pgwire/message"
)

func deduceResponseType(in *fullMessage) (message.Response, error) {
	var out message.Response
	switch in.id {
	case 'R':
		return deduceAuth(in)
	case 'S':
		out = &message.ParameterStatus{}
	case 'Z':
		out = &message.ReadyForQuery{}
	case 'K':
		out = &message.BackendKeyData{}
	case 'T':
		out = &message.RowDescription{}
	case 'D':
		out = &message.DataRow{}
	case 'C':
		out = &message.CommandComplete{}
	case 'E':
		out = &message.ErrorResponse{}
	default:
		out = nil
	}

	var err error
	if out == nil {
		err = errorNotSupported(in)
	}
	return out, err
}

//                                  length+4   peekInt32  data
// AuthenticationOk:                Int32(8),  Int32(0)
// AuthenticationKerberosV5:        Int32(8),  Int32(2)
// AuthenticationCleartextPassword: Int32(8),  Int32(3)
// AuthenticationGSS:               Int32(8),  Int32(7)
// AuthenticationSSPI:              Int32(8),  Int32(9)

// AuthenticationMD5Password:       Int32(12), Int32(5),  Byte4

// AuthenticationGSSContinue:       Int32,     Int32(8),  Byten
// AuthenticationSASL:              Int32,     Int32(10), String
// AuthenticationSASLContinue:      Int32,     Int32(11), Byten
// AuthenticationSASLFinal:         Int32,     Int32(12), Byten

func deduceAuth(in *fullMessage) (message.Response, error) {
	var out message.Response

	spec := in.body.PeekInt32()
	switch in.length {
	case 4:
		switch spec {
		case 0:
			out = &message.AuthOk{}
		case 2:
			out = &message.AuthKerberosV5{}
		case 3:
			out = &message.AuthCleartextPassword{}
		case 7:
			out = &message.AuthGSS{}
		case 9:
			out = &message.AuthSSPI{}
		}
	case 8:
		out = &message.AuthMD5Password{}
	default:
		switch spec {
		case 10:
			out = &message.AuthSASL{}
		}
	}

	var err error
	if out == nil {
		err = errorNotSupported(in)
	}
	return out, err
}

func errorNotSupported(in *fullMessage) error {
	return fmt.Errorf("[%x / '%c'] id is not supported", in.id, in.id)
}
