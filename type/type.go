/*
Basic json-schema generator based on Go types, for easy interchange of Go
structures between diferent languages.
*/
package datatype

import (
	"fmt"
	"reflect"
	"strings"
)

func BuildMedadata(data map[string]interface{}) string {
	metadata := (`{"Tag":"name=parquet-go","Fields":[`)
	for name, val := range data {
		if val == nil {
			metadata += Read(name, reflect.TypeOf(""), val)
		} else {
			metadata += Read(name, reflect.TypeOf(val), val)
		}
	}
	metadata = strings.TrimSuffix(metadata, ",")
	metadata += `]}`
	return metadata
}

func Read(name string, t reflect.Type, val any) string {

	switch t.Kind() {
	case reflect.Slice:
		return readFromSlice(name, t, val)
	case reflect.Map:
		return readFromMap(name, t, val)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return fmt.Sprintf(`{"Tag":"name=%s, type=INT64,repetitiontype=OPTIONAL"},`, name)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf(`{"Tag":"name=%s, type=FLOAT,repetitiontype=OPTIONAL"},`, name)
	case reflect.Bool:
		return fmt.Sprintf(`{"Tag":"name=%s, type=BOOLEAN,repetitiontype=OPTIONAL"},`, name)
	case reflect.String:
		return fmt.Sprintf(`{"Tag":"name=%s, type=BYTE_ARRAY, convertedtype=UTF8, repetitiontype=OPTIONAL"},`, name)
	}
	fmt.Println("ELseee caseee")
	return ""
}

func readFromSlice(name string, t reflect.Type, val any) string {

	new_val := val.([]interface{})
	mapData := fmt.Sprintf(`{"Tag":"name=%s, type=LIST,repetitiontype=OPTIONAL","Fields":[`, name)

	for _, v := range new_val {
		mapData += Read("element", reflect.TypeOf(v), v)
		break
	}
	mapData = strings.TrimSuffix(mapData, ",")
	mapData += `]},`
	return mapData
}

func readFromMap(name string, t reflect.Type, val any) string {
	new_val := val.(map[string]interface{})
	mapData := fmt.Sprintf(`{"Tag":"name=%s,repetitiontype=OPTIONAL","Fields":[`, name)

	fmt.Println("parent :", name)
	for key, v := range new_val {
		fmt.Println("    map:", key)
		mapData += Read(key, reflect.TypeOf(v), v)
	}
	mapData = strings.TrimSuffix(mapData, ",")
	mapData += `]},`
	return mapData

}
