package util

// All this disabled for now until I can put more time into figuring out embedded structs & their tags

// JSONAPISchema is an interface used for generating the proper JSON API response packet
// type JSONAPISchema interface {
// 	GetAPITag(lookup string) string
// }

// func convertType(val reflect.Value) (interface{}, error) {

// 	switch val.Kind() {
// 	case reflect.String:
// 		return val.String(), nil
// 	case reflect.Int:
// 		return val.Int(), nil
// 	case reflect.Bool:
// 		return val.Bool(), nil
// 	case reflect.Struct:
// 		var fill = make(map[string]interface{})
// 		for i := 0; i < val.NumField(); i++ {
// 			conv, err := convertType(val.Field(i))
// 			if err != nil {
// 				log.Error(err.Error())
// 			}
// 			name := GetAPITag(val, val.Type().Field(i).Name)
// 			fill[name] = conv
// 		}
// 		return fill, nil
// 	}

// 	// Didn't match anything, just return an empty string :/
// 	return "", errors.New("Kind unable to be matched")
// }

// // MarshalResponse takes an object that implements the JSONAPISchema interface and marshals it to a map[string]interface{}
// func MarshalResponse(s JSONAPISchema) map[string]interface{} {
// 	var response = make(map[string]interface{})
// 	var data = make(map[string]interface{})
// 	// Create attributes & meta maps before adding to main response map
// 	var attributes = make(map[string]interface{})
// 	var meta = make(map[string]interface{})

// 	typ := reflect.TypeOf(s)
// 	val := reflect.ValueOf(s)

// 	for i := 0; i < typ.NumField(); i++ {
// 		fieldName := typ.Field(i).Name
// 		if typ.Field(i).Anonymous {
// 			// Anonymous field, don't try to access
// 			continue
// 		}
// 		// Get the jsonapi tags for this element
// 		tags := s.GetAPITag(fieldName)
// 		// Split the tags on the , character
// 		split := strings.Split(tags, ",")

// 		// Anything after the first element is tags, figure out which we want
// 		for _, tag := range split[1:] {
// 			switch tag {
// 			case "attributes":
// 				// Attribute
// 				value, err := convertType(val.Field(i))
// 				if err != nil {
// 					log.Error(err.Error())
// 					continue
// 				}
// 				attributes[split[0]] = value
// 			case "meta":
// 				// Meta information about the request
// 				value, err := convertType(val.Field(i))
// 				if err != nil {
// 					log.Error(err.Error())
// 					continue
// 				}
// 				meta[split[0]] = value
// 			case "primary":
// 				// It's the primary key/record ID
// 				response["id"] = val.Field(i).String()
// 			default: // Ignore any other tags
// 			}
// 		}
// 	}

// 	data["attributes"] = attributes
// 	response["data"] = data

// 	// Only add it if there's anything *to* add
// 	if len(meta) != 0 {
// 		response["meta"] = meta
// 	}

// 	return response
// }

// GetAPITag takes the object given and returns the jsonapi tag value for the specified field
// func GetAPITag(obj interface{}, lookup string) string {
// 	field, ok := reflect.TypeOf(obj).FieldByName(lookup)
// 	if !ok {
// 		log.Warn("Uh, stuff happened")
// 		return ""
// 	}
// 	tag, ok := field.Tag.Lookup("jsonapi")
// 	if !ok {
// 		log.Warn("Aw snap")
// 		return ""
// 	}
// 	return tag
// }
