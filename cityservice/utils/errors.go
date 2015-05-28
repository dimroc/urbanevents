package utils

import (
	"fmt"
)

func Check(e error) {
	if e != nil {
		Logger.Critical(fmt.Sprint(e))
		panic(e)
	}
}
