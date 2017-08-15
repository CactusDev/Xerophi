package util

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"

	log "github.com/Sirupsen/logrus"
)

// MapStringToInterface converts map[string]string to a map[string]interface{}
func MapStringToInterface(vars map[string]string) map[string]interface{} {
	var converted = make(map[string]interface{})

	for k, v := range vars {
		converted[k] = v
	}

	return converted
}

// NiceError factors away the erroring of a function into a clean single-line function call
func NiceError(ctx *gin.Context, err error, code int) {
	log.Error(err.Error())
	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.AbortWithStatus(code)
		return
	}
	ctx.AbortWithStatusJSON(code, ve)
}
