package util

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

