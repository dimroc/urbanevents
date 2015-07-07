package utils

import (
	"encoding/json"
	"fmt"
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
	if len(baseUrl) == 0 && len(os.Getenv("TUTUM_SERVICE_FQDN")) > 0 {
		baseUrl = fmt.Sprintf("http://%s", os.Getenv("TUTUM_SERVICE_FQDN"))
		if len(os.Getenv("PORT")) > 0 {
			baseUrl = fmt.Sprintf("%s:%s", baseUrl, os.Getenv("PORT"))
		}
	}

	return baseUrl
}
