package fbUtil

func flatten(firestoreFields interface{}) map[string]interface{} {
	flat := make(map[string]interface{})
	if mapped, ok := firestoreFields.(map[string]interface{}); ok {
		for key, meta := range mapped {
			metaMapped := meta.(map[string]interface{})
			for firestoreType, value := range metaMapped {
				switch {
				case firestoreType == "mapValue":
					if mapMapped, ok := value.(map[string]interface{}); ok {
						flat[key] = flatten(mapMapped["fields"])
					}
				case firestoreType == "arrayValue":
					flatArray := make([]interface{}, 0)
					arrayValuesMapped := value.(map[string]interface{})
					array := arrayValuesMapped["values"].([]interface{})
					for _, element := range array {
						elementMapped := element.(map[string]interface{})
						for _, elementValue := range elementMapped {
							flatArray = append(flatArray, elementValue)
						}
					}
					flat[key] = flatArray
				case firestoreType == "integerValue" || firestoreType == "doubleValue" ||
					firestoreType == "booleanValue" || firestoreType == "stringValue":
					flat[key] = value
				}
			}
		}
	}
	return flat
}
