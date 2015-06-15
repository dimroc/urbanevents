package main

import (
	"github.com/codegangsta/negroni"
	"github.com/dimroc/urbanevents/cityservice/cityrecorder"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/rs/cors"
	"os"
)

var (
	settingsFilename = GetenvOrDefault("CITYSERVICE_SETTINGS", "config/conf1.json")
	port             = GetenvOrDefault("PORT", "58080")
)

const (
	CTX_SETTINGS_KEY           = "city.settings"
	CTX_ELASTIC_CONNECTION_KEY = "city.elasticconnection"
)

func main() {
	Logger.Info("Running in " + GO_ENV + " with settings " + settingsFilename)
	settings, settingsErr := cityrecorder.LoadSettings(settingsFilename)
	Check(settingsErr)

	recorder := cityrecorder.NewTweetRecorder(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_TOKEN"),
		os.Getenv("TWITTER_TOKEN_SECRET"),
	)

	pusher := cityrecorder.NewPusherFromURL(os.Getenv("PUSHER_URL"))
	elastic := cityrecorder.NewBulkElasticConnection(os.Getenv("ELASTICSEARCH_URL"))
	defer elastic.Close()

	broadcaster := cityrecorder.NewBroadcastWriter(pusher, elastic)
	if GO_ENV == "development" {
		broadcaster.Push(cityrecorder.StdoutWriter)
	}

	for _, city := range settings.Cities {
		Logger.Debug("Configuring city: " + city.String())
		go recorder.Record(city, broadcaster)
	}

	router := mux.NewRouter()
	apiRoutes := router.PathPrefix("/api/v1").Subrouter()
	apiRoutes.HandleFunc("/settings", SettingsHandler).Methods("GET")
	apiRoutes.HandleFunc("/cities", CitiesHandler).Methods("GET")
	apiRoutes.HandleFunc("/cities/{city}", CityHandler).Methods("GET")

	n := negroni.Classic()
	n.Use(cors.Default())
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.Use(SettingsMiddleware(settings))
	n.Use(ElasticMiddleware(elastic))
	n.UseHandler(context.ClearHandler(router))
	n.Run(":" + port)
}
