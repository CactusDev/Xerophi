package command

import (
	"github.com/CactusDev/Xerophi/schemas"
	"github.com/CactusDev/Xerophi/util"
)

type meta map[string]interface{}

// ResponseSchema is the schema for the data that will be sent out to the client
type ResponseSchema struct {
	ID        string                  `jsonapi:"primary,command"`
	Arguments []schemas.MessagePacket `jsonapi:"attr,arguments"`
	Count     int                     `jsonapi:"attr,count"`
	CreatedAt string                  `jsonapi:"meta,createdAt"`
	Enabled   bool                    `jsonapi:"attr,enabled"`
	Name      string                  `jsonapi:"attr,name"`
	Response  EmbeddedResponseSchema  `jsonapi:"attr,response"`
	Token     string                  `jsonapi:"meta,token"`
}

func (rs ResponseSchema) JSONAPIMeta() *meta {
	return &meta{
		"createdAt": rs.CreatedAt,
		"token":     rs.Token,
	}
}

// ClientSchema is the schema the data from the client will be marshalled into
type ClientSchema struct {
	Arguments []schemas.MessagePacket `json:"arguments" validate:"required"`
	Enabled   bool                    `json:"enabled" validate:"required"`
	Response  EmbeddedResponseSchema  `json:"response" validate:"required"`
	// Ignore these fields in user input, they will be filled automatically by the API
	ID        string `json:"id" validate:"-"`
	Count     int    `json:"count" validate:"-"`
	CreatedAt string `json:"createdAt" validate:"-"`
	Token     string `json:"token" validate:"-"`
	Name      string `json:"name" validate:"-"`
}

// EmbeddedResponseSchema is the schema that is stored under the response key in ResponseSchema
type EmbeddedResponseSchema struct {
	Action  bool                    `jsonapi:"attr,action" validate:"required"`
	Message []schemas.MessagePacket `jsonapi:"attr,message" validate:"required,gt=0"`
	Role    int                     `jsonapi:"attr,role" validate:"gte=0,lte=256"`
	Target  string                  `jsonapi:"attr,target"`
	User    string                  `jsonapi:"attr,user"`
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (r ResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(r, lookup, "jsonapi")
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (r EmbeddedResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(r, lookup, "jsonapi")
}
