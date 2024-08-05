package request_test

import (
	"encoding/hex"
	"testing"

	"github.com/starnuik/pogoda/pkg/pgwire/message/request"
	"github.com/stretchr/testify/assert"
)

func TestStartupEncode(t *testing.T) {
	assert := assert.New(t)

	tt := []struct {
		user     string
		database string
		wantHex  string
	}{
		{"meow", "meow", "000000210003000075736572006d656f77006461746162617365006d656f770000"},
	}

	for _, t := range tt {
		want, _ := hex.DecodeString(t.wantHex)

		req := request.Startup{
			User:     t.user,
			Database: t.database,
		}
		have := req.Encode()

		assert.Equal(want, have)
	}
}
