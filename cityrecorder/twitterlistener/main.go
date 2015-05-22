package main

import (
	"flag"
	"github.com/dimroc/urban-events/cityrecorder/cityrecorder"
	"github.com/dimroc/urban-events/cityrecorder/flagvalidator"
	"log"
)

var (
	consumerKey    = flag.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret = flag.String("consumer-secret", "", "Twitter Consumer Key")
	token          = flag.String("token", "", "Twitter Access Token")
	tokenSecret    = flag.String("token-secret", "", "Twitter Token Secret")
)

func main() {
	flag.Parse()
	flags := []string{"consumer-key", "consumer-secret", "token", "token-secret"}
	flagvalidator.ValidateFlags(flags)

	recorder := cityrecorder.TweetRecorder{
		ConsumerKey:    *consumerKey,
		ConsumerSecret: *consumerSecret,
		Token:          *token,
		TokenSecret:    *tokenSecret,
	}

	//city := cityrecorder.City{
	//Key:       "nyc",
	//Display:   "New York City",
	//Locations: "-74.3,40.462,-73.65,40.95",
	//}

	//cities := []cityrecorder.City{city}
	//settings := cityrecorder.Settings{
	//Cities: cities,
	//}

	settings, err := cityrecorder.LoadSettings()
	if err != nil {
		log.Fatal(err)
	}

	for _, city := range settings.Cities {
		recorder.Start(city)
	}
}
