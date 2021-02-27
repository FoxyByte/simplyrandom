package main

func (r *random) jsonMap(json map[string]interface{}, key string) (map[string]interface{}, error) {
	value := json[key]
	if value == nil {
		return nil, errJSON
	}

	newMap, ok := value.(map[string]interface{})
	if !ok {
		return nil, errJSON
	}

	return newMap, nil
}
