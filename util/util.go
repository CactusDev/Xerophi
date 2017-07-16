package util

// MapStringToInterface converts map[string]string to a map[string]interface{}
func MapStringToInterface(vars map[string]string) map[string]interface{} {
	var converted = make(map[string]interface{})

	for k, v := range vars {
		converted[k] = v
	}

	return converted
}
