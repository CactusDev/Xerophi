package rethink

import (
	r "gopkg.in/gorethink/gorethink.v4"
)

// GetSingle returns a single object from the current table via a filter key
func (c *Connection) GetSingle(filter map[string]interface{}, table string) (interface{}, error) {
	res, err := r.Table(table).Filter(filter).Run(c.Session)
	if err != nil {
		return nil, err
	}
	var response interface{}
	res.One(&response)

	if response == nil {
		return response, nil
	}

	if response.(map[string]interface{})["deletedAt"].(float64) != 0 {
		// Don't include anything that has a non-zero deletedAt (soft deleted)
		return nil, nil
	}

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

	if response.(map[string]interface{})["deletedAt"].(float64) != 0 {
		// Don't include anything that has a non-zero deletedAt (soft deleted)
		return nil, nil
	}

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

// GetByFilter is a like GetMultiple, except it has the ability to filter the results first
func (c *Connection) GetByFilter(table string, filter map[string]interface{}, limit int) ([]interface{}, error) {
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
