package client

func (c *client) SetTemp(key string, value interface{}) {
	c.valuesMtx.Lock()
	defer c.valuesMtx.Unlock()

	c.values[key] = value
}

func (c *client) GetTemp(key string) interface{} {
	c.valuesMtx.RLock()
	defer c.valuesMtx.RUnlock()

	return c.values[key]
}
