package trust

import (
	"encoding/json"
	"time"

	"github.com/CactusDev/Xerophi/types"
	"github.com/CactusDev/Xerophi/util"
)

// ResponseSchema is the schema for the data that will be sent out to the client
type ResponseSchema struct {
	ID string `jsonapi:"primary,trust"`

	Active   bool   `jsonapi:"attr,active"`
	UserID   string `jsonapi:"attr,userId"`
	UserName string `jsonapi:"attr,userName"`

	CreatedAt string `jsonapi:"meta,createdAt"`
	Token     string `jsonapi:"meta,token"`
}

// ClientSchema is the schema the data from the client will be marshalled into
type ClientSchema struct {
	UserName string `json:"userName"`
}

// CreationSchema is all the data required for a new trust to be created
type CreationSchema struct {
	ClientSchema
	// Ignore these fields in user input, they will be filled automatically by the API
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	DeletedAt float64   `json:"deletedAt"`
	Token     string    `json:"token"`
	UserID    string    `json:"userId"`
}

// UpdateSchema is ClientSchema that is used when updating
type UpdateSchema struct {
	// Has to be a pointer otherwise Go assumes a false value means it hasn't
	// been assigned and thus it's omitted
	Active   *bool  `json:"active,omitempty"`
	UserName string `json:"userName,omitempty"`
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
