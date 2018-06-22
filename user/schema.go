package user

import (
	"encoding/json"

	"github.com/CactusDev/Xerophi/types"
	"github.com/CactusDev/Xerophi/util"
)

// Database is all the data stored in the DB for a user record
type Database struct {
	ID        string  // RethinkDB Record UUID
	Hash      string  // Argon2 password hash
	Salt      string  // Random salt used for hashing
	Token     string  // The user's token
	DeletedAt float64 // A unix timestamp value for soft deletion, 0 if active
	CreatedAt string  // A timestamp string for when the record was created
	UserID    int     // The internal static numeric ID for the user
	UserName  string  // The username for the user, can be changed unlike UserID
	Service   string  // The default service this user uses
}

// ResponseSchema is the schema for the data that will be sent out to the client
type ResponseSchema struct {
	ID string `jsonapi:"primary,user"`

	Service  string `jsonapi:"attr,service"`
	UserName string `jsonapi:"attr,username"`
	UserID   int    `jsonapi:"attr,userId"`

	Token     string `jsonapi:"meta,token"`
	CreatedAt string `jsonapi:"meta,createdAt"`
}

// ConfigSchema is the data returned for a user's config
type ConfigSchema struct {
	ID string `jsonapi:"primary,config"`

	Announce EventMessages `jsonapi:"attr,announce"`
	Services []Service     `jsonapi:"attr,services"`
	Spam     Spam          `jsonapi:"attr,spam"`

	// Meta
	Token     string `jsonapi:"meta,token"`
	CreatedAt string `jsonapi:"meta,createdAt"`
}

// EventMessages is all the event messages under the Announcements field
type EventMessages struct {
	Follow EventMessage `json:"follow,omitempty" jsonapi:"attr,follow"`
	Host   EventMessage `json:"host,omitempty" jsonapi:"attr,host"`
	Join   EventMessage `json:"join,omitempty" jsonapi:"attr,join"`
	Leave  EventMessage `json:"leave,omitempty" jsonapi:"attr,leave"`
	Sub    EventMessage `json:"sub,omitempty" jsonapi:"attr,sub"`
}

// EventMessage provides the configuration options for announcing different events
type EventMessage struct {
	// Darn Go being silly with having to make it a pointer to not ignore false
	Announce *bool  `json:"announce,omitempty" jsonapi:"attr,announce"`
	Message  string `json:"message,omitempty" jsonapi:"attr,message"`
}

// Service represents a service the bot is connecting to
type Service struct {
	IsOAuth     *bool    `json:"isOAuth" jsonapi:"attr,isOAuth"`
	Name        string   `json:"name" jsonapi:"attr,name"`
	Permissions []string `json:"permissions" jsonapi:"attr,permissions"`
	Username    string   `json:"username" jsonapi:"attr,username"`
}

// Spam is the data regarding spam handling for the bot
type Spam struct {
	AllowLinks      *bool      `json:"allowLinks,omitempty" jsonapi:"attr,allowLinks"`
	AutoTimeout     *bool      `json:"autoTimeout,omitempty" jsonapi:"attr,autoTimeout"`
	Blacklist       []string   `json:"blacklist,omitempty" jsonapi:"attr,blacklist"`
	Whitelist       []string   `json:"whitelist,omitempty" jsonapi:"attr,whitelist"`
	WhitelistedUrls []string   `json:"whitelistedUrls,omitempty" jsonapi:"attr,whitelistedUrls"`
	AllowUrls       SpamAction `json:"allowUrls,omitempty" jsonapi:"attr,allowUrls"`
	MaxEmoji        SpamAction `json:"maxEmoji,omitempty" jsonapi:"attr,maxEmoji"`
	MaxCapsScore    SpamAction `json:"maxCapsScore,omitempty" jsonapi:"attr,maxCapsScore"`
}

// SpamAction is a set of configuration options used for what do with a certain
// spam event
type SpamAction struct {
	Action   string `json:"action,omitempty" jsonapi:"attr,action"`
	Value    int    `json:"value,omitempty" jsonapi:"attr,value"`
	Warnings int    `json:"warnings,omitempty" jsonapi:"attr,warnings"`
}

// UpdateSchema is ClientSchema that is used when updating
type UpdateSchema struct {
	Announce EventMessages `json:"announce,omitempty" jsonapi:"attr,announce"`
	Services []Service     `json:"services,omitempty" jsonapi:"attr,services"`
	Spam     Spam          `json:"spam,omitempty" jsonapi:"attr,spam"`
}

// JSONAPIMeta returns a meta object for the response
func (rs ResponseSchema) JSONAPIMeta() *types.Meta {
	return &types.Meta{
		"createdAt": rs.CreatedAt,
		"token":     rs.Token,
	}
}

// JSONAPIMeta returns a meta object for the response
func (cs ConfigSchema) JSONAPIMeta() *types.Meta {
	return &types.Meta{
		"createdAt": cs.CreatedAt,
		"token":     cs.Token,
	}
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (rs ResponseSchema) GetAPITag(lookup string) string {
	return util.FieldTag(rs, lookup, "jsonapi")
}

// GetAPITag allows each of these types to implement the JSONAPISchema interface
func (cs ConfigSchema) GetAPITag(lookup string) string {
	return util.FieldTag(cs, lookup, "jsonapi")
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
