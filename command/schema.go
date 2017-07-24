package command

import (
	"github.com/CactusDev/Xerophi/schemas"
	"github.com/CactusDev/Xerophi/util"
)

// ResponseSchema is the schema for the data that will be sent out to the client
type ResponseSchema struct {
	ID        string                  `jsonapi:"id,primary"`
	Arguments []schemas.MessagePacket `jsonapi:"arguments,attr"`
	Count     int                     `jsonapi:"count,attr"`
	CreatedAt string                  `jsonapi:"createdAt,attr"`
	Enabled   bool                    `jsonapi:"enabled,attr"`
	Name      string                  `jsonapi:"name,attr"`
	Response  EmbeddedResponseSchema  `jsonapi:"response,attr"`
	Token     string                  `jsonapi:"token,attr"`
}

// ClientSchema is the schema the data from the client will be marshalled into
type ClientSchema struct{}

// EmbeddedResponseSchema is the schema that is stored under the response key in ResponseSchema
type EmbeddedResponseSchema struct {
	Action  bool                    `jsonapi:"action"`
	Message []schemas.MessagePacket `jsonapi:"message"`
	Role    int                     `jsonapi:"role"`
	Target  string                  `jsonapi:"target"`
	User    string                  `jsonapi:"user"`
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (r ResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(r, lookup, "jsonapi")
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (r EmbeddedResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(r, lookup, "jsonapi")
}
