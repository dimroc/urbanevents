package cityrecorder

import (
	"encoding/json"
	"fmt"
	. "github.com/dimroc/urbanevents/cityservice/utils"
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

type logWriter struct{}

func NewLogWriter() Writer {
	return &logWriter{}
}

func (l logWriter) Write(g GeoEvent) error {
	val, err := ToJsonString(g)
	if err == nil {
		Logger.Info(val)
	}

	return err
}

type BroadcastWriter struct {
	Writers []Writer
}

func NewBroadcastWriter(writers ...Writer) *BroadcastWriter {
	return &BroadcastWriter{Writers: writers}
}

func (b *BroadcastWriter) Write(g GeoEvent) error {
	var err error
	for _, writer := range b.Writers {
		newErr := writer.Write(g)
		if newErr != nil {
			Logger.Warning("Encountered Error: %s", newErr)
			err = newErr
		}
	}

	return err
}

func (b *BroadcastWriter) Push(writer Writer) {
	b.Writers = append(b.Writers, writer)
}
