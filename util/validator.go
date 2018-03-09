package util

import (
	"encoding/json"
	"fmt"
	"os"

	jschema "github.com/xeipuuv/gojsonschema"
)

// APIError is an alias for map[string]string to make cleaner code
type APIError map[string]string

// HandleJSONErrors handles the JSON decoding errors that occur and converts
// them into more human-readable errors
func HandleJSONErrors(source []byte, err error) APIError {
	var resp = make(APIError)
	switch err.(type) {
	case *json.SyntaxError:
		conv, ok := err.(*json.SyntaxError)
		resp["document"] = fmt.Sprintf("Invalid JSON at offset %d", conv.Offset)
		offset, _ := GetFromOffset(string(source), int(conv.Offset))
		resp["offset"] = offset
		if !ok {
			// We failed to convert to the specific type
			return resp
		}
	case *json.UnmarshalTypeError:
		conv, ok := err.(*json.UnmarshalTypeError)
		resp[conv.Field] = fmt.Sprintf(
			"Invalid type for field %s. Expected %s and received %s.",
			conv.Field, conv.Type.String(), conv.Value)
		if !ok {
			// We failed to convert to the specific type
			return resp
		}
	}

	return resp
}

// ValidateInput valids the data provided against the provided JSON schema
func ValidateInput(source []byte, schema string) (APIError, error) {
	var errors = make(APIError)

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// TODO: Make schemaLoader static so we're not constantly re-creating a schema
	// 			 each time
	schemaLoader := jschema.NewReferenceLoader("file://" + path + schema)
	dataLoader := jschema.NewBytesLoader(source)

	res, err := jschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return nil, err
	}

	if res.Errors() != nil {
		for _, err := range res.Errors() {
			errors[err.Field()] = err.Description()
		}
		// There were errors, return those
		return errors, nil
	}

	// Passed validation
	return nil, err
}
