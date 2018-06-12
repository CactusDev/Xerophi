package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/CactusDev/Xerophi/types"
	jschema "github.com/xeipuuv/gojsonschema"
)

// APIError is an alias for map[string]string to make cleaner code
type APIError struct {
	Data map[string]interface{}
}

// Error allows APIError to be returned as an error object
func (e APIError) Error() string {
	var errorString string
	// We know that APIError is map[string]interface{}, so no error catching on type-assertion
	for field, msg := range e.Data {
		errorString += fmt.Sprintf("[%s] \"%s\" ", field, msg)
	}
	return errorString
}

// ValidateAndMap takes an io reader (most like io.ReadCloser from the request),
// reads in the data, and then validates it against the JSONSchema provided.
// Then it converts it into the provided schema via the interface (it's always a
// struct as we use it) to apply the JSON  stuff we want and then returns a map
func ValidateAndMap(in io.Reader, schemaPath string, schema types.Schema) (map[string]interface{}, error) {
	var conv map[string]interface{}
	// Retrieve the data from the reader, most likely the request body
	bodyData, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	// Validate the data
	err = ValidateInput(bodyData, schemaPath)
	// It's not an APIError and an actual error exists
	if validateErr, ok := err.(APIError); !ok && err != nil {
		return nil, err
	} else if ok {
		// It's a validation error
		return nil, validateErr
	}

	// Dump the body data into the schema and get the bytes back
	schemaBytes, err := schema.DumpBody(bodyData)
	if err != nil {
		return nil, err
	}

	// Unmarshal our byte array into the output
	if err := json.Unmarshal(schemaBytes, &conv); err != nil {
		return nil, err
	}

	return conv, nil
}

// ValidateInput valids the data provided against the provided JSON schema
// Will only return an error if there's a problem with the data
func ValidateInput(source []byte, schema string) error {
	errors := APIError{Data: make(map[string]interface{})}

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	// TODO: Make schemaLoader static so we're not constantly re-creating a schema
	// 			 each time
	schemaLoader := jschema.NewReferenceLoader("file://" + path + schema)
	dataLoader := jschema.NewBytesLoader(source)

	res, err := jschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return err
	}

	if res.Errors() != nil {
		for _, jschemaErr := range res.Errors() {
			errors.Data[jschemaErr.Field()] = jschemaErr.Description()
		}
		// There were errors, return those
		return errors
	}

	// Passed validation
	return nil
}
