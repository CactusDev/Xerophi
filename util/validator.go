package util

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v8"
)

// APIError is an alias for map[string]interface{} for cleaner code
type APIError map[string]interface{}

func convertNamespace(namespace string) []string {
	var split []string
	// Iterate over the namespace string
	// Found a period
	// 	Store everything to the left
	// 	Delete everything to the left now
	// Continue iteration

	return split
}

// ConvertToError converts the error given to a key paired value
// that's human-readable. Returns a nil error object if it was successful
func ConvertToError(err error) (APIError, error) {
	response := APIError{}

	switch err.(type) {
	case *json.UnmarshalTypeError:
		ve, ok := err.(*json.UnmarshalTypeError)
		if !ok {
			log.Error("SPLODEY NOT OKAY!")
		}
		curErr := APIError{}
		curErr["raw"] = ve
		curErr["field"] = convertNamespace(ve.Field)
		curErr["error"] = fmt.Sprintf("Invalid type %s for field %s, require %s.",
			ve.Value, ve.Field, ve.Type.String())

		response[strings.ToLower(ve.Field)] = curErr

	case validator.ValidationErrors:
		for _, err := range err.(validator.ValidationErrors) {
			curErr := APIError{}
			curErr["raw"] = err
			curErr["field"] = convertNamespace(err.NameNamespace)
			switch err.Tag {
			case "required":
				curErr["error"] = "Required field not included"
			}
			response[strings.ToLower(err.Field)] = curErr
		}
	}

	return response, nil
}
