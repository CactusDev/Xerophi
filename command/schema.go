package command

import (
	"encoding/json"
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
}

// ClientSchema is the schema the data from the client will be marshalled into
type ClientSchema struct {
	Arguments []schemas.MessagePacket `json:"arguments"`
	Enabled   *bool                   `json:"enabled"`
	Response  EmbeddedResponseSchema  `json:"response"`
}

// CreationSchema is all the data required for a new command to be created
type CreationSchema struct {
	ClientSchema
	// Ignore these fields in user input, they will be filled automatically by the API
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"createdAt"`
	DeletedAt float64   `json:"deletedAt"`
	Token     string    `json:"token"`
	Name      string    `json:"name"`
	Enabled   bool      `json:"enabled"`
}

// UpdateSchema is ClientSchema that is used when updating
type UpdateSchema struct {
	Arguments []schemas.MessagePacket      `json:"arguments,omitempty"`
	Enabled   *bool                        `json:"enabled,omitempty"`
	Response  UpdateEmbeddedResponseSchema `json:"response,omitempty"`
}

// EmbeddedResponseSchema is the schema that is stored under the response key in ResponseSchema
type EmbeddedResponseSchema struct {
	Action  bool                    `json:"action" jsonapi:"attr,action"`
	Message []schemas.MessagePacket `json:"message" jsonapi:"attr,message"`
	Role    int                     `json:"role" jsonapi:"attr,role"`
	Target  string                  `json:"target" jsonapi:"attr,target"`
	User    string                  `json:"user" jsonapi:"attr,user"`
}

// UpdateEmbeddedResponseSchema is the schema that is stored under the response key in UpdateSchema
type UpdateEmbeddedResponseSchema struct {
	Action  *bool                   `json:"action,omitempty" jsonapi:"attr,action"`
	Message []schemas.MessagePacket `json:"message,omitempty" jsonapi:"attr,message"`
	Role    *int                    `json:"role,omitempty" jsonapi:"attr,role"`
	Target  *string                 `json:"target,omitempty" jsonapi:"attr,target"`
	User    *string                 `json:"user,omitempty" jsonapi:"attr,user"`
}

// JSONAPIMeta returns a meta object for the response
func (rs ResponseSchema) JSONAPIMeta() *types.Meta {
	return &types.Meta{
		"createdAt": rs.CreatedAt,
		"token":     rs.Token,
	}
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (rs ResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(rs, lookup, "jsonapi")
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (r EmbeddedResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(r, lookup, "jsonapi")
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (r UpdateEmbeddedResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(r, lookup, "jsonapi")
}

// DumpBody dumps the body data bytes into this specific schema and returns
// the bytes from this
func (cs CreationSchema) DumpBody(data []byte) ([]byte, error) {
	// Unmarshal the byte slice into the provided schema
	if err := json.Unmarshal(data, &cs); err != nil {
		return nil, err
	}

	// Marshal the unmarshalled byte slice back into a byte array
	schemaBytes, err := json.Marshal(cs)
	if err != nil {
		return nil, err
	}

	return schemaBytes, nil
}

// DumpBody dumps the body data bytes into this specific schema and returns
// the bytes from this
func (rs ResponseSchema) DumpBody(data []byte) ([]byte, error) {
	// Unmarshal the byte slice into the provided schema
	if err := json.Unmarshal(data, &rs); err != nil {
		return nil, err
	}

	// Marshal the unmarshalled byte slice back into a byte array
	schemaBytes, err := json.Marshal(rs)
	if err != nil {
		return nil, err
	}

	return schemaBytes, nil
}

// DumpBody dumps the body data bytes into this specific schema and returns
// the bytes from this
func (cs ClientSchema) DumpBody(data []byte) ([]byte, error) {
	// Unmarshal the byte slice into the provided schema
	if err := json.Unmarshal(data, &cs); err != nil {
		return nil, err
	}

	// Marshal the unmarshalled byte slice back into a byte array
	schemaBytes, err := json.Marshal(cs)
	if err != nil {
		return nil, err
	}

	return schemaBytes, nil
}

// DumpBody dumps the body data bytes into this specific schema and returns
// the bytes from this
func (us UpdateSchema) DumpBody(data []byte) ([]byte, error) {
	// Unmarshal the byte slice into the provided schema
	if err := json.Unmarshal(data, &us); err != nil {
		return nil, err
	}

	// Marshal the unmarshalled byte slice back into a byte array
	schemaBytes, err := json.Marshal(us)
	if err != nil {
		return nil, err
	}

	return schemaBytes, nil
}
