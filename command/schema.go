package command

import (
	"github.com/CactusDev/Xerophi/schemas"
	"github.com/CactusDev/Xerophi/util"
)

// ResponseSchema is the schema for the data that will be sent out to the client
type ResponseSchema struct {
	ID        string                  `jsonapi:"id,primary"`
	Arguments []schemas.MessagePacket `jsonapi:"arguments,attributes"`
	Count     int                     `jsonapi:"count,attributes"`
	CreatedAt string                  `jsonapi:"createdAt,attributes"`
	Enabled   bool                    `jsonapi:"enabled,attributes"`
	Name      string                  `jsonapi:"name,attributes"`
	Response  EmbeddedResponseSchema  `jsonapi:"response,attributes"`
	Token     string                  `jsonapi:"token,attributes"`
}

// ClientSchema is the schema the data from the client will be marshalled into
type ClientSchema struct {
}

// EmbeddedResponseSchema is the schema that is stored under the response key in ResponseSchema
type EmbeddedResponseSchema struct {
	Action bool `jsonapi:"action"`
}

// GetAPITag returns the jsonapi tag value for the specified field
func (r EmbeddedResponseSchema) GetAPITag(lookup string) string {
	return util.GetAPITag(r, lookup)
}

// GetAPITag returns the jsonapi tag value for the specified field
func (r ResponseSchema) GetAPITag(lookup string) string {
	return util.GetAPITag(r, lookup)
}
