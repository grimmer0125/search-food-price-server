package util

func GetStringProperty(body map[string]interface{}, key string) string {
	if val, ok := body[key]; ok {
		if val != nil {
			return val.(string)
		}
	}

	return ""
}
