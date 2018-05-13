package command

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/CactusDev/Xerophi/rethink"
	"github.com/CactusDev/Xerophi/types"
	"github.com/CactusDev/Xerophi/util"

	"github.com/gin-gonic/gin"

	mapstruct "github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

// Command is the struct that implements the handler interface for the command resource
type Command struct {
	Conn  *rethink.Connection // The RethinkDB connection
	Table string              // The database table we're using
}

// Routes returns the routing information for this endpoint
func (c *Command) Routes() []types.RouteDetails {
	return []types.RouteDetails{
		types.RouteDetails{
			Enabled: true, Path: "", Verb: "GET",
			Handler: c.GetAll,
		},
		types.RouteDetails{
			Enabled: true, Path: "/:name", Verb: "GET",
			Handler: c.GetSingle,
		},
		types.RouteDetails{
			Enabled: true, Path: "/:name", Verb: "PATCH",
			Handler: c.Update,
		},
		types.RouteDetails{
			Enabled: true, Path: "/:name", Verb: "POST",
			Handler: c.Create,
		},
		types.RouteDetails{
			Enabled: true, Path: "/:name", Verb: "DELETE",
			Handler: c.Delete,
		},
	}
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
		return response, rethink.RetrievalResult{
			Success: false, SoftDeleted: false, Message: ""}
	}

	// Decode the response from the DB into the response schema object
	if err = mapstruct.Decode(fromDB, &response); err != nil {
		return response, err
	}

	if fromDB.(map[string]interface{})["deletedAt"].(float64) != 0 {
		return response, rethink.RetrievalResult{true, true, ""}
	}

	return response, rethink.RetrievalResult{true, false, ""}
}

// GetAll returns all records associated with the token
func (c *Command) GetAll(ctx *gin.Context) {
	token := strings.ToLower(html.EscapeString(ctx.Param("token")))
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

	ctx.Header("x-total-count", fmt.Sprint(len(decoded)))
	ctx.JSON(http.StatusOK, response)
}

// GetSingle returns a single record
func (c *Command) GetSingle(ctx *gin.Context) {
	token := html.EscapeString(ctx.Param("token"))
	name := html.EscapeString(ctx.Param("name"))
	filter := map[string]interface{}{"token": token, "name": name}

	res, err := c.ReturnOne(filter)
	retRes, ok := err.(rethink.RetrievalResult)
	// If !ok AND then err != nil then we have an actual error and not a RetRes
	if !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// If we find one then we're just going to return right away
	if retRes.Success && !retRes.SoftDeleted {
		ctx.Header("x-total-count", "1")
		ctx.JSON(http.StatusOK, util.MarshalResponse(res))
		return
	}

	// None were found Jim, 404 that boyo
	ctx.AbortWithStatus(http.StatusNotFound)
	return
}

// Create creates a new record
func (c *Command) Create(ctx *gin.Context) {
	// Declare default values
	createVals := CreationSchema{
		CreatedAt: time.Now().UTC(),
		DeletedAt: 0,
		Token:     strings.ToLower(html.EscapeString(ctx.Param("token"))),
		Name:      html.EscapeString(ctx.Param("name")),
		Enabled:   true,
	}

	// TODO: Could this be check pulled out into a decorator/middleware of sorts?
	// Do an initial check if it exists
	filter := map[string]interface{}{
		"token": createVals.Token, "name": createVals.Name}
	res, err := c.ReturnOne(filter)

	// Check if it's a RetrievalResult, or an actual error
	if retRes, ok := err.(rethink.RetrievalResult); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	} else if retRes.Success {
		if !retRes.SoftDeleted {
			// We have a record that already exists, abort
			ctx.AbortWithStatusJSON(http.StatusConflict, util.MarshalResponse(res))
			return
		}
		// It exists and is soft-deleted. Remove that one and then create a new one
		_, err := c.Conn.Delete(c.Table, res.ID)
		if err != nil {
			util.NiceError(ctx, err, http.StatusInternalServerError)
			return
		}
	}

	// No records already exist that match, go ahead with creation
	// Passed validation, put in the user data & prepare the data we're using
	createData, err := util.ValidateAndMap(
		ctx.Request.Body, "/command/createSchema.json", createVals)

	if validateErr, ok := err.(util.APIError); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	} else if ok {
		// It's a validation error
		ctx.AbortWithStatusJSON(http.StatusBadRequest, validateErr.Data)
	}

	// Attempt to create the new resource
	if _, err := c.Conn.Create(c.Table, createData); err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
		return
	}

	response, err := c.ReturnOne(filter)
	// Actual error, not a RetrievalResult
	if _, ok := err.(rethink.RetrievalResult); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Aaaand success
	ctx.Header("x-total-count", "1")
	ctx.JSON(http.StatusCreated, util.MarshalResponse(response))
}

// Update handles the updating of a record if the record exists
func (c *Command) Update(ctx *gin.Context) {
	// Get the data we need from the request
	token := html.EscapeString(ctx.Param("token"))
	name := html.EscapeString(ctx.Param("name"))

	// Check if the resource that we want to edit exists
	filter := map[string]interface{}{"token": token, "name": name}
	resp, err := c.ReturnOne(filter)
	if retRes, ok := err.(rethink.RetrievalResult); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	} else if retRes.Success && retRes.SoftDeleted {
		// Record "doesn't exist", abort with a 404
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Made it past the checks, record exists
	// Passed validation, put in the user data & prepare the data we're using
	var updateVals ClientSchema
	updateData, err := util.ValidateAndMap(
		ctx.Request.Body, "/command/schema.json", updateVals)

	if validateErr, ok := err.(util.APIError); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	} else if ok {
		// It's a validation error
		ctx.AbortWithStatusJSON(http.StatusBadRequest, validateErr.Data)
	}

	// Attempt to update the new resource
	_, err = c.Conn.Update(c.Table, resp.ID, updateData)
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Retrieve the newly updated record
	response, err := c.ReturnOne(filter)
	// If !ok AND then err != nil then we have an actual error and not a RetRes
	if _, ok := err.(rethink.RetrievalResult); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Success
	ctx.Header("x-total-count", "1")
	ctx.JSON(http.StatusOK, util.MarshalResponse(response))
}

// Delete soft-deletes a record
func (c *Command) Delete(ctx *gin.Context) {
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

	rs, valid := resp[0].(map[string]interface{})
	if !valid {
		log.Errorf("[%s] - Unable to typecast response to correct type", c.Table)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = c.Conn.Disable(c.Table, rs["id"].(string))
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Success
	ctx.Header("x-resource-id-removed", rs["id"].(string))
	ctx.Status(http.StatusOK)
}
