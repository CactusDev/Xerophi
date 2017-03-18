package rethink

import rethink "gopkg.in/gorethink/gorethink.v3"

// Connection defines a connection to a RethinkDB table within a database
type Connection struct {
	DB      string
	Table   string
	URL     string
	Session *rethink.Session
}

// Connect connects you to Rethink
func (c *Connection) Connect() error {
	connArgs := rethink.ConnectOpts{
		Address:  c.URL,
		Database: c.DB,
	}

	session, err := rethink.Connect(connArgs)
	if err != nil {
		return err
	}
	c.Session = session

	return nil
}

// Close the current session
func (c *Connection) Close() error {
	err := c.Session.Close()
	if err != nil {
		return err
	}
	c.Session = nil

	return nil
}

// GetSingle returns a single object from the current table via a filter key
func (c *Connection) GetSingle(field string, value interface{}) (interface{}, error) {
	filter := make(map[string]interface{})
	filter[field] = value

	res, err := rethink.Table(c.Table).Filter(filter).Run(c.Session)
	if err != nil {
		return nil, err
	}
	var response interface{}
	res.One(&response)

	return response, nil
}

// GetByUUID returns a single object from the current table via the GetByUUID
func (c *Connection) GetByUUID(uuid string) (interface{}, error) {
	res, err := rethink.Table(c.Table).Get(uuid).Run(c.Session)
	defer res.Close()
	if err != nil {
		return nil, err
	}
	var response interface{}
	res.One(&response)

	return response, nil
}

// GetTable returns all the rescords in a table
func (c *Connection) GetTable() ([]interface{}, error) {
	res, err := rethink.Table(c.Table).Run(c.Session)
	defer res.Close()
	if err != nil {
		return nil, err
	}
	var response []interface{}
	res.All(&response)
	return response, nil
}
