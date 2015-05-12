package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/azr/anaconda"
	"github.com/dimroc/urban-events/flagvalidator"
	"github.com/dimroc/urban-events/geoevent"
	"log"
	"net/url"
)

var (
	consumerKey    = flag.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret = flag.String("consumer-secret", "", "Twitter Consumer Key")
	token          = flag.String("token", "", "Twitter Access Token")
	tokenSecret    = flag.String("token-secret", "", "Twitter Token Secret")
	// Defaults to NYC
	locations = flag.String("locations", "-74.3,40.462,-73.65,40.95", "Twitter geographic bounding box")
)

func tweetPusher() chan<- anaconda.Tweet { // return send only channel
	outbox := make(chan anaconda.Tweet)
	go func() {
		for tweet := range outbox {
			g := geoevent.NewFromTweet(tweet)
			jsonOut, err := json.Marshal(g)
			if err == nil {
				fmt.Println(string(jsonOut))
			}
		}
	}()

	return outbox
}

func main() {
	flag.Parse()
	flags := []string{"consumer-key", "consumer-secret", "token", "token-secret"}
	flagvalidator.ValidateFlags(flags)

	anaconda.SetConsumerKey(*consumerKey)
	anaconda.SetConsumerSecret(*consumerSecret)
	api := anaconda.NewTwitterApi(*token, *tokenSecret)

	outbox := tweetPusher()

	v := url.Values{}
	v.Set("locations", *locations)
	stream := api.PublicStreamFilter(v)

	for {
		select {
		case <-stream.Quit:
			log.Println("Quitting")
		case elem := <-stream.C:
			t, ok := elem.(anaconda.Tweet)
			if ok {
				outbox <- t
			}
		}
	}
}
