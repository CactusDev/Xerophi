package quote

import (
	"time"

	"github.com/CactusDev/Xerophi/types"
	"github.com/CactusDev/Xerophi/util"
)

// ResponseSchema is the schema for the data that will be sent out to the client
type ResponseSchema struct {
	ID        string `jsonapi:"primary,command"`
	CreatedAt string `jsonapi:"meta,createdAt"`
	Enabled   bool   `jsonapi:"attr,enabled"`
	QuoteID   int    `jsonapi:"attr,quoteId"`
	Quote     string `jsonapi:"attr,quote"`
	Token     string `jsonapi:"meta,token"`
}

// ClientSchema is the schema the data from the client will be marshalled into
type ClientSchema struct {
	Quote string `json:"quote"`
}

// CreationSchema is all the data required for a new quote to be created
type CreationSchema struct {
	ClientSchema
	// Ignore these fields in user input, they will be filled automatically by the API
	CreatedAt time.Time `json:"createdAt"`
	DeletedAt float64   `json:"deletedAt"`
	Token     string    `json:"token"`
	QuoteID   int       `json:"quoteId"`
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (cs CreationSchema) GetAPITag(lookup string) string {
	return util.FieldTag(cs, lookup, "jsonapi")
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (rs ResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(rs, lookup, "jsonapi")
}

// JSONAPIMeta returns a meta object for the response
func (rs ResponseSchema) JSONAPIMeta() *types.Meta {
	return &types.Meta{
		"createdAt": rs.CreatedAt,
		"token":     rs.Token,
	}
}
