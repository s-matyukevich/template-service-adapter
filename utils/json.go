package utils

func ConvertToJsonCompatibleMap(obj map[string]interface{}) map[string]interface{} {
	for k, v := range obj {
		obj[k] = convert(v)
	}
	return obj
}
func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}
