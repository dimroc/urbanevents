package utils

import (
	logging "github.com/op/go-logging"
	"log"
	"os"
)

var (
	Logger = newLogger()
)

func newLogger() *logging.Logger {
	newLogger := logging.MustGetLogger("cityrecorder")
	format := logging.MustStringFormatter(
		"%{color}%{level} [%{shortfunc}] %{color:reset}%{message}",
	)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	loggingLevel := os.Getenv("CITYSERVICE_LOGLEVEL")
	if len(loggingLevel) > 0 {
		newLogger.Notice("Configuring log level to %s", loggingLevel)
		// Only errors and more severe messages should be sent to backend1
		backendLeveled := logging.AddModuleLevel(backendFormatter)

		level, err := logging.LogLevel(loggingLevel)
		if err != nil {
			log.Panic(err)
		}

		backendLeveled.SetLevel(level, "")
		logging.SetBackend(backendLeveled)
	} else {
		logging.SetBackend(backendFormatter)
	}

	return newLogger
}

func RequestTracer(method, url, body string) {
	Logger.Debug("Requesting %s %s", method, url)
	Logger.Debug("Request body: %s", body)
}
