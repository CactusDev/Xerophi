package command

import (
	"encoding/json"
	"html"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/CactusDev/Xerophi/rethink"
	"github.com/CactusDev/Xerophi/util"

	"github.com/gin-gonic/gin"

	mapstruct "github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

var schema = "command/commandSchema.json"

// Command is the struct that implements the handler interface for the command resource
type Command struct {
	Conn           *rethink.Connection // The RethinkDB connection
	Table          string              // The database table we're using
	CreationSchema struct{}            // The schema that the incoming JSON will be decoded into
	ResponseSchema struct{}            // The schema that will sent in response
}

// ReturnOne retrieves a single record given the filter provided
func (c *Command) ReturnOne(filter map[string]interface{}) (ResponseSchema, error) {
	var response ResponseSchema

	// Retrieve a single record from the DB based on the filter
	fromDB, err := c.Conn.GetSingle(filter, c.Table)
	if err != nil {
		return response, err
	}
	// Was anything returned?
	if fromDB == nil {
		// Return nothing, it's not an error but there's nothing there
		return response, nil
	}

	// Decode the response from the DB into the response schema object
	if err = mapstruct.Decode(fromDB, &response); err != nil {
		return response, err
	}

	// It has been successfully populated, set this to true
	response.Populated = true
	return response, nil
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

	res, err := c.ReturnOne(filter)
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	var status int
	if res.Populated {
		status = http.StatusOK
	} else {
		status = http.StatusNotFound
	}

	ctx.JSON(status, util.MarshalResponse(res))
}

// Create creates a new record
func (c *Command) Create(ctx *gin.Context) {
	// Declare default values
	vals := ClientSchema{
		Enabled: true,
	}
	createVals := CreationSchema{
		CreatedAt: time.Now().UTC(),
	}

	// Update the token and name values from the request
	createVals.Token = strings.ToLower(ctx.Param("token"))
	createVals.Name = ctx.Param("name")
	// Get the JSON values
	createVals.ClientSchema = vals

	// Check if it exists yet
	filter := map[string]interface{}{"token": createVals.Token, "name": createVals.Name}
	// if res, _ := c.ReturnOne(filter); res.Populated {
	// 	// It exists already, error out, can't edit from this endpoint
	// 	ctx.AbortWithStatusJSON(http.StatusConflict, util.MarshalResponse(res))
	// 	return
	// }

	// Validate the data provided
	str, _ := ioutil.ReadAll(ctx.Request.Body)
	// Bind the JSON values in the request to the ClientSchema object
	err := ctx.BindJSON(&vals)
	// We can continue validating the input so as to get both syntax and schema
	// errors
	// Validate the JSON values binding
	if err != nil {
		// Check if it's a validation error and if so convert to a human-readable
		// error message
		resp, convErr := util.ValidateInput(string(str), schema)
		// convErr is nil so we have an error that needs to be returned to the user
		if convErr == nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}
	}

	var toCreate map[string]interface{}
	// Unmarshal the JSON data into the values we'll use to create the resource
	if data, err := json.Marshal(createVals); err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	} else {
		json.Unmarshal(data, &toCreate)
	}

	// Attempt to create the new resource and check if it errored at all
	if _, err := c.Conn.Create(c.Table, toCreate); err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
		return
	}

	// Retrieve the newly created record
	res, err := c.ReturnOne(filter)
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Aaaand success
	ctx.JSON(http.StatusCreated, util.MarshalResponse(res))
}

// Delete removes a record
func (c *Command) Delete(ctx *gin.Context) {}
