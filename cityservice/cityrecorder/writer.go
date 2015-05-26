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

type broadcastWriter struct {
	Writers []Writer
}

func NewBroadcastWriter(writers ...Writer) Writer {
	return &broadcastWriter{Writers: writers}
}

func (b *broadcastWriter) Write(g GeoEvent) error {
	var err error
	for _, writer := range b.Writers {
		newErr := writer.Write(g)
		if newErr != nil {
			err = newErr
		}
	}
	return err
}
