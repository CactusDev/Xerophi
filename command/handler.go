package command

import (
	"encoding/json"
	"net/http"
	"strings"

	"gopkg.in/go-playground/validator.v9"

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
		return
	}
	if resp == nil {
		// Command doesn't exist, pass control to c.Create
		c.Create(ctx)
		return
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
	var decoded = make([]map[string]interface{}, len(fromDB))
	for pos, record := range fromDB {
		// If there's an issue decoding it, just log it and move on to the next record
		if err := mapstruct.Decode(record, &respDecode); err != nil {
			log.Error(err.Error())
			continue
		}
		marshalled := util.MarshalResponse(respDecode)
		decoded[pos] = map[string]interface{}{
			"id":         marshalled["data"].(map[string]interface{})["id"],
			"attributes": marshalled["data"].(map[string]interface{})["attributes"],
			"meta":       marshalled["meta"],
		}
	}
	var response = make(map[string]interface{})

	response["data"] = decoded

	ctx.JSON(http.StatusOK, response)
}

// GetSingle returns a single record
func (c *Command) GetSingle(ctx *gin.Context) {
	filter := map[string]interface{}{"token": strings.ToLower(ctx.Param("token")), "name": ctx.Param("name")}
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
	var vals ClientSchema

	err := ctx.Bind(&vals)

	if err != nil {
		ve, ok := err.(validator.ValidationErrors)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		var errors = make([]map[string]interface{}, len(ve))
		for pos, vErr := range ve {
			errors[pos] = map[string]interface{}{
				vErr.Namespace(): vErr.Value(),
			}
		}
		errResp := map[string]interface{}{
			"errors": errors,
		}
		ctx.JSON(http.StatusBadRequest, errResp)
		// util.NiceError(ctx, err, http.StatusBadRequest)
		return
	}
	vals.Token = strings.ToLower(ctx.Param("token"))
	vals.Name = ctx.Param("name")

	log.Debugf("%+v", vals)

	var toCreate map[string]interface{}

	if data, err := json.Marshal(vals); err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	} else {
		json.Unmarshal(data, &toCreate)
	}

	log.Debugf("%+v", toCreate)

	// if _, err := c.Conn.Create(c.Table, toCreate); err != nil {
	// 	util.NiceError(ctx, err, http.StatusBadRequest)
	// 	return
	// }

	// Pass control off to GetSingle since we don't want to duplicate logic
	c.GetSingle(ctx)
}

// Delete removes a record
func (c *Command) Delete(ctx *gin.Context) {}
