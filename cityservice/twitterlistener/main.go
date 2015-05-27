package main

import (
	"github.com/dimroc/urban-events/cityservice/cityrecorder"
	"log"
	"os"
	"sync"
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
		log.Panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	for _, city := range settings.Cities {
		go recorder.Record(city, cityrecorder.StdoutWriter) // Blocking
	}

	wg.Wait()
}
