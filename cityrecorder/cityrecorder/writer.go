package cityrecorder

import (
	"encoding/json"
	"fmt"
)

type Writer interface {
	Write(GeoEvent) error
}

type StdoutWriter struct{}

func (w StdoutWriter) Write(g GeoEvent) error {
	jsonOut, err := json.Marshal(g)
	if err == nil {
		fmt.Println(string(jsonOut))
	}

	return err
}
