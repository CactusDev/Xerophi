package command

import (
	"net/http"

	"github.com/CactusDev/Xerophi/rethink"

	"github.com/gin-gonic/gin"

	log "github.com/Sirupsen/logrus"
)

// Command is the struct that implements the handler interface for the command resource
type Command struct {
	Conn           *rethink.Connection // The RethinkDB connection
	Table          string              // The database table we're using
	CreationSchema struct{}            // The schema that the incoming JSON will be decoded into
	ResponseSchema struct{}            // The schema that will sent in response
}

// Update handles the updating of a record
func (c *Command) Update(ctx *gin.Context) {}

// GetAll returns all records associated with the token
func (c *Command) GetAll(ctx *gin.Context) {
	log.Debug("Ohai")
	resp, err := c.Conn.GetAll(c.Table)
	if err != nil {
		log.Error(err.Error())
	}
	log.Debug("getall")
	ctx.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

// GetSingle returns a single record
func (c *Command) GetSingle(ctx *gin.Context) {}

// Create creates a new record after checking that the record does not already exist, if so it will pass control to c.Update
func (c *Command) Create(ctx *gin.Context) {}

// Delete removes a record
func (c *Command) Delete(ctx *gin.Context) {}
