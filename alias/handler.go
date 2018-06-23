package alias

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

// Alias is the struct that implements the handler interface for the alias resource
type Alias struct {
	Conn  *rethink.Connection // The RethinkDB connection
	Table string              // The database table we're using
}

// Routes returns the routing information for this endpoint
func (a *Alias) Routes() []types.RouteDetails {
	return []types.RouteDetails{
		types.RouteDetails{
			Enabled: true, Path: "/:name", Verb: "DELETE",
			Handler: a.Delete,
			Scopes:  []string{},
		},
		types.RouteDetails{
			Enabled: true, Path: "/:name", Verb: "POST",
			Handler: a.Create,
			Scopes:  []string{"user:manage"},
		},
	}
}

// ReturnOne retrieves a single record given the filter provided
func (a *Alias) ReturnOne(filter map[string]interface{}) (ResponseSchema, error) {
	var response ResponseSchema

	// Retrieve a single record from the DB based on the filter
	fromDB, err := a.Conn.GetSingle(a.Table, filter)
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
func (a *Alias) GetSingle(ctx *gin.Context) {
	token := html.EscapeString(ctx.Param("token"))
	filter := map[string]interface{}{"token": token}

	res, err := a.ReturnOne(filter)
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

// Create creates a new record
func (a *Alias) Create(ctx *gin.Context) {
	token := strings.ToLower(html.EscapeString(ctx.Param("token")))
	name := html.EscapeString(ctx.Param("name"))

	// Declare default values
	createVals := CreationSchema{
		CreatedAt: time.Now().UTC(),
		DeletedAt: 0,
		Token:     token,
		Name:      name,
		Enabled:   true,
	}

	// TODO: Could this be check pulled out into a decorator/middleware of sorts?
	// Do an initial check if it exists
	filter := map[string]interface{}{
		"token": createVals.Token, "name": createVals.Name}
	res, err := a.ReturnOne(filter)

	// Check if it's a RetrievalResult, or an actual error
	if retRes, ok := err.(rethink.RetrievalResult); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	} else if retRes.Success {
		if !retRes.SoftDeleted {
			// It exists already but isn't soft-deleted, error out
			// can't edit from this endpoint
			ctx.AbortWithStatusJSON(http.StatusConflict, util.MarshalResponse(res))
			return
		}
		// It exists and is soft-deleted. Remove that one and then create a new one
		_, err := a.Conn.Delete(a.Table, res.ID)
		if err != nil {
			util.NiceError(ctx, err, http.StatusInternalServerError)
			return
		}
	}

	// No records already exist that match, go ahead with creation
	// Passed validation, put in the user data & prepare the data we're using
	createData, err := util.ValidateAndMap(
		ctx.Request.Body, "/alias/createSchema.json", createVals)

	if validateErr, ok := err.(util.APIError); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	} else if ok {
		// It's a validation error
		ctx.AbortWithStatusJSON(http.StatusBadRequest, validateErr.Data)
		return
	}

	// Validation passed, set the command reference
	lookupFilter := map[string]interface{}{
		"token": token, "name": createData["commandName"]}
	lookupRes, err := a.Conn.GetSingle("commands", lookupFilter)
	if err != nil {
		// Don't want to create repeat without a command to be repeating
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	} else if lookupRes == nil || rethink.IsSoftDeleted(lookupRes) {
		util.NiceError(
			ctx,
			fmt.Errorf(
				"Missing referenced command %s", createData["commandName"]),
			http.StatusBadRequest)
		return
	}

	// Neither test failed, attempt to update createData with the UUID for command
	// Okay without a test (for now) because we know there is a record that was
	// returned from the DB that is non-null and thus must have an ID key
	createData["command"] = lookupRes.(map[string]interface{})["id"]

	// Attempt to create the new resource
	if _, err := a.Conn.Create(a.Table, createData); err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
		return
	}

	response, err := a.ReturnOne(filter)
	// Actual error, not a RetrievalResult
	if _, ok := err.(rethink.RetrievalResult); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Aaaand success
	ctx.Header("X-Total-Count", "1")
	ctx.JSON(http.StatusCreated, util.MarshalResponse(response))
}

// Delete soft-deletes a record
func (a *Alias) Delete(ctx *gin.Context) {
	token := html.EscapeString(ctx.Param("token"))
	name := html.EscapeString(ctx.Param("name"))
	filter := map[string]interface{}{"token": token, "name": name}
	resp, err := a.Conn.GetByFilter(a.Table, filter, 1)

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
		log.Errorf("[%s] - Unable to typecast response to correct type", a.Table)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = a.Conn.Disable(a.Table, rs["id"].(string))
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Success
	ctx.Header("X-Resource-ID-Removed", rs["id"].(string))
	ctx.Status(http.StatusOK)
}

// Stub endpoints

// GetAll is a stub to implement the Handler interface
func (a *Alias) GetAll(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusTeapot)
}

// Update is a stub to implement the Handler interface
func (a *Alias) Update(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusTeapot)
}
