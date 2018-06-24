package types

import (
	"github.com/gin-gonic/gin"

	"github.com/CactusDev/Xerophi/rethink"
	//"github.com/CactusDev/Xerophi/secure"
)

// Meta is just a wrapper for map[string]interface{} to be used for JSONAPI meta
type Meta map[string]interface{}

// Schema is an interface for schemas so we can pass all different schemas
// for our parsing and validation functions
type Schema interface {
	DumpBody(data []byte) ([]byte, error)
}

// RouteDetails gives us the info needed to automatically create handlers
type RouteDetails struct {
	Enabled bool            // Whether the function is enabled
	Handler gin.HandlerFunc // Handler function
	Path    string          // /:name/delete - Extension path off of group base
	Verb    string          // GET/POST/PATCH/etc. - HTTP request verb
	Scopes  []string        // Scopes required to interact with this endpoint
}

// Handler is the collection of methods required for a type to be a handler
type Handler interface {
	Routes() []RouteDetails
}

// DatabaseInfo keeps track of the information each handler requires
type DatabaseInfo struct {
	Table      string
	Connection *rethink.Connection
	Meta       map[string]interface{}
	Schema     map[string]interface{}
}
