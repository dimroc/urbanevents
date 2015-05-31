package utils

import (
	"os"
)

var (
	GO_ENV = getGoEnvironment()
)

func getGoEnvironment() string {
	if len(os.Getenv("GO_ENV")) > 0 {
		return os.Getenv("GO_ENV")
	} else {
		return "development"
	}
}

func GetenvOrDefault(env_var string, default_value string) string {
	temp := os.Getenv(env_var)
	if len(temp) == 0 {
		return default_value
	} else {
		return temp
	}
}
