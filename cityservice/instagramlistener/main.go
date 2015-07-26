package main

import (
	"flag"
	"github.com/dimroc/urbanevents/cityservice/citylib"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	settingsFilename = flag.String("settings", "config/nyc.json", "Path to the settings file")
	baseUrl          = flag.String("baseurl", "", "The base url of the service used for callbacks")
)

func main() {
	flag.Parse()
	ValidateFlags([]string{"settings", "baseurl"})
	elastic := citylib.NewElasticConnection(os.Getenv("ELASTICSEARCH_URL"))
	hoodEnricher := citylib.NewHoodEnricher(elastic)
	recorder := citylib.NewInstagramRecorder(
		os.Getenv("INSTAGRAM_CLIENT_ID"),
		os.Getenv("INSTAGRAM_CLIENT_SECRET"),
		citylib.StdoutWriter,
		hoodEnricher,
	)

	defer recorder.Close()
	recorder.DeleteAllSubscriptions()
	printExistingSubscriptions(recorder)

	settings, err := citylib.LoadSettings(*settingsFilename)
	if err != nil {
		log.Panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	Logger.Debug("Listening for instagram events")
	router := mux.NewRouter()
	router.Handle("/api/v1/callbacks/instagram/{city}", recorder).Methods("GET", "POST")
	go func() { log.Fatal(http.ListenAndServe(":58080", router)) }()

	Logger.Debug("Adding Subscriptions if necessary")
	recorder.Subscribe(*baseUrl+"/api/v1/callbacks/instagram/", settings.Cities)

	printExistingSubscriptions(recorder)
	wg.Wait()
}

func printExistingSubscriptions(recorder *citylib.InstagramRecorder) {
	Logger.Debug("Existing Subscriptions:")
	subscriptions := recorder.GetSubscriptions()
	for _, subscription := range subscriptions {
		Logger.Debug(ToJsonStringUnsafe(subscription))
	}
}
