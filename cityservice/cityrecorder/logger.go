package cityrecorder

import (
	logging "github.com/op/go-logging"
	"os"
)

var (
	logger = newLogger()
)

func newLogger() *logging.Logger {
	newLogger := logging.MustGetLogger("cityrecorder")
	format := logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{level} ▶ %{shortfunc} %{color:reset}%{message}",
	)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
	return newLogger
}
