package main

import (
	"flag"
	"github.com/dimroc/urban-events/cityrecorder/cityrecorder"
	"github.com/dimroc/urban-events/cityrecorder/flagvalidator"
)

var (
	consumerKey    = flag.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret = flag.String("consumer-secret", "", "Twitter Consumer Key")
	token          = flag.String("token", "", "Twitter Access Token")
	tokenSecret    = flag.String("token-secret", "", "Twitter Token Secret")
	// Defaults to NYC
	locations = flag.String("locations", "-74.3,40.462,-73.65,40.95", "Twitter geographic bounding box")
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

	recorder.Start(*locations)
}
