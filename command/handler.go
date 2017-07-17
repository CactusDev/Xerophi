package command

import (
	"github.com/CactusDev/Xerophi/rethink"
	"github.com/gin-gonic/gin"
)

// Command is the struct that implements the handler interface for the command resource
type Command struct {
	Conn *rethink.Connection
}

// Update handles the updating of a record
func (c *Command) Update(ctx *gin.Context) {}

// GetAll returns all records associated with the token
func (c *Command) GetAll(ctx *gin.Context) {

}

// GetSingle returns a single record
func (c *Command) GetSingle(ctx *gin.Context) {}

// Create creates a new record after checking that the record does not already exist, if so it will pass control to c.Update
func (c *Command) Create(ctx *gin.Context) {}

// Delete removes a record
func (c *Command) Delete(ctx *gin.Context) {}
