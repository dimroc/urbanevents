package citylib

import (
	"os"
)

const (
	CTX_SETTINGS_KEY                  = "city.settings"
	CTX_ELASTIC_CONNECTION_KEY        = "city.elasticconnection"
	ES_TypeName                string = "geoevent"
)

var (
	ES_IndexName = os.Getenv("GO_ENV") + "-geoevents-write"
)
