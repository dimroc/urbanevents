package main

import (
	"flag"
	"github.com/dimroc/urbanevents/cityservice/cityrecorder"
	"github.com/dimroc/urbanevents/cityservice/utils"
	"log"
	"os"
	"sync"
)

var (
	settingsFilename = flag.String("settings", "config/nyc.json", "Path to the settings file")
)

func main() {
	flag.Parse()
	utils.ValidateFlags([]string{"settings"})

	elastic := cityrecorder.NewElasticConnection(os.Getenv("ELASTICSEARCH_URL"))
	hoodEnricher := cityrecorder.NewHoodEnricher(elastic)
	recorder := cityrecorder.NewTweetRecorder(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_TOKEN"),
		os.Getenv("TWITTER_TOKEN_SECRET"),
		hoodEnricher,
	)

	settings, err := cityrecorder.LoadSettings(*settingsFilename)
	if err != nil {
		log.Panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	for _, city := range settings.Cities {
		go recorder.Record(city, cityrecorder.StdoutWriter)
	}

	wg.Wait()
}
