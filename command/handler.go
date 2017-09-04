package command

import (
	"encoding/json"
	"html"
	"net/http"
	"strings"

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

// Update handles the updating of a record if the record exists
func (c *Command) Update(ctx *gin.Context) {
	token := html.EscapeString(ctx.Param("token"))
	name := html.EscapeString(ctx.Param("name"))
	filter := map[string]interface{}{"token": token, "name": name}
	resp, err := c.Conn.GetByFilter(c.Table, filter, 1)

	if err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
		return
	}
	if resp == nil {
		// Resource doesn't exist, return a 404
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}
	// Command exists, lets update it
	// Bind the JSON from the request
	var updateData map[string]interface{}
	err = ctx.BindJSON(updateData)
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}
	// What ID are we using to update?
	id, err := util.GetResourceID(resp)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	updated, err := c.Conn.Update(c.Table, id, updateData)
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"original": resp, "updated": updated})
}

// GetAll returns all records associated with the token
func (c *Command) GetAll(ctx *gin.Context) {
	token := html.EscapeString(ctx.Param("token"))
	filter := map[string]interface{}{"token": token}
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
	token := html.EscapeString(ctx.Param("token"))
	name := html.EscapeString(ctx.Param("name"))
	filter := map[string]interface{}{"token": token, "name": name}
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
		switch err.(type) {
		case *json.UnmarshalTypeError:
			ve, ok := err.(*json.UnmarshalTypeError)
			if !ok {
				log.Error("SPLODEY NOT OKAY!")
			}
			log.Warn("ve:\t", ve.Field)
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err)
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
