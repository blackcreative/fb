package fb

import (
	"fmt"
	"strconv"

	"google.golang.org/genproto/googleapis/type/latlng"
)

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	Fields interface{} `json:"fields"`
}

// String is a wrapper on strings to make it easier to unwrap Firestore String Arrays
type String struct {
	Value string `json:"stringValue"`
}

// Number is a wrapper on strings to make it easier to unwrap Firestore String Arrays
type Number struct {
	Value int `json:"numberValue"`
}

// GetGeoPointValue extracts a geopoint (*latlng.LatLng) value from a FirestoreValue
func (v FirestoreValue) GetGeoPointValue(name string) (*latlng.LatLng, error) {
	if mappedField, ok := getMappedFieldFromFirestoreValue(name, v); ok {
		if geopoint, ok := mappedField[`json:"geoPointValue"`].(*latlng.LatLng); ok {
			return geopoint, nil
		}
	}
	return nil, fmt.Errorf("Error extracting value named %s from %+v as geopoint. Also check to make sure \"stringValue\" is still a part of FirestoreEvent.FirestoreValue Json", name, v.Fields)
}

// GetStringValue extracts a string value from a FirestoreValue
func (v FirestoreValue) GetStringValue(name string) (string, error) {
	if mappedField, ok := getMappedFieldFromFirestoreValue(name, v); ok {
		if strValue, ok := mappedField["stringValue"].(string); ok {
			return strValue, nil
		}
	}
	return "", fmt.Errorf("Error extracting value named %s from %+v as string. Also check to make sure \"stringValue\" is still a part of FirestoreEvent.FirestoreValue Json", name, v.Fields)
}

// GetIntegerValue extracts an integer value from a FirestoreValue
func (v FirestoreValue) GetIntegerValue(name string) (int, error) {
	if mappedField, ok := getMappedFieldFromFirestoreValue(name, v); ok {
		if strValue, ok := mappedField["integerValue"].(string); ok {
			if value, err := strconv.Atoi(strValue); err == nil {
				return value, nil
			}
		}
	}
	return 0, fmt.Errorf("Error extracting value named %s from %+v as integer. Also check to make sure \"integerValue\" is still a part of FirestoreEvent.FirestoreValue Json", name, v.Fields)
}

// GetStringArray extracts a []string value from a FirestoreValue
func (v FirestoreValue) GetStringArray(name string) (*[]string, error) {
	if array, err := getArrayFromFirestoreValue(name, v); err != nil {
		stringArray := make([]string, len(array))
		for index, any := range array {
			if fbString, ok := any.(String); ok {
				stringArray[index] = fbString.Value
			}
		}
		return &stringArray, nil
	} else {
		return nil, err
	}
}

// GetIntegerArray extracts an integer value from a FirestoreValue
func (v FirestoreValue) GetIntegerArray(name string) (*[]int, error) {
	if array, err := getArrayFromFirestoreValue(name, v); err != nil {
		intArray := make([]int, len(array))
		for index, any := range array {
			if fbNumber, ok := any.(Number); ok {
				intArray[index] = fbNumber.Value
			}
		}
		return &intArray, nil
	} else {
		return nil, err
	}
}

func getMappedFieldFromFirestoreValue(name string, v FirestoreValue) (map[string]interface{}, bool) {
	fields, ok := v.Fields.(map[string]interface{})
	mapped, ok := fields[name].(map[string]interface{})
	return mapped, ok
}

func getArrayFromFirestoreValue(name string, v FirestoreValue) ([]interface{}, error) {
	if mappedField, ok := getMappedFieldFromFirestoreValue(name, v); ok {
		arrayWrapper := mappedField["arrayValue"].(map[string]interface{})
		if array, ok := arrayWrapper["values"].([]interface{}); ok {
			return array, nil
		}
	}
	return nil, fmt.Errorf("Error extracting value named %s from %+v. Also check to make sure \"arrayValue\" and \"values\" are still a part of FirestoreEvent.FirestoreValue Json", name, v.Fields)
}
