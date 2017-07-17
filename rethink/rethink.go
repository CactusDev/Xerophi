package rethink

import r "gopkg.in/gorethink/gorethink.v3"

// ConnectionOpts is what we need to connect to a RethinkDB server
type ConnectionOpts struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Connection defines a connection to a RethinkDB instance
type Connection struct {
	DB      string         // The Rethink DB to use
	Opts    ConnectionOpts // Connection options for connecting to the Rethink server
	Session *r.Session     // The connected session
}

// Database is a set of methods that must be implemented for an object to implement the Database interface
// Makes it easier if we change databases in the future
type Database interface {
	Connect() error
	Close() error
	GetSingle(table string) (interface{}, error)
	GetMultiple(table string, limit int) ([]interface{}, error)
	GetAll(table string) ([]interface{}, error)
	GetByUUID(uuid string, table string) (interface{}, error)
	GetByFilter(table string, filter map[string]interface{}, limit int) ([]interface{}, error)
	Monitor(table string) (*r.Cursor, error)
	FilteredMonitor(table string, filter map[string]interface{}) (*r.Cursor, error)
}

// Connect connects you to Rethink
func (c *Connection) Connect() error {
	session, err := r.Connect(r.ConnectOpts{
		Address:  c.Opts.Host,
		Database: c.DB,
		Username: c.Opts.User,
		Password: c.Opts.Password,
	})
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
