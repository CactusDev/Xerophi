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
		util.NiceError(ctx, err, http.StatusBadRequest)
	}
	if resp == nil {
		// Command doesn't exist, pass control to c.Create
		c.Create(ctx)
	}
	// Command exists, update it
	// ctx.JSON(http.StatusOK, gin.H{"data"})
}

// GetAll returns all records associated with the token
func (c *Command) GetAll(ctx *gin.Context) {
	filter := map[string]interface{}{"token": ctx.Param("token")}
	fromDB, err := c.Conn.GetByFilter(c.Table, filter, 0)
	if err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
		return
	}
	if fromDB == nil {
		ctx.JSON(http.StatusNotFound, make([]struct{}, 0))
		return
	}

	var respDecode ResponseSchema
	response := make([]map[string]interface{}, len(fromDB))
	var decoded []ResponseSchema
	for pos, record := range fromDB {
		// If there's an issue decoding it, just log it and move on to the next record
		if err := mapstruct.Decode(record, &respDecode); err != nil {
			log.Error(err.Error())
			continue
		}
		decoded[pos] = respDecode
	}
	// foo := util.MarshalResponse(decoded...)

	ctx.JSON(http.StatusOK, response)
}

// GetSingle returns a single record
func (c *Command) GetSingle(ctx *gin.Context) {
	filter := map[string]interface{}{"token": ctx.Param("token"), "name": ctx.Param("name")}
	fromDB, err := c.Conn.GetSingle(filter, c.Table)
	if err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
		return
	}
	if fromDB == nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var response ResponseSchema
	if err = mapstruct.Decode(fromDB, &response); err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, util.MarshalResponse(response))
}

// Create creates a new record
func (c *Command) Create(ctx *gin.Context) {
	var vals map[string]interface{}
	ctx.BindJSON(&vals)
	vals["token"] = ctx.Param("token")
	vals["name"] = ctx.Param("name")

	record, err := c.Conn.Create(c.Table, vals)
	if err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
		return
	}

	ctx.JSON(200, record)
}

// Delete removes a record
func (c *Command) Delete(ctx *gin.Context) {}
