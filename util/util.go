package util

import (
	"errors"
	"fmt"
	"reflect"
)

// SetField sets the value of a field in the object
func SetField(obj interface{}, name string, value interface{}) error {
	structVal := reflect.ValueOf(obj).Elem()
	structFieldVal := structVal.FieldByName(name)

	if !structFieldVal.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}
	if !structFieldVal.CanSet() {
		return fmt.Errorf("Unable to set %s field value", name)
	}

	structFieldType := structFieldVal.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		invalidTypeError := errors.New("Provided value does not match type")
		return invalidTypeError
	}

	structFieldVal.Set(val)
	return nil
}

// FillStruct takes a map and an interface (usually a struct) and fills it
func FillStruct(m map[string]interface{}, obj interface{}) error {
	for key, v := range m {
		err := SetField(obj, key, v)
		if err != nil {
			return err
		}
	}
	return nil
}
