package tools

import (
	"encoding/json"
	"fmt"
)

// 结构体转为json
func Struct2Json(obj any) string {
	str, err := json.Marshal(obj)
	if err != nil {
		panic(fmt.Sprintf("[Struct2Json]转换异常: %v", err))
	}
	return string(str)
}

// json转为结构体
func Json2Struct(str string, obj any) {
	// 将json转为结构体
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		panic(fmt.Sprintf("[Json2Struct]转换异常: %v", err))
	}
}

// json interface转为结构体
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
