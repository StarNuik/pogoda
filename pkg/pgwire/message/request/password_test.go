package request_test

import (
	"encoding/hex"
	"testing"

	"github.com/starnuik/pogoda/pkg/pgwire/message/request"
	"github.com/stretchr/testify/assert"
)

func TestPasswordEncode(t *testing.T) {
	assert := assert.New(t)

	tt := []struct {
		from    string
		wantHex string
	}{
		{"password", "700000000d70617373776f726400"},
		{"", "700000000500"},
	}

	for _, t := range tt {
		want, _ := hex.DecodeString(t.wantHex)

		req := request.Password{
			Password: t.from,
		}
		have := req.Encode()

		assert.Equal(want, have)
	}
}
