/*
Basic json-schema generator based on Go types, for easy interchange of Go
structures between diferent languages.
*/
package datatype

import (
	"fmt"
	"strings"
)

func BuildMedadata(data map[string]interface{}) string {
	metadata := (`{"Tag":"name=parquet-go","Fields":[`)
	for name, val := range data {
		mapVal := val.(map[string]interface{})
		if val == nil {
			metadata += Read(name, mapVal)
		} else {
			metadata += Read(name, mapVal)
		}
	}
	metadata = strings.TrimSuffix(metadata, ",")
	metadata += `]}`
	return metadata
}

func Read(name string, val map[string]interface{}) string {
	switch val["type"] {
	case "LIST":
		return readFromSlice(name, val["value"].(map[string]interface{}))
	case "MAP":
		return readFromMap(name, val["value"].(map[string]interface{}))
	case "INT":
		return fmt.Sprintf(`{"Tag":"name=%s, type=INT64,repetitiontype=OPTIONAL"},`, name)
	case "FLOAT":
		return fmt.Sprintf(`{"Tag":"name=%s, type=FLOAT,repetitiontype=OPTIONAL"},`, name)
	case "BOOLEAN":
		return fmt.Sprintf(`{"Tag":"name=%s, type=BOOLEAN,repetitiontype=OPTIONAL"},`, name)
	case "BYTE_ARRAY":
		return fmt.Sprintf(`{"Tag":"name=%s, type=BYTE_ARRAY, convertedtype=UTF8, repetitiontype=OPTIONAL"},`, name)
	}

	fmt.Println("ELseee caseee", val["type"])
	return ""
}

func readFromSlice(name string, val map[string]interface{}) string {

	mapData := fmt.Sprintf(`{"Tag":"name=%s, type=LIST,repetitiontype=OPTIONAL","Fields":[`, name)
	mapData += Read("element", val)
	mapData = strings.TrimSuffix(mapData, ",")
	mapData += `]},`
	return mapData
}

func readFromMap(name string, val map[string]interface{}) string {
	mapData := fmt.Sprintf(`{"Tag":"name=%s,repetitiontype=OPTIONAL","Fields":[`, name)
	for key, v := range val {
		mapV := v.(map[string]interface{})
		mapData += Read(key, mapV)
	}
	mapData = strings.TrimSuffix(mapData, ",")
	mapData += `]},`
	return mapData

}
