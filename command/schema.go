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
	CreatedAt string                  `jsonapi:"createdAt,meta"`
	Enabled   bool                    `jsonapi:"enabled,attr"`
	Name      string                  `jsonapi:"name,attr"`
	Response  EmbeddedResponseSchema  `jsonapi:"response,attr"`
	Token     string                  `jsonapi:"token,meta"`
}

// ClientSchema is the schema the data from the client will be marshalled into
type ClientSchema struct{}

// EmbeddedResponseSchema is the schema that is stored under the response key in ResponseSchema
type EmbeddedResponseSchema struct {
	Action  bool                    `jsonapi:"action,attr"`
	Message []schemas.MessagePacket `jsonapi:"message,attr"`
	Role    int                     `jsonapi:"role,attr"`
	Target  string                  `jsonapi:"target,attr"`
	User    string                  `jsonapi:"user,attr"`
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (r ResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(r, lookup, "jsonapi")
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (r EmbeddedResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(r, lookup, "jsonapi")
}
