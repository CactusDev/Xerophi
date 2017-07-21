package util

import (
	"reflect"

	log "github.com/Sirupsen/logrus"
)

// MapStringToInterface converts map[string]string to a map[string]interface{}
func MapStringToInterface(vars map[string]string) map[string]interface{} {
	var converted = make(map[string]interface{})

	for k, v := range vars {
		converted[k] = v
	}

	return converted
}

// GetAPITag takes the object given and returns the jsonapi tag value for the specified field
func GetAPITag(obj interface{}, lookup string) string {
	field, ok := reflect.TypeOf(obj).FieldByName(lookup)
	if !ok {
		log.Warn("Uh, stuff happened")
		return ""
	}
	tag, ok := field.Tag.Lookup("jsonapi")
	if !ok {
		log.Warn("Aw snap")
		return ""
	}
	return tag
}
