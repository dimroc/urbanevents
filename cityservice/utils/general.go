package utils

import (
	"encoding/json"
	"os"
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

func GetBaseUrl() string {
	baseUrl := os.Getenv("BASEURL")
	if len(baseUrl) == 0 {
		baseUrl = os.Getenv("TUTUM_SERVICE_FQDN")
	}

	return baseUrl
}
