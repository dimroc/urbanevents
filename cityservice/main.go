package main

import (
	"github.com/codegangsta/negroni"
	"github.com/dimroc/urbanevents/cityservice/cityrecorder"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"os"
	"time"
)

var (
	settingsFilename = GetenvOrDefault("CITYSERVICE_SETTINGS", "config/nyc.json")
	enableInstagram  = GetenvOrDefault("CITYSERVICE_INSTAGRAM", "true")
	enableTwitter    = GetenvOrDefault("CITYSERVICE_TWITTER", "true")
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

	// Configure Geoevent Writers
	eventpusher := cityrecorder.NewEventPusher()
	elastic := cityrecorder.NewBulkElasticConnection(os.Getenv("ELASTICSEARCH_URL"))
	hoodEnricher := cityrecorder.NewHoodEnricher(elastic)
	defer eventpusher.Close()
	defer elastic.Close()

	tweetRecorder := cityrecorder.NewTweetRecorder(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_TOKEN"),
		os.Getenv("TWITTER_TOKEN_SECRET"),
		hoodEnricher,
	)

	broadcaster := cityrecorder.NewBroadcastWriter(eventpusher, elastic)
	if GO_ENV == "development" {
		broadcaster.Push(cityrecorder.StdoutWriter)
	}

	instagramRecorder := cityrecorder.NewInstagramRecorder(
		os.Getenv("INSTAGRAM_CLIENT_ID"),
		os.Getenv("INSTAGRAM_CLIENT_SECRET"),
		broadcaster,
		hoodEnricher,
	)
	defer instagramRecorder.Close()

	if twitterEnabled() {
		for _, city := range settings.Cities {
			Logger.Debug("Recording tweets for city: " + city.String())
			go tweetRecorder.Record(city, broadcaster)
		}
	}

	router := mux.NewRouter()
	apiRoutes := router.PathPrefix("/api/v1").Subrouter()
	apiRoutes.Handle("/events", eventpusher)
	apiRoutes.HandleFunc("/settings", SettingsHandler).Methods("GET")
	apiRoutes.HandleFunc("/cities", CitiesHandler).Methods("GET")
	apiRoutes.HandleFunc("/cities/{city}", CityHandler).Methods("GET")
	apiRoutes.Handle("/callbacks/instagram/{city}", instagramRecorder).Methods("GET", "POST")

	if instagramEnabled() {
		timer := time.NewTimer(time.Second)
		go func() {
			<-timer.C
			Logger.Debug("Subscribing to instagram geographies")

			instagramRecorder.DeleteAllSubscriptions()
			instagramRecorder.Subscribe(GetBaseUrl()+"/api/v1/callbacks/instagram/", settings.Cities)
		}()
	}

	n := negroni.Classic()
	n.Use(cors.Default())
	n.Use(SettingsMiddleware(settings))
	n.Use(ElasticMiddleware(elastic))
	n.UseHandler(context.ClearHandler(router))
	n.Run(":" + port)
}

func instagramEnabled() bool {
	return enableInstagram != "false"
}

func twitterEnabled() bool {
	return enableTwitter != "false"
}
