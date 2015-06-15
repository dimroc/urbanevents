package utils

import (
	"encoding/json"
)

func ToJsonString(value interface{}) (string, error) {
	b, err := json.Marshal(value)
	if err == nil {
		return string(b), err
	} else {
		return "", err
	}
}

func ToJsonStringUnsafe(value interface{}) string {
	rval, err := ToJsonString(value)
	Check(err)
	return rval
}
