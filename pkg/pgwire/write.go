package pgwire

import "github.com/starnuik/pogoda/pkg/pgwire/message"

func (c *Wire) Write(in message.Request) error {
	bytes := in.Encode()
	_, err := c.w.Write(bytes)
	return err
}
