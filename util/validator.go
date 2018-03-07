package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
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
	log.Debug(string(source))
	var errors = make(APIError)
	schemaStream, err := ioutil.ReadFile("./" + schema)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	schemaLoader := jschema.NewBytesLoader(schemaStream)
	dataLoader := jschema.NewBytesLoader(source)
	res, err := jschema.Validate(schemaLoader, dataLoader)

	log.Error(err)

	if err != nil {
		return errors, err
	}

	if res.Errors() != nil {
		for _, err := range res.Errors() {
			errors[err.Field()] = err.String()
		}
	}

	return errors, err
}
