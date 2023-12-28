package main

import "strconv"

func mapToStringMap(currentMap map[string]interface{}) map[string]string {
	newMap := map[string]string{}
	for key := range currentMap {
		switch currentMap[key].(type) {
		case string:
			newMap[key] = currentMap[key].(string)
		case int:
			newMap[key] = strconv.Itoa(currentMap[key].(int))
		case bool:
			newMap[key] = strconv.FormatBool(currentMap[key].(bool))
		}
	}
	return newMap
}
