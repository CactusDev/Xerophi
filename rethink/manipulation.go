package rethink

import (
	"time"

	log "github.com/Sirupsen/logrus"
	r "gopkg.in/gorethink/gorethink.v3"
)

// Update takes the table the record is in, the UUID of the record, and the data to update it with - then updates the record (°Д°）
func (c *Connection) Update(table string, uid string, data map[string]interface{}) (interface{}, error) {
	resp, err := r.Table(table).Update(data).RunWrite(c.Session)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp, nil
}

// Create takes the table the record is in and the data to update it with, and creates a new record @(・0・)@
func (c *Connection) Create(table string, data map[string]interface{}) (interface{}, error) {
	resp, err := r.Table(table).Insert(data).RunWrite(c.Session)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp, nil
}

// Delete ... well, it deletes a record. Softly.
func (c *Connection) Delete(table string, uid string) (interface{}, error) {
	// Check if the record exists
	resp, err := r.Table(table).Update(map[string]interface{}{"deletedAt": time.Now().UTC().Unix()}).RunWrite(c.Session)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp, nil
}
