package cityrecorder

import (
	"encoding/json"
	"fmt"
)

var (
	StdoutWriter = newStdoutWriter()
)

type Writer interface {
	Write(GeoEvent) error
}

type stdoutWriter struct{}

func newStdoutWriter() *stdoutWriter {
	return &stdoutWriter{}
}

func (w stdoutWriter) Write(g GeoEvent) error {
	jsonOut, err := json.Marshal(g)
	if err == nil {
		fmt.Println(string(jsonOut))
	}

	return err
}
