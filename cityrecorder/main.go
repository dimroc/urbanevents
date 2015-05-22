package main

import (
	"encoding/json"
	"fmt"
	"github.com/dimroc/urban-events/cityrecorder/cityrecorder"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"os"
)

var (
	settings, settingsErr = cityrecorder.LoadSettings()
)

func stdoutLoggingHandler(h http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, h)
}

func main() {
	if settingsErr != nil {
		log.Fatal(settingsErr)
	}

	fmt.Println("Loaded Settings")

	recorder := cityrecorder.NewTweetRecorder(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_TOKEN"),
		os.Getenv("TWITTER_TOKEN_SECRET"),
	)

	pusher := cityrecorder.NewPusherFromURL(os.Getenv("PUSHER_URL"))

	for _, city := range settings.Cities {
		fmt.Println("Configuring city:", city)
		go recorder.Record(city, pusher)
	}

	router := mux.NewRouter()
	stdChain := alice.New(stdoutLoggingHandler, handlers.CompressHandler) //.Then(finalHandler)

	router.Handle("/api/v1/settings", stdChain.Then(http.HandlerFunc(Settings)))

	fmt.Println("Running server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Settings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(settings); err != nil {
		panic(err)
	}
}
