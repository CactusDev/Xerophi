package command

import (
	"time"

	"github.com/CactusDev/Xerophi/schemas"
	"github.com/CactusDev/Xerophi/types"
	"github.com/CactusDev/Xerophi/util"
)

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
	Populated bool
}

// JSONAPIMeta returns a meta object for the response
func (rs ResponseSchema) JSONAPIMeta() *types.Meta {
	return &types.Meta{
		"createdAt": rs.CreatedAt,
		"token":     rs.Token,
	}
}

// ClientSchema is the schema the data from the client will be marshalled into
type ClientSchema struct {
	Arguments []schemas.MessagePacket `json:"arguments" binding:"required"`
	Enabled   bool                    `json:"enabled" binding:"omitempty"`
	Response  EmbeddedResponseSchema  `json:"response" binding:"required"`
}

// CreationSchema is all the data required for a new command to be created
type CreationSchema struct {
	ClientSchema `binding:"required"`
	// Ignore these fields in user input, they will be filled automatically by the API
	Count     int       `json:"count" binding:"-"`
	CreatedAt time.Time `json:"createdAt" binding:"-"`
	Token     string    `json:"token" binding:"-"`
	Name      string    `json:"name" binding:"-"`
}

// EmbeddedResponseSchema is the schema that is stored under the response key in ResponseSchema
type EmbeddedResponseSchema struct {
	Action  bool                    `json:"action" jsonapi:"attr,action" binding:"exists"`
	Message []schemas.MessagePacket `json:"message" jsonapi:"attr,message" binding:"required,gt=0"`
	Role    int                     `json:"role" jsonapi:"attr,role" binding:"gte=0,lte=256"`
	Target  string                  `json:"target" jsonapi:"attr,target"`
	User    string                  `json:"user" jsonapi:"attr,user"`
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (rs ResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(rs, lookup, "jsonapi")
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (r EmbeddedResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(r, lookup, "jsonapi")
}
