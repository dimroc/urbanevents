package main

import (
	"github.com/dimroc/urban-events/cityservice/cityrecorder"
	"log"
	"os"
)

func main() {
	recorder := cityrecorder.NewTweetRecorder(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_TOKEN"),
		os.Getenv("TWITTER_TOKEN_SECRET"),
	)

	settings, err := cityrecorder.LoadSettings()
	if err != nil {
		log.Fatal(err)
	}

	for _, city := range settings.Cities {
		recorder.Record(city, cityrecorder.StdoutWriter) // Blocking
	}
}
