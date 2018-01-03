package util

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"

	log "github.com/sirupsen/logrus"
)

// GetFromOffset generates a human-readable line/character error from a
// JSON error's offset
func GetFromOffset(input string, offset int) (string, error) {
	var line = 1 // Silly humans tending to start counting at 1
	var character = 0

	lf := rune(0x0A)

	if offset > len(input) || offset < 0 {
		return "", fmt.Errorf("Offset beyond bounds of input")
	}

	for pos, char := range input {
		if char == lf {
			// Newline
			character = 0
			line++
		}
		character++
		// We're at the problem-causing character
		if pos == offset {
			break
		}
	}

	return fmt.Sprintf("JSON syntax error at line %d, character %d", line, character), nil
}

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

func GetResourceID(data interface{}) (string, error) {
	var mapped map[string]interface{}
	var isMap bool

	switch data.(type) {
	case []interface{}:
		mapped, isMap = data.([]interface{})[0].(map[string]interface{})
		if !isMap {
			return "", errors.New("Data could not be converted to a map")
		}
	case map[string]interface{}:
		mapped, isMap = data.(map[string]interface{})
		if !isMap {
			return "", errors.New("Data could not be converted to a map")
		}
	}

	id, ok := mapped["id"].(string)
	if id != "" && ok {
		return id, nil
	} else if !ok {
		return "", errors.New("id field is not type string")
	}

	return "", errors.New("No id field, unable to retrieve the ID")
}
