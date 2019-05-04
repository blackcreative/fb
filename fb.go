package fb

import (
	"fmt"
	"strconv"
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

type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GetGeoPointValue extracts a geopoint (*latlng.LatLng) value from a FirestoreValue
func (v FirestoreValue) GetGeoPointValue(name string) (*GeoPoint, error) {
	if mappedField, ok := getMappedFieldFromFirestoreValue(name, v); ok {
		if geoPointWrapper, ok := mappedField["geoPointValue"].(map[string]float64); ok {
			var geoPoint = GeoPoint{
				Latitude:  geoPointWrapper["latitude"],
				Longitude: geoPointWrapper["longitude"]}
			return &geoPoint, nil
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

func getArrayFromFirestoreValue(name string, v FirestoreValue) ([]interface{}, error) {
	if mappedField, ok := getMappedFieldFromFirestoreValue(name, v); ok {
		arrayWrapper := mappedField["arrayValue"].(map[string]interface{})
		if array, ok := arrayWrapper["values"].([]interface{}); ok {
			return array, nil
		}
	}
	return nil, fmt.Errorf("Error extracting value named %s from %+v. Also check to make sure \"arrayValue\" and \"values\" are still a part of FirestoreEvent.FirestoreValue Json", name, v.Fields)
}

func getMappedFieldFromFirestoreValue(name string, v FirestoreValue) (map[string]interface{}, bool) {
	fields, ok := v.Fields.(map[string]interface{})
	mapped, ok := fields[name].(map[string]interface{})
	return mapped, ok
}
