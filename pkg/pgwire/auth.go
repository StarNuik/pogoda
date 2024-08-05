package pgwire

import (
	"fmt"

	"github.com/starnuik/naive-pgd/pkg/pgwire/message"
)

func (c *Conn) Auth(user string, password string, database string) error {
	req := startupRequest(user, database)
	err := c.Write(&req)
	if err != nil {
		return err
	}

	res, err := c.Read()
	if err != nil {
		return err
	}

	auth := c.deduceAuthFlow(res)
	err = auth(password)
	if err != nil {
		return err
	}

	for {
		res, err := c.Read()
		if err != nil {
			return err
		}
		if _, is := res.(*message.ReadyForQuery); is {
			break
		}
	}

	return nil
}

func startupRequest(user string, database string) message.Startup {
	out := message.Startup{}
	out.User = user
	if database != "" {
		out.Database = database
	} else {
		out.Database = user
	}
	return out
}

func (c *Conn) deduceAuthFlow(in message.Response) func(string) error {
	var out func(string) error

	switch in.(type) {
	case *message.AuthOk:
		// no-op
		out = func(string) error { return nil }
	case *message.AuthCleartextPassword:
		out = c.authCleartext
	default:
		out = nil
	}
	return out
}

func (c *Conn) authCleartext(pass string) error {
	req := message.Password{
		Password: pass,
	}
	err := c.Write(&req)
	if err != nil {
		return err
	}

	res, err := c.Read()
	if err != nil {
		return err
	}
	switch res.(type) {
	case *message.AuthOk:
		return nil
	default:
		return fmt.Errorf("auth: unexpected message %#v", res)
	}
}
