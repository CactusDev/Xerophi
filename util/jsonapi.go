package util

import (
	"reflect"
	"strings"
)

// JSONAPISchema is an interface used for generating the proper JSON API response packet
type JSONAPISchema interface {
	GetAPITag(lookup string) string
}

// MarshalResponse takes an object that implements the JSONAPISchema interface and marshals it to a map[string]interface{}
func MarshalResponse(s JSONAPISchema) map[string]interface{} {
	var response = make(map[string]interface{})
	var data = make(map[string]interface{})
	// Create attributes & meta maps before adding to main response map
	var attributes = make(map[string]interface{})
	var meta = make(map[string]interface{})

	ift := reflect.TypeOf(s)
	ifv := reflect.ValueOf(s)

	for i := 0; i < ift.NumField(); i++ {
		fieldName := ift.Field(i).Name
		if ift.Field(i).Anonymous {
			// Anonymous field, don't try to access
			continue
		}
		// Get the jsonapi tags for this element
		tags := s.GetAPITag(fieldName)
		// Split the tags on the , character
		split := strings.Split(tags, ",")

		// Anything after the first element is tags, figure out which we want
		for _, tag := range split[1:] {
			value := ifv.Field(i).Interface()
			switch tag {
			case "attributes":
				// Attribute
				attributes[split[0]] = value
			case "meta":
				// Meta information about the request
				meta[split[0]] = value
			case "primary":
				// It's the primary key/record ID
				response["id"] = ifv.Field(i).String()
			default: // Ignore any other tags
			}
		}
	}

	data["attributes"] = attributes
	response["data"] = data

	// Only add it if there's anything *to* add
	if len(meta) != 0 {
		response["meta"] = meta
	}

	return response
}

// ReturnTags takes an interface and a string to look up the tag for.
// If the first argument passed is a struct, then it starts looking for the
// tag given by the second argument & will return any it finds in map form
func ReturnTags(obj interface{}, lookup string) map[string]interface{} {
	response := make(map[string]interface{})
	ift := reflect.TypeOf(obj)
	ifv := reflect.ValueOf(obj)

	// Our code only works with structs
	if ift.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < ift.NumField(); i++ {
		// Skip the current field if it's anonymous
		if ift.Field(i).Anonymous {
			continue
		}
		fieldName := ift.Field(i).Name
		switch ifv.Field(i).Kind() {
		case reflect.Struct:
			// Iterate over the
			resp := ReturnTags(ifv.Field(i).Interface(), lookup)
			response[fieldName] = resp
		default:
			tag, ok := ift.Field(i).Tag.Lookup(lookup)
			if !ok {
				response[fieldName] = nil
			}
			response[fieldName] = tag
		}
	}

	return response
}
