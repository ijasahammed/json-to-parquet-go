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
			// field := make(map[string]interface{}, 0)
			if val == nil {
				_, output[name] = Read(name, reflect.TypeOf(""), val)
			} else {
				_, output[name] = Read(name, reflect.TypeOf(val), val)
			}
			// field["name"] = name
			// output[name] = field
		}
	}
	return output
}

func Read(name string, t reflect.Type, val any) (string, map[string]interface{}) {

	switch t.Kind() {
	case reflect.Slice:
		return "LIST", map[string]interface{}{"type": "LIST", "value": readFromSlice(name, t, val)}
	case reflect.Map:
		return "MAP", map[string]interface{}{"type": "MAP", "value": readFromMap(name, t, val)}
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

func readFromSlice(name string, t reflect.Type, val any) map[string]interface{} {

	resultData := make(map[string]interface{}, 0)
	new_val := val.([]interface{})
	for _, v := range new_val {
		_, resultData = Read("element", reflect.TypeOf(v), v)
		break
	}
	return resultData
}

func readFromMap(name string, t reflect.Type, val any) map[string]interface{} {
	new_val := val.(map[string]interface{})
	mapData := make(map[string]interface{}, 0)

	for key, v := range new_val {
		_, mapData[key] = Read(key, reflect.TypeOf(v), v)
	}
	return mapData

}
