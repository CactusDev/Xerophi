package command

import (
	"net/http"

	"github.com/CactusDev/Xerophi/rethink"
	"github.com/CactusDev/Xerophi/util"

	"github.com/gin-gonic/gin"

	log "github.com/Sirupsen/logrus"
	mapstruct "github.com/mitchellh/mapstructure"
)

// Command is the struct that implements the handler interface for the command resource
type Command struct {
	Conn           *rethink.Connection // The RethinkDB connection
	Table          string              // The database table we're using
	CreationSchema struct{}            // The schema that the incoming JSON will be decoded into
	ResponseSchema struct{}            // The schema that will sent in response
}

// Update handles the updating of a record if the record exists, otherwise create it
func (c *Command) Update(ctx *gin.Context) {
	filter := map[string]interface{}{"token": ctx.Param("token"), "name": ctx.Param("name")}
	resp, err := c.Conn.GetByFilter(c.Table, filter, 1)
	if err != nil {
		log.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if resp == nil {
		// Command doesn't exist, pass control to c.Create
		c.Create(ctx)
	}
	// Command exists, update it
	ctx.Data(http.StatusTeapot, "text", []byte("Ohai"))
}

// GetAll returns all records associated with the token
func (c *Command) GetAll(ctx *gin.Context) {
	filter := map[string]interface{}{"token": ctx.Param("token")}
	resp, err := c.Conn.GetByFilter(c.Table, filter, 0)
	if err != nil {
		log.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if resp == nil {
		resp = make([]interface{}, 0)
	}
	ctx.JSON(http.StatusOK, gin.H{"data": resp})
}

// GetSingle returns a single record
func (c *Command) GetSingle(ctx *gin.Context) {
	filter := map[string]interface{}{"token": ctx.Param("token"), "name": ctx.Param("name")}
	resp, err := c.Conn.GetSingle(filter, c.Table)
	if err != nil {
		log.Error(err.Error())
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
	if resp == nil {
		resp = make([]interface{}, 0)
	}
	var response ResponseSchema
	err = mapstruct.Decode(resp, &response)
	if err != nil {
		log.Error(err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{"data": util.MarshalResponse(response)})
}

// Create creates a new record
func (c *Command) Create(ctx *gin.Context) {
	var vals map[string]interface{}
	ctx.BindJSON(&vals)
	resp, err := c.Conn.Create(c.Table, vals) // data)
	if err != nil {
		log.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	ctx.JSON(http.StatusCreated, resp)
}

// Delete removes a record
func (c *Command) Delete(ctx *gin.Context) {}
