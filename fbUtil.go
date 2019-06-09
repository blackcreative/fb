package fbUtil

// Flatten flattens a FirestoreValue.Fields object
func Flatten(firestoreFields interface{}) map[string]interface{} {
	var flat map[string]interface{}
	if mapped, ok := firestoreFields.(map[string]interface{}); ok {
		flat = make(map[string]interface{})
		for key, meta := range mapped {
			metaMapped := meta.(map[string]interface{})
			for firestoreType, value := range metaMapped {
				switch {
				case firestoreType == "mapValue":
					if mapMapped, ok := value.(map[string]interface{}); ok {
						flat[key] = Flatten(mapMapped["fields"])
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
