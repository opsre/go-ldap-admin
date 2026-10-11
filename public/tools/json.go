package tools

import (
	"encoding/json"
	"fmt"
)

// Convert a struct to JSON
func Struct2Json(obj any) string {
	str, err := json.Marshal(obj)
	if err != nil {
		panic(fmt.Sprintf("[Struct2Json] conversion error: %v", err))
	}
	return string(str)
}

// Convert JSON to a struct
func Json2Struct(str string, obj any) {
	// Convert JSON to a struct
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		panic(fmt.Sprintf("[Json2Struct] conversion error: %v", err))
	}
}

// Convert a JSON interface to a struct
func JsonI2Struct(str any, obj any) {
	JsonStr := str.(string)
	Json2Struct(JsonStr, obj)
}

// Convert json string to map
func JsonToMap(jsonStr string) (m map[string]string, err error) {
	err = json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}
	return
}

// Convert map to json string
func MapToJson(m map[string]string) (string, error) {
	result, err := json.Marshal(m)
	if err != nil {
		return "", nil
	}
	return string(result), nil
}
