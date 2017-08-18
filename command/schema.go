package command

import (
	"github.com/CactusDev/Xerophi/schemas"
	"github.com/CactusDev/Xerophi/util"
)

// ResponseSchema is the schema for the data that will be sent out to the client
type ResponseSchema struct {
	ID        string          `jsonapi:"primary,id"`
	Name      string          `jsonapi:"attr,name"`
	Response  schemas.Message `jsonapi:"attr,response"`
	Count     int             `jsonapi:"attr,count"`
	Enabled   bool            `jsonapi:"attr,enabled"`
	CreatedAt string          `jsonapi:"meta,createdAt"`
	Token     string          `jsonapi:"meta,token"`
}

// ClientSchema is the schema the data from the client will be marshalled into
type ClientSchema struct {
	Context schemas.Context `json:"context" validate:"required"`
	Enabled bool            `json:"enabled" validate:"required"`
	// Ignore these fields in user input, they will be filled automatically by the API
	ID    string `json:"id" validate:"-"`
	Count int    `json:"count" validate:"-"`
	Token string `json:"token" validate:"-"`
	Name  string `json:"name" validate:"-"`
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (r ResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(r, lookup, "jsonapi")
}
