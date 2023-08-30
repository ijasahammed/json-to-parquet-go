/*
Basic json-schema generator based on Go types, for easy interchange of Go
structures between diferent languages.
*/
package jsonstruct

import (
	"fmt"
	"reflect"
)

func BuildJsonStruct(data []map[string]interface{}) map[string]interface{} {

	output := make(map[string]interface{}, 0)
	for _, v := range data {
		for name, val := range v {
			preValue := make(map[string]interface{}, 0)

			if _, exist := output[name]; exist {
				keyPre := output[name].(map[string]interface{})
				if keyPre["type"] == "MAP" || keyPre["type"] == "LIST" {
					preValue = keyPre["value"].(map[string]interface{})
				}

			}
			if val == nil {
				_, output[name] = Read(name, reflect.TypeOf(""), val, preValue)
			} else {
				_, output[name] = Read(name, reflect.TypeOf(val), val, preValue)
			}
		}
	}
	return output
}

func Read(name string, t reflect.Type, val any, preValue map[string]interface{}) (string, map[string]interface{}) {

	switch t.Kind() {
	case reflect.Slice:
		return "LIST", map[string]interface{}{"type": "LIST", "value": readFromSlice(name, t, val, preValue)}
	case reflect.Map:
		return "MAP", map[string]interface{}{"type": "MAP", "value": readFromMap(name, t, val, preValue)}
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return "INT", map[string]interface{}{"type": "INT", "value": 1}
	case reflect.Float32, reflect.Float64:
		return "FLOAT", map[string]interface{}{"type": "FLOAT", "value": 1.0}
	case reflect.Bool:
		return "BOOLEAN", map[string]interface{}{"type": "BOOLEAN", "value": true}
	case reflect.String:
		return "BYTE_ARRAY", map[string]interface{}{"type": "BYTE_ARRAY", "value": "s"}
	}
	fmt.Println("ELseee caseee")
	return "", nil
}

func readFromSlice(name string, t reflect.Type, val any, preValue map[string]interface{}) map[string]interface{} {

	resultData := make(map[string]interface{}, 0)
	if typeVal, exists := preValue["type"]; !exists || !(typeVal == "MAP") {
		preValue = make(map[string]interface{}, 0)
	} else {
		preValue = preValue["value"].(map[string]interface{})
	}
	new_val := val.([]interface{})
	var dataType string
	for _, v := range new_val {
		dataType, resultData = Read("element", reflect.TypeOf(v), v, preValue)
		if dataType != "MAP" {
			break
		} else {
			preValue = resultData["value"].(map[string]interface{})
		}
	}
	return resultData
}

func readFromMap(name string, t reflect.Type, val any, preValue map[string]interface{}) map[string]interface{} {
	new_val := val.(map[string]interface{})
	mapData := preValue

	for key, v := range new_val {
		keyPreValue := make(map[string]interface{}, 0)
		if _, exits := mapData[key]; exits {
			keyPre := mapData[key].(map[string]interface{})
			if keyPre["type"] == "MAP" || keyPre["type"] == "LIST" {
				keyPreValue = keyPre["value"].(map[string]interface{})
			}
		}
		if len(keyPreValue) > 0 && reflect.TypeOf(v).Kind() == reflect.Slice && len(v.([]interface{})) == 0 {
			continue
		}
		_, mapData[key] = Read(key, reflect.TypeOf(v), v, keyPreValue)

	}
	return mapData

}
