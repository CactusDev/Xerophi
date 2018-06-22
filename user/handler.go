package user

import (
	"html"
	"net/http"
	"strings"

	"github.com/CactusDev/Xerophi/rethink"
	"github.com/CactusDev/Xerophi/types"
	"github.com/CactusDev/Xerophi/util"

	"github.com/gin-gonic/gin"

	mapstruct "github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

// User is the struct that implements the handler interface for the user resource
type User struct {
	Conn  *rethink.Connection // The RethinkDB connection
	Table string              // The database table we're using
}

// Routes returns the routing information for this endpoint
func (c *User) Routes() []types.RouteDetails {
	return []types.RouteDetails{
		types.RouteDetails{
			Enabled: true, Path: "", Verb: "GET",
			Handler: c.GetSingle,
			Scopes:  []string{},
		},
		types.RouteDetails{
			Enabled: true, Path: "/config", Verb: "GET",
			Handler: c.GetConfig,
			Scopes:  []string{},
		},
		types.RouteDetails{
			Enabled: true, Path: "/config", Verb: "PATCH",
			Handler: c.Update,
			Scopes:  []string{"user:manage"},
		},
	}
}

// ReturnOne retrieves a single record given the filter provided
func (c *User) ReturnOne(filter map[string]interface{}) (ResponseSchema, error) {
	var response ResponseSchema

	// Retrieve a single record from the DB based on the filter
	fromDB, err := c.Conn.GetSingle(c.Table, filter)
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

// ReturnConfig retrieves a single config record given the filter provided
func (c *User) ReturnConfig(filter map[string]interface{}) (ConfigSchema, error) {
	var response ConfigSchema

	// Retrieve a single record from the DB based on the filter
	fromDB, err := c.Conn.GetSingle("configs", filter)
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

// GetSingle returns a single record
func (c *User) GetSingle(ctx *gin.Context) {
	token := html.EscapeString(ctx.Param("token"))
	filter := map[string]interface{}{"token": token}

	res, err := c.ReturnOne(filter)
	retRes, ok := err.(rethink.RetrievalResult)
	// If !ok AND then err != nil then we have an actual error and not a RetRes
	if !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// If we find one then we're just going to return right away
	if retRes.Success && !retRes.SoftDeleted {
		ctx.Header("X-Total-Count", "1")
		ctx.JSON(http.StatusOK, util.MarshalResponse(res))
		return
	}

	// None were found Jim, 404 that boyo
	ctx.AbortWithStatus(http.StatusNotFound)
	return
}

// Update handles the updating of a record if the record exists
func (c *User) Update(ctx *gin.Context) {
	// Get the data we need from the request
	token := strings.ToLower(html.EscapeString(ctx.Param("token")))

	// Check if the resource that we want to edit exists
	filter := map[string]interface{}{"token": token}
	resp, err := c.ReturnConfig(filter)
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
	var updateVals UpdateSchema
	updateData, err := util.ValidateAndMap(
		ctx.Request.Body, "/user/schema.json", updateVals)

	log.Debug(updateData)

	if validateErr, ok := err.(util.APIError); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	} else if ok {
		// It's a validation error
		ctx.AbortWithStatusJSON(http.StatusBadRequest, validateErr.Data)
		return
	}

	// Attempt to update the new resource
	_, err = c.Conn.Update("configs", resp.ID, updateData)
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Retrieve the newly updated record
	response, err := c.ReturnConfig(filter)
	// If !ok AND then err != nil then we have an actual error and not a RetRes
	if _, ok := err.(rethink.RetrievalResult); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Success
	ctx.Header("X-Total-Count", "1")
	ctx.JSON(http.StatusOK, util.MarshalResponse(response))
}

// GetConfig returns a single config record
func (c *User) GetConfig(ctx *gin.Context) {
	token := html.EscapeString(ctx.Param("token"))
	filter := map[string]interface{}{"token": token}

	res, err := c.ReturnConfig(filter)
	retRes, ok := err.(rethink.RetrievalResult)
	// If !ok AND then err != nil then we have an actual error and not a RetRes
	if !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// If we find one then we're just going to return right away
	if retRes.Success && !retRes.SoftDeleted {
		ctx.Header("X-Total-Count", "1")
		ctx.JSON(http.StatusOK, util.MarshalResponse(res))
		return
	}

	// None were found Jim, 404 that boyo
	ctx.AbortWithStatus(http.StatusNotFound)
	return
}

// Stubs - non-used functions but required for implementing Handler interface

// Create is a stub required to implement Handler interface, but it's not actually used
func (c *User) Create(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusTeapot)
}

// GetAll is a stub required to implement Handler interface, but it's not actually used
func (c *User) GetAll(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusTeapot)
}

// Delete is a stub required to implement Handler interface, but it's not actually used
func (c *User) Delete(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusTeapot)
}
