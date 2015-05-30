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
