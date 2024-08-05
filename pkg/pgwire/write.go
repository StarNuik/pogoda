package pgwire

import "github.com/starnuik/pogoda/pkg/pgwire/message"

func (c *Conn) Write(in message.Request) error {
	bytes := in.Bytes()
	_, err := c.w.Write(bytes)
	return err
}
