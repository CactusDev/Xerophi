package util

import (
	"strings"

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
func ConvertToError(errors error) (APIError, error) {
	response := APIError{}

	if _, ok := errors.(validator.ValidationErrors); !ok {
		return nil, errors
	}

	for _, err := range errors.(validator.ValidationErrors) {
		curErr := APIError{}
		curErr["raw"] = err
		curErr["field"] = convertNamespace(err.NameNamespace)
		switch err.Tag {
		case "required":
			curErr["error"] = "Required field not included"
		}
		response[strings.ToLower(err.Field)] = curErr
	}

	return response, nil
}
