package types

import (
	"github.com/CactusDev/Xerophi/rethink"
	"github.com/gin-gonic/gin"
)

// Meta is just a wrapper for map[string]interface{} to be used for JSONAPI meta
type Meta map[string]interface{}

// Handler is the collection of methods required for a type to be a handler
type Handler interface {
	Update(*gin.Context)
	GetAll(*gin.Context)
	GetSingle(*gin.Context)
	Create(*gin.Context)
	Delete(*gin.Context)
}

// DatabaseInfo keeps track of the information each handler requires
type DatabaseInfo struct {
	Table      string
	Connection *rethink.Connection
	Meta       map[string]interface{}
	Schema     map[string]interface{}
}
