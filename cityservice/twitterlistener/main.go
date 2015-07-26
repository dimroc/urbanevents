package main

import (
	"flag"
	"github.com/dimroc/urbanevents/cityservice/citylib"
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

	elastic := citylib.NewElasticConnection(os.Getenv("ELASTICSEARCH_URL"))
	hoodEnricher := citylib.NewHoodEnricher(elastic)
	frenchEnricher := citylib.NewFrenchEnricher()
	gramEnricher := citylib.NewInstagramLinkEnricher(
		os.Getenv("INSTAGRAM_CLIENT_ID"),
		os.Getenv("INSTAGRAM_CLIENT_SECRET"),
		os.Getenv("INSTAGRAM_CLIENT_ACCESS_TOKEN"),
	)

	broadcastEnricher := citylib.NewBroadcastEnricher(hoodEnricher, frenchEnricher, gramEnricher)

	recorder := citylib.NewTwitterRecorder(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_TOKEN"),
		os.Getenv("TWITTER_TOKEN_SECRET"),
		broadcastEnricher,
	)

	settings, err := citylib.LoadSettings(*settingsFilename)
	if err != nil {
		log.Panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	for _, city := range settings.Cities {
		go recorder.Record(city, citylib.StdoutWriter)
	}

	wg.Wait()
}
