package rethink

import (
	r "gopkg.in/gorethink/gorethink.v4"
)

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
	GetByUUID(table string, uid string) (interface{}, error)
	GetByFilter(table string, filter map[string]interface{}, limit int) ([]interface{}, error)
	GetRandom(table string, filter map[string]interface{}) (interface{}, error)
	Update(table string, uid string, data map[string]interface{}) (interface{}, error)
	Create(table string, data map[string]interface{}) (interface{}, error)
	Delete(table string, uid string) (interface{}, error)  // Hard deletion
	Disable(table string, uid string) (interface{}, error) // Soft deletion
	Exists(table string, filter map[string]interface{}) (interface{}, error)
	Status(table string) (interface{}, error)
}

// Issue is the schema for any responses from RethinkDB will be in
// for any active issues
type Issue struct {
	ID          string
	Type        string
	Critical    bool
	Info        map[string]interface{}
	Description string
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

// Status returns status of the table specified
func (c *Connection) Status() ([]Issue, error) {
	var issues []Issue

	// Retrieve everything from the current issues admin table
	res, err := r.DB("rethinkdb").Table("current_issues").Run(c.Session)
	if err != nil {
		return nil, err
	}
	res.All(&issues)

	return issues, nil
}
