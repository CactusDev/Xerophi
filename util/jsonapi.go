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
// Sub-structs will be placed automatically under their parent (meta/attr) so there is no need to have that tag on
// any sub-struct
func MarshalResponse(s JSONAPISchema) map[string]interface{} {
	var response = make(map[string]interface{})

	var data = make(map[string]interface{})
	ift := reflect.TypeOf(s)
	ifv := reflect.ValueOf(s)
	attributes, meta, id, objectType := pullVals(ift, ifv)

	data["attributes"] = attributes
	data["id"] = id
	data["type"] = objectType
	if len(meta) > 0 {
		response["meta"] = meta
	}

	response["data"] = data

	return response
}

func pullVals(ift reflect.Type, ifv reflect.Value) (map[string]interface{}, map[string]interface{}, string, string) {
	var attr = make(map[string]interface{})
	var meta = make(map[string]interface{})
	var recordType = ""
	var id = ""
	for i := 0; i < ift.NumField(); i++ {
		var value interface{}
		split := GetTags(ift.Field(i))
		if split == nil {
			// It's an anonymous field, ignore it
			continue
		}
		value = ifv.Field(i).Interface()
		if ifv.Field(i).Kind() == reflect.Struct {
			var subID = ""
			var subMeta map[string]interface{}
			value, subMeta, subID, _ = pullVals(ift.Field(i).Type, ifv.Field(i))
			if len(subMeta) > 0 {
				meta[split[0]] = subMeta
			}
			if subID != "" && id == "" {
				id = subID
			}
		}
		// Anything after the first element is tags, figure out which we want
		for _, tag := range split {
			// Need to set the keys w/ their names here if it's a struct
			switch tag {
			case "attr":
				// Attribute
				attr[split[1]] = value
			case "meta":
				// Meta information about the request
				meta[split[1]] = value
			case "primary":
				// It's the primary key/record ID & record type
				id = ifv.Field(i).String()
				recordType = split[1]
			default: // Ignore any other tags
			}
		}
	}

	return attr, meta, id, recordType
}

// GetTags takes a reflect.StructField object and returns a slice of the associated tags
func GetTags(obj reflect.StructField) []string {
	if obj.Anonymous {
		// Anonymous field, don't try to access
		return nil
	}
	// Get the jsonapi tags for this element
	tags := obj.Tag.Get("jsonapi")
	// Split the tags on the , character
	split := strings.Split(tags, ",")

	return split
}

// FieldTag looks up a single field in the provided object and returns the tag for it
func FieldTag(obj interface{}, lookup string, tag string) string {
	ift := reflect.TypeOf(obj)

	// Our code only works with structs
	if ift.Kind() != reflect.Struct {
		return ""
	}

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