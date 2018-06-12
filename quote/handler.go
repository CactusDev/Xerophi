package quote

import (
	"fmt"
	"html"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CactusDev/Xerophi/rethink"
	"github.com/CactusDev/Xerophi/types"
	"github.com/CactusDev/Xerophi/util"

	"github.com/gin-gonic/gin"

	mapstruct "github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

// Quote is the struct that implements the handler interface for the quote resource
type Quote struct {
	Conn  *rethink.Connection // The RethinkDB connection
	Table string              // The database table we're using
}

// Routes returns the routing information for this endpoint
func (q *Quote) Routes() []types.RouteDetails {
	return []types.RouteDetails{
		types.RouteDetails{
			Enabled: true, Path: "", Verb: "GET",
			Handler: q.GetAll,
		},
		types.RouteDetails{
			Enabled: true, Path: "", Verb: "POST",
			Handler: q.Create,
		},
		types.RouteDetails{
			Enabled: true, Path: "/:quoteId", Verb: "GET",
			Handler: q.GetSingle,
		},
		types.RouteDetails{
			Enabled: true, Path: "/:quoteId", Verb: "PATCH",
			Handler: q.Update,
		},
		types.RouteDetails{
			Enabled: true, Path: "/:quoteId", Verb: "DELETE",
			Handler: q.Delete,
		},
	}
}

// ReturnOne retrieves a single record given the filter provided
func (q *Quote) ReturnOne(filter map[string]interface{}) (ResponseSchema, error) {
	var response ResponseSchema

	// Retrieve a single record from the DB based on the filter
	fromDB, err := q.Conn.GetSingle(q.Table, filter)
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
func (q *Quote) GetAll(ctx *gin.Context) {
	// Get parameter values
	token := html.EscapeString(ctx.Param("token"))
	filter := map[string]interface{}{"token": token}

	// Get all the records that match our filter
	fromDB, err := q.Conn.GetByFilter(q.Table, filter, 0)
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

	ctx.Header("X-Total-Count", fmt.Sprint(len(decoded)))
	ctx.JSON(http.StatusOK, response)
}

// GetSingle returns a single record
func (q *Quote) GetSingle(ctx *gin.Context) {
	token := html.EscapeString(ctx.Param("token"))
	filter := map[string]interface{}{"token": token}

	if ctx.Param("quoteId") == "random" {
		q.GetRandom(ctx)
		return
	}

	quoteID, err := strconv.Atoi(ctx.Param("quoteId"))
	if err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
		return
	}
	filter["quoteId"] = quoteID

	res, err := q.ReturnOne(filter)
	retRes, ok := err.(rethink.RetrievalResult)
	// If !ok AND then err != nil then we have an actual error and not a RetRes
	if !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	if retRes.Success && !retRes.SoftDeleted {
		ctx.Header("X-Total-Count", "1")
		ctx.JSON(http.StatusOK, util.MarshalResponse(res))
		return
	}

	// None were found Jim, 404 that boyo
	ctx.AbortWithStatus(http.StatusNotFound)
}

// GetRandom retrieves a random quote if any exist
func (q *Quote) GetRandom(ctx *gin.Context) {
	var resp ResponseSchema
	token := html.EscapeString(ctx.Param("token"))
	// Don't want to retrieve a soft-deleted quote
	filter := map[string]interface{}{"token": token, "deletedAt": 0}
	fromDB, err := q.Conn.GetRandom(q.Table, filter)
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}
	// None exist that match the filter, oh well
	if fromDB == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
	}
	// We made it past the checks, at least one exists, return that
	if err = mapstruct.Decode(fromDB, &resp); err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}
	ctx.Header("X-Total-Count", "1")
	ctx.JSON(http.StatusOK, util.MarshalResponse(resp))
	return
}

// Create creates a new record
func (q *Quote) Create(ctx *gin.Context) {
	// Declare default values
	createVals := CreationSchema{
		CreatedAt: time.Now().UTC(),
		DeletedAt: 0,
		Token:     strings.ToLower(html.EscapeString(ctx.Param("token"))),
		QuoteID:   rand.Intn(10),
		Enabled:   true,
		// QuoteID: nextQuoteNum from user object for :token
	}

	// // Check if it exists yet
	filter := map[string]interface{}{
		"token": createVals.Token, "quoteId": createVals.QuoteID}
	res, err := q.ReturnOne(filter)

	// If !ok AND then err != nil then we have an actual error and not a RetRes
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
		_, err := q.Conn.Delete(q.Table, res.ID)
		if err != nil {
			util.NiceError(ctx, err, http.StatusInternalServerError)
			return
		}
	}

	// No records already exist that match, go ahead with creation
	// Passed validation, put in the user data & prepare the data we're using
	createData, err := util.ValidateAndMap(
		ctx.Request.Body, "/quote/createSchema.json", createVals)

	if validateErr, ok := err.(util.APIError); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	} else if ok {
		// It's a validation error
		ctx.AbortWithStatusJSON(http.StatusBadRequest, validateErr.Data)
	}

	// Attempt to create the new resource
	if _, err := q.Conn.Create(q.Table, createData); err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
		return
	}

	// Retrieve the newly created record
	res, err = q.ReturnOne(filter)
	// If !ok AND then err != nil then we have an actual error and not a RetRes
	if _, ok := err.(rethink.RetrievalResult); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Aaaand success
	ctx.Header("X-Total-Count", "1")
	ctx.JSON(http.StatusCreated, util.MarshalResponse(res))
}

// Update handles the updating of a record if the record exists
func (q *Quote) Update(ctx *gin.Context) {
	// Get the data we need from the request
	token := strings.ToLower(html.EscapeString(ctx.Param("token")))

	// Check if the resource that we want to edit exists
	filter := map[string]interface{}{"token": token}
	quoteID, err := strconv.Atoi(ctx.Param("quoteId"))
	if err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
		return
	}
	filter["quoteId"] = quoteID
	resp, err := q.ReturnOne(filter)

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
		ctx.Request.Body, "/quote/schema.json", updateVals)

	if validateErr, ok := err.(util.APIError); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	} else if ok {
		// It's a validation error
		ctx.AbortWithStatusJSON(http.StatusBadRequest, validateErr.Data)
	}

	// Attempt to update the new resource
	_, err = q.Conn.Update(q.Table, resp.ID, updateData)
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Retrieve the newly updated record
	response, err := q.ReturnOne(filter)
	// If !ok AND then err != nil then we have an actual error and not a RetRes
	if _, ok := err.(rethink.RetrievalResult); !ok && err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Success
	ctx.Header("X-Total-Count", "1")
	ctx.JSON(http.StatusOK, util.MarshalResponse(response))
}

// Delete soft-deletes a record
func (q *Quote) Delete(ctx *gin.Context) {
	token := strings.ToLower(html.EscapeString(ctx.Param("token")))
	filter := map[string]interface{}{"token": token}
	// Convert the quoteId parameter to an int
	quoteID, err := strconv.Atoi(ctx.Param("quoteId"))
	if err != nil {
		util.NiceError(ctx, err, http.StatusBadRequest)
		return
	}
	filter["quoteId"] = quoteID
	resp, err := q.Conn.GetByFilter(q.Table, filter, 1)

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
		log.Errorf("[%s] - Unable to typecast response to correct type", q.Table)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Soft-delete the record
	_, err = q.Conn.Disable(q.Table, rs["id"].(string))
	if err != nil {
		util.NiceError(ctx, err, http.StatusInternalServerError)
		return
	}

	// Success
	ctx.Header("X-Resource-ID-Removed", rs["id"].(string))
	ctx.Status(http.StatusOK)
}
