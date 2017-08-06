package util

import (
	"reflect"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// JSONAPISchema is an interface used for generating the proper JSON API response packet
type JSONAPISchema interface {
	GetAPITag(lookup string) string
}

// MarshalResponse takes an object that implements the JSONAPISchema interface and marshals it to a map[string]interface{}
// Sub-structs will be placed automatically under their parent (meta/attr) so there is no need to have that tag on
// any sub-struct
func MarshalResponse(s JSONAPISchema) map[string]interface{} {
	var response = make(map[string]interface{})
	var data = make(map[string]interface{})
	// Create attributes & meta maps before adding to main response map
	var attributes = make(map[string]interface{})
	var meta = make(map[string]interface{})

	ift := reflect.TypeOf(s)
	ifv := reflect.ValueOf(s)

	for i := 0; i < ift.NumField(); i++ {
		split := GetTags(ift.Field(i), s)
		if split == nil {
			// It's an anonymous field, ignore it
			continue
		}
		value := ifv.Field(i).Interface()
		// Anything after the first element is tags, figure out which we want
		for _, tag := range split[1:] {
			// Need to set the keys w/ their names here if it's a struct
			switch tag {
			case "attr":
				// Attribute
				attributes[split[0]] = value
			case "meta":
				// Meta information about the request
				meta[split[0]] = value
			case "primary":
				// It's the primary key/record ID
				attributes["id"] = ifv.Field(i).String()
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

// GetTags takes a reflect.StructField object and returns the appropriate jsonapi tagged field name,
// and a slice of the associated tags
func GetTags(obj reflect.StructField, s JSONAPISchema) []string {
	fieldName := obj.Name
	if obj.Anonymous {
		// Anonymous field, don't try to access
		return nil
	}
	if obj.Type.Kind() == reflect.Struct {

	}
	// Get the jsonapi tags for this element
	tags := s.GetAPITag(fieldName)
	// Split the tags on the , character
	split := strings.Split(tags, ",")
	log.Debugf("[%s]\t%s", fieldName, split)

	return split
}

// FieldTag looks up a single field in the provided object and returns the tag for it
func FieldTag(obj interface{}, lookup string, tag string) string {
	ift := reflect.TypeOf(obj)

	// Our code only works with structs
	if ift.Kind() != reflect.Struct {
		return ""
	}

	// TODO: If it's a struct, then we need to iterate into it and get the tags for it
	field, ok := ift.FieldByName(lookup)
	if !ok {
		return ""
	}
	if field.Type.Kind() == reflect.Struct {

	}
	resp, ok := field.Tag.Lookup(tag)
	if !ok {
		return ""
	}

	return resp
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
