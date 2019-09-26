package util

import (
	"encoding/json"
)

func GetMapFromJsonString(jsonString string) (ret map[string]interface{}) {
	var i interface{}
	json.Unmarshal([]byte(jsonString), &i)
	ret = i.(map[string]interface{})
	return
}

func GetJsonStringFromMap(i interface{}) string {
	jsonString, _ := json.Marshal(i)
	return string(jsonString)
}
