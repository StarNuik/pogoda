package message

import (
	"fmt"

	"github.com/starnuik/naive-pgd/pkg/pgwire/internal"
)

type ErrorResponse struct {
	// Severity string
	// Code string
	// Message string
	Fields map[byte]string
}

// type ErrorField struct {
// 	Description string
// }

func (m *ErrorResponse) Populate(body internal.ReadBuf) error {
	m.Fields = make(map[byte]string)
	for /* end := bytes.IndexByte(body, '\000'); end > 0; */ {
		key := body.Byte()
		if key == '\000' {
			break
		}
		field := body.String()
		m.Fields[key] = field
	}
	return nil
}

func (m *ErrorResponse) Error() string {
	severity := m.Fields['V']
	code := m.Fields['C']
	message := m.Fields['M']
	return fmt.Sprintf("%s (%s): %s", severity, code, message)
}
