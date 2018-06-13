package rethink

import (
	"math/rand"

	log "github.com/sirupsen/logrus"
	r "gopkg.in/gorethink/gorethink.v4"
)

// RetrievalResult is a specific type of error object used for tracking
// the result of a database retrieval operation
type RetrievalResult struct {
	Success     bool
	SoftDeleted bool
	Message     string
}

func (rr RetrievalResult) Error() string {
	return rr.Message
}

// GetSingle returns a single object from the current table via a filter key
func (c *Connection) GetSingle(table string, filter Filter) (interface{}, error) {
	res, err := r.Table(table).Filter(filter).Run(c.Session)
	if err != nil {
		return nil, err
	}
	var response interface{}
	res.One(&response)

	// Probably want retrieval logic to be pure and not deal with any of these
	// checks
	// if response.(map[string]interface{})["deletedAt"].(float64) != 0 {
	// 	// Don't include anything that has a non-zero deletedAt (soft deleted)
	// 	return response, RetrievalResult{true, true, "Requested object is soft-deleted"}
	// }

	return response, nil
}

// GetByUUID returns a single object from the current table via the uuid
func (c *Connection) GetByUUID(uuid string, table string) (interface{}, error) {
	res, err := r.Table(table).Get(uuid).Run(c.Session)
	defer res.Close()
	if err != nil {
		return nil, err
	}
	var response interface{}
	res.One(&response)

	return response, nil
}

// GetAll returns all the record in a table, a wrapper around GetMultiple
func (c *Connection) GetAll(table string) ([]interface{}, error) {
	response, err := c.GetMultiple(table, 0) // 0 means all
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetMultiple returns multiple records from a table
func (c *Connection) GetMultiple(table string, limit int) ([]interface{}, error) {
	query := r.Table(table)
	if limit > 0 {
		query = r.Table(table).Limit(limit)
	}
	res, err := query.Run(c.Session)
	defer res.Close()
	if err != nil {
		return nil, err
	}
	var fromDB []interface{}
	var response []interface{}
	res.All(&fromDB)

	for _, val := range fromDB {
		if val.(map[string]interface{})["deletedAt"].(float64) != 0 {
			// Don't include anything that has a non-zero deletedAt (soft deleted)
			continue
		}
		response = append(response, val)
	}

	return response, nil
}

// GetByFilter is like GetMultiple, except it has the ability to filter the results first
func (c *Connection) GetByFilter(table string, filter Filter, limit int) ([]interface{}, error) {
	query := r.Table(table)
	if len(filter) > 0 {
		query = query.Filter(filter)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	res, err := query.Run(c.Session)
	defer res.Close()
	if err != nil {
		return nil, err
	}
	var fromDB []interface{}
	var response []interface{}
	res.All(&fromDB)

	for _, val := range fromDB {
		if val.(map[string]interface{})["deletedAt"].(float64) != 0 {
			// Don't include anything that has a non-zero deletedAt (soft deleted)
			continue
		}
		response = append(response, val)
	}

	return response, nil
}

// GetRandom retrieves a single random record from the table given the filter
func (c *Connection) GetRandom(table string, filter Filter) (interface{}, error) {
	response, err := c.GetByFilter(table, filter, 0)
	if err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, nil
	}

	return response[rand.Intn(len(response))], nil
}

// GetTotalRecords returns the total number of records in a table
func (c *Connection) GetTotalRecords(table string, filter Filter) (int, error) {
	res, err := r.Table(table).Filter(filter).Count().Run(c.Session)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	var response int
	res.One(&response)

	return response, nil
}

// IsSoftDeleted checks whether the supplied record's deletedAt key is non-zero
func IsSoftDeleted(record interface{}) bool {
	mapped, ok := record.(map[string]interface{})
	if !ok {
		// Not a map, so can't be soft-deleted ¯\_(ツ)_/¯
		return false
	}

	deletedAt, ok := mapped["deletedAt"]
	if ok {
		// deletedAt key exists, go ahead with check
		val, ok := deletedAt.(float64)
		// deletedAt was converted -> float64 and was non-zero => soft-deleted
		return ok && val != 0.0
	}

	// Didn't pass conversion tests
	return false
}
