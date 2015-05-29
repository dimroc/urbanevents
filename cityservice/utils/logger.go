package utils

import (
	logging "github.com/op/go-logging"
	"os"
)

var (
	Logger = newLogger()
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

func RequestTracer(method, url, body string) {
	Logger.Debug("Requesting %s %s", method, url)
	Logger.Debug("Request body: %s", body)
}
