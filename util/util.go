package util

import "time"

func NowMilliseconds() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}

func MergeCopy(source,target map[string]interface{}) map[string]interface{} {
	if target == nil {
		target = make(map[string]interface{})
	}
	if source == nil {
		return target
	}
	for key, ele := range DeepCopy(source) {
		target[key] = ele
	}
	return target
}

func DeepCopy(value map[string]interface{}) map[string]interface{} {
	newCopy := deepCopy(value)
	if newMap, ok := newCopy.(map[string]interface{}); ok {
		return newMap
	}
	return nil
}

func deepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = deepCopy(v)
		}
		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = deepCopy(v)
		}
		return newSlice
	}
	return value
}

