package cityrecorder

import (
	"encoding/json"
	"fmt"
)

var (
	StdoutWriter = NewStdoutWriter()
)

type Writer interface {
	Write(GeoEvent) error
}

type stdoutWriter struct{}

func NewStdoutWriter() *stdoutWriter {
	return &stdoutWriter{}
}

func (w stdoutWriter) Write(g GeoEvent) error {
	jsonOut, err := json.Marshal(g)
	if err == nil {
		fmt.Println(string(jsonOut))
	}

	return err
}
