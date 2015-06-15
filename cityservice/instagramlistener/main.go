package main

import (
	"flag"
	"github.com/dimroc/urbanevents/cityservice/cityrecorder"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	settingsFilename = flag.String("settings", "", "Path to the settings file")
)

func main() {
	flag.Parse()
	ValidateFlags([]string{"settings"})
	recorder := cityrecorder.NewInstagramRecorder(
		os.Getenv("INSTAGRAM_CLIENT_ID"),
		os.Getenv("INSTAGRAM_CLIENT_SECRET"),
		cityrecorder.StdoutWriter,
	)

	defer recorder.Close()
	recorder.DeleteAllSubscriptions()
	printExistingSubscriptions(recorder)

	settings, err := cityrecorder.LoadSettings(*settingsFilename)
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
	recorder.Subscribe(settings.Cities)

	printExistingSubscriptions(recorder)
	wg.Wait()
}

func printExistingSubscriptions(recorder *cityrecorder.InstagramRecorder) {
	Logger.Debug("Existing Subscriptions:")
	subscriptions := recorder.GetSubscriptions()
	for _, subscription := range subscriptions {
		Logger.Debug(ToJsonStringUnsafe(subscription))
	}
}
