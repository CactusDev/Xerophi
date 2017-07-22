package util

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	log "github.com/Sirupsen/logrus"
	mapstruct "github.com/mitchellh/mapstructure"
)

// MarshalResponse parses out the data provided into the appropriate JSON API packet format
func MarshalResponse(attributes, meta map[string]interface{}, errors error) map[string]interface{} {
	var response = make(map[string]interface{})

	// Want the id separate and not included argument data/attributes
	if val, exists := attributes["id"]; exists {
		response["id"] = val
		delete(attributes, "id")
	}

	response["data"] = attributes

	// Only include errors key if there's anything to share
	if errors != nil {
		// TODO: Handle this better, deal with more in-depth errors
		response["errors"] = []string{errors.Error()}
	}
	// Only include meta key if there's anything to share
	if len(meta) != 0 {
		response["meta"] = meta
	}
	return response
}

// NiceResponse moves the ugliness of a proper response away from the handler function
func NiceResponse(ctx *gin.Context, schema interface{}, data interface{}, many bool) {
	var multiResponse []map[string]interface{}
	var response map[string]interface{}
	err := mapstruct.Decode(data, &schema)
	bytes, err := json.Marshal(schema)
	if many {
		err = json.Unmarshal(bytes, &multiResponse)
	} else {
		err = json.Unmarshal(bytes, &response)
	}

	if err != nil {
		log.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, MarshalResponse(nil, nil, err))
		return
	}
	if many {
		ctx.JSON(http.StatusOK, MarshalResponse(multiResponse, nil, nil))
	}
	ctx.JSON(http.StatusOK, MarshalResponse(response, nil, nil))
}

// NiceError factors away the erroring of a function into a clean single-line function call
func NiceError(ctx *gin.Context, err error, code int) {
	log.Error(err.Error())
	ctx.AbortWithError(code, err)
}
