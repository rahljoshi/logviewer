package parser

import "fmt"

func hasKey(keys map[string]struct{}, key string) bool {
	_, ok := keys[key]
	return ok
}

func stringifyValue(value interface{}) string {
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return fmt.Sprint(value)
	}
}
