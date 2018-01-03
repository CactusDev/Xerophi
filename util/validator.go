package util

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	jschema "github.com/xeipuuv/gojsonschema"
)

// ValidateInput valids the data provided against the provided JSON schema
func ValidateInput(source string, schema string) (map[string]string, error) {
	var errors = make(map[string]string)
	schemaStream, err := ioutil.ReadFile("./" + schema)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	schemaLoader := jschema.NewBytesLoader(schemaStream)
	dataLoader := jschema.NewStringLoader(source)
	res, err := jschema.Validate(schemaLoader, dataLoader)

	for _, err := range res.Errors() {
		errors[err.Field()] = err.String()
	}

	return errors, err
}
