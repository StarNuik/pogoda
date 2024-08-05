package pgp

// todo: panics -> errors

type ReadyForQuery struct {
	Status byte
}

func (m *ReadyForQuery) populate(body readBuf) error {
	// unsafe
	m.Status = body.byte()
	return nil
}

type ParameterStatus struct {
	Key   string
	Value string
}

func (m *ParameterStatus) populate(body readBuf) error {
	// unsafe
	m.Key = body.string()
	m.Value = body.string()
	return nil
}

type BackendKeyData struct {
	Process int32
	Secret  int32
}

func (m *BackendKeyData) populate(body readBuf) error {
	// unsafe
	m.Process = body.int32()
	m.Secret = body.int32()
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

func (m *RowDescription) populate(body readBuf) error {
	len := body.int16()
	for range len {
		m.Fields = append(m.Fields, RowField{
			Name:         body.string(),
			TableOid:     body.int32(),
			ColumnNumber: body.int16(),
			TypeOid:      body.int32(),
			TypeSize:     body.int16(),
			TypeModifier: body.int32(),
			Format:       body.int16(),
		})
	}
	return nil
}

type DataRow struct {
	// RowsCount int16
	Values []string
}

func (m *DataRow) populate(body readBuf) error {
	len := body.int16()
	for range len {
		len := body.int32()
		val := body.bytes(len)
		// todo: force the connection to use 'text' format
		// todo: nil / SqlNull handling?
		m.Values = append(m.Values, string(val))
	}
	return nil
}

type CommandComplete struct {
	Tag string
}

func (m *CommandComplete) populate(body readBuf) error {
	m.Tag = body.string()
	return nil
}

type ErrorResponse struct {
	Fields []ErrorField
}

type ErrorField struct {
	Code        byte
	Description string
}

func (m *ErrorResponse) populate(body readBuf) error {
	for /* end := bytes.IndexByte(body, '\000'); end > 0; */ {
		field := ErrorField{}
		field.Code = body.byte()
		if field.Code == '\000' {
			break
		}
		field.Description = body.string()
		m.Fields = append(m.Fields, field)
	}
	return nil
}

type AuthOk struct{}
type AuthKerberosV5 struct{}
type AuthCleartextPassword struct{}
type AuthGSS struct{}
type AuthSSPI struct{}

func (*AuthOk) populate(_ readBuf) error                { return nil }
func (*AuthKerberosV5) populate(_ readBuf) error        { return nil }
func (*AuthCleartextPassword) populate(_ readBuf) error { return nil }
func (*AuthGSS) populate(_ readBuf) error               { return nil }
func (*AuthSSPI) populate(_ readBuf) error              { return nil }

type AuthMD5Password struct {
	Salt []byte
}

// type AuthGSSContinue struct {
// 	Data []byte
// }

type AuthSASL struct {
	Name string
}

// type AuthSASLContinue struct {
// 	Data []byte
// }

// type AuthSASLFinal struct {
// 	Outcome []byte
// }

func (m *AuthMD5Password) populate(body readBuf) error {
	m.Salt = body.bytes(4)
	return nil
}

// func (m *AuthGSSContinue) populate(body readBuf) error { return nil }
func (m *AuthSASL) populate(body readBuf) error {
	m.Name = body.string()
	return nil
}

// func (m *AuthSASLContinue) populate(body readBuf) error { return nil }
// func (m *AuthSASLFinal) populate(body readBuf) error    { return nil }
