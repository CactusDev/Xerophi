package command

import "github.com/CactusDev/Xerophi/schemas"

// ResponseSchema is the schema for the data that will be sent out to the client
type ResponseSchema struct {
	ID        string                  `json:"id" jsonapi:"id,primary"`
	Arguments []schemas.MessagePacket `json:"arguments" jsonapi:"arguments,attributes"`
	Count     int                     `json:"count" jsonapi:"count,attributes"`
	CreatedAt string                  `json:"createdAt" jsonapi:"createdAt,attributes"`
	Enabled   bool                    `json:"enabled" jsonapi:"enabled,attributes"`
	Name      string                  `json:"name" jsonapi:"name,attributes"`
	Response  EmbeddedResponseSchema  `json:"response" jsonapi:"response,attributes"`
	Token     string                  `json:"token" jsonapi:"token,attributes"`
}

// ClientSchema is the schema the data from the client will be marshalled into
type ClientSchema struct{}

// EmbeddedResponseSchema is the schema that is stored under the response key in ResponseSchema
type EmbeddedResponseSchema struct {
	Action  bool                    `json:"action" jsonapi:"action"`
	Message []schemas.MessagePacket `json:"message"`
	Role    int                     `json:"role"`
	Target  string                  `json:"target"`
	User    string                  `json:"user"`
}
