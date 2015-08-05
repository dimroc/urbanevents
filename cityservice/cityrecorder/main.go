package main

import (
	"github.com/codegangsta/negroni"
	"github.com/dimroc/urbanevents/cityservice/citylib"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"os"
	"time"
)

var (
	settingsFilename = GetenvOrDefault("CITYSERVICE_SETTINGS", "../config/nyc.json")
	enableInstagram  = GetenvOrDefault("CITYSERVICE_INSTAGRAM", "true")
	enableTwitter    = GetenvOrDefault("CITYSERVICE_TWITTER", "true")
	port             = GetenvOrDefault("PORT", "58080")
)

func main() {
	Logger.Info("Running in " + GO_ENV + " with settings " + settingsFilename + " and base url: " + GetBaseUrl())
	settings, settingsErr := citylib.LoadSettings(settingsFilename)
	Check(settingsErr)

	// Configure Geoevent Writers
	eventpusher := citylib.NewEventPusher()
	elastic := citylib.NewElasticConnection(os.Getenv("ELASTICSEARCH_URL"))
	hoodEnricher := citylib.NewHoodEnricher(elastic)
	frenchEnricher := citylib.NewFrenchEnricher()
	gramEnricher := citylib.NewInstagramLinkEnricher(
		os.Getenv("INSTAGRAM_CLIENT_ID"),
		os.Getenv("INSTAGRAM_CLIENT_SECRET"),
		os.Getenv("INSTAGRAM_CLIENT_ACCESS_TOKEN"),
	)
	broadcastEnricher := citylib.NewBroadcastEnricher(hoodEnricher, frenchEnricher, gramEnricher)

	defer eventpusher.Close()
	defer elastic.Close()

	tweetRecorder := citylib.NewTwitterRecorder(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_TOKEN"),
		os.Getenv("TWITTER_TOKEN_SECRET"),
		broadcastEnricher,
	)

	logWriter := citylib.NewLogWriter()
	broadcaster := citylib.NewBroadcastWriter(eventpusher, elastic, logWriter)

	instagramRecorder := citylib.NewInstagramRecorder(
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
	apiRoutes.HandleFunc("/settings", citylib.SettingsHandler).Methods("GET")
	apiRoutes.HandleFunc("/cities", citylib.CitiesHandler).Methods("GET")
	apiRoutes.HandleFunc("/cities/{city}", citylib.CityHandler).Methods("GET")
	apiRoutes.Handle("/callbacks/instagram/{city}", instagramRecorder).Methods("GET", "POST")

	if instagramEnabled() && len(GetBaseUrl()) > 0 {
		timer := time.NewTimer(time.Second)
		go func() {
			<-timer.C
			Logger.Debug("Subscribing to instagram geographies")

			instagramRecorder.DeleteAllSubscriptions()
			instagramRecorder.Subscribe(GetBaseUrl()+"/api/v1/callbacks/instagram/", settings.Cities)
		}()
	} else if len(GetBaseUrl()) == 0 {
		Logger.Warning("Unable to subscribe to instagram, no base url for callback")
	}

	n := negroni.Classic()
	n.Use(cors.Default())
	n.Use(citylib.SettingsMiddleware(settings))
	n.Use(citylib.ElasticMiddleware(elastic))
	n.UseHandler(context.ClearHandler(router))
	n.Run(":" + port)
}

func instagramEnabled() bool {
	return enableInstagram != "false"
}

func twitterEnabled() bool {
	return enableTwitter != "false"
}
