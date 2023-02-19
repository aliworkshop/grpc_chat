package client

func (c *client) closeHandler(code int, text string) error {
	c.Stop()
	return nil
}
