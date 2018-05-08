package types

import (
	"github.com/gin-gonic/gin"

	"github.com/CactusDev/Xerophi/rethink"
	//"github.com/CactusDev/Xerophi/secure"
)

// Meta is just a wrapper for map[string]interface{} to be used for JSONAPI meta
type Meta map[string]interface{}

// RouteDetails gives us the info needed to automatically create handlers
type RouteDetails struct {
	Enabled bool
	Handler gin.HandlerFunc
	Path    string
	Verb    string
	// Protected secure.AuthDetails	// Information on whether authentication is required
}

// Handler is the collection of methods required for a type to be a handler
type Handler interface {
	Update(*gin.Context)
	GetAll(*gin.Context)
	GetSingle(*gin.Context)
	Create(*gin.Context)
	Delete(*gin.Context)
	Routes() []RouteDetails
}

// DatabaseInfo keeps track of the information each handler requires
type DatabaseInfo struct {
	Table      string
	Connection *rethink.Connection
	Meta       map[string]interface{}
	Schema     map[string]interface{}
}
