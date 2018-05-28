package redis

// Exists takes a string and returns whether or not the key exists
func (c *Connection) Exists(key string) (bool, error) {
	res, err := c.Session.Get(key).Result()
	if res == "" {
		return false, err
	}
	return true, err
}

// Retrieve
