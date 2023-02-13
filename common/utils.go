package common

import "encoding/json"

func ToMap(i interface{}) (map[string]interface{}, error) {

	bytes, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ToMapString(i interface{}) (string, error) {

	bytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
