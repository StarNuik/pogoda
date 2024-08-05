package message

import "github.com/starnuik/naive-pgd/pkg/pgwire/internal"

// todo: panics -> errors

type ReadyForQuery struct {
	Status byte
}

func (m *ReadyForQuery) Populate(body internal.ReadBuf) error {
	// unsafe
	m.Status = body.Byte()
	return nil
}

type ParameterStatus struct {
	Key   string
	Value string
}

func (m *ParameterStatus) Populate(body internal.ReadBuf) error {
	// unsafe
	m.Key = body.String()
	m.Value = body.String()
	return nil
}

type BackendKeyData struct {
	Process int32
	Secret  int32
}

func (m *BackendKeyData) Populate(body internal.ReadBuf) error {
	// unsafe
	m.Process = body.Int32()
	m.Secret = body.Int32()
	return nil
}

type RowDescription struct {
	// FieldsCount int16
	Fields []RowField
}

type RowField struct {
	Name         string
	TableOid     int32
	ColumnNumber int16
	TypeOid      int32
	TypeSize     int16
	TypeModifier int32
	Format       int16
}

func (m *RowDescription) Populate(body internal.ReadBuf) error {
	len := body.Int16()
	for range len {
		m.Fields = append(m.Fields, RowField{
			Name:         body.String(),
			TableOid:     body.Int32(),
			ColumnNumber: body.Int16(),
			TypeOid:      body.Int32(),
			TypeSize:     body.Int16(),
			TypeModifier: body.Int32(),
			Format:       body.Int16(),
		})
	}
	return nil
}

type DataRow struct {
	// RowsCount int16
	Values []string
}

func (m *DataRow) Populate(body internal.ReadBuf) error {
	len := body.Int16()
	for range len {
		len := body.Int32()
		val := body.Bytes(len)
		// todo: force the connection to use 'text' format
		// todo: nil / SqlNull handling?
		m.Values = append(m.Values, string(val))
	}
	return nil
}

type CommandComplete struct {
	Tag string
}

func (m *CommandComplete) Populate(body internal.ReadBuf) error {
	m.Tag = body.String()
	return nil
}

type AuthOk struct{}
type AuthKerberosV5 struct{}
type AuthCleartextPassword struct{}
type AuthGSS struct{}
type AuthSSPI struct{}

func (*AuthOk) Populate(_ internal.ReadBuf) error                { return nil }
func (*AuthKerberosV5) Populate(_ internal.ReadBuf) error        { return nil }
func (*AuthCleartextPassword) Populate(_ internal.ReadBuf) error { return nil }
func (*AuthGSS) Populate(_ internal.ReadBuf) error               { return nil }
func (*AuthSSPI) Populate(_ internal.ReadBuf) error              { return nil }

type AuthMD5Password struct {
	Salt []byte
}

type AuthGSSContinue struct {
	Data []byte
}

type AuthSASL struct {
	Name string
}

type AuthSASLContinue struct {
	Data []byte
}

type AuthSASLFinal struct {
	Outcome []byte
}

func (m *AuthMD5Password) Populate(body internal.ReadBuf) error {
	m.Salt = body.Bytes(4)
	return nil
}

func (m *AuthGSSContinue) Populate(body internal.ReadBuf) error {
	m.Data = body.BytesRemainder()
	return nil
}

func (m *AuthSASL) Populate(body internal.ReadBuf) error {
	m.Name = body.String()
	return nil
}

func (m *AuthSASLContinue) Populate(body internal.ReadBuf) error {
	m.Data = body.BytesRemainder()
	return nil
}

func (m *AuthSASLFinal) Populate(body internal.ReadBuf) error {
	m.Outcome = body.BytesRemainder()
	return nil
}
