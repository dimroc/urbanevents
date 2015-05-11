package main

import (
	"flag"
	"fmt"
	"github.com/azr/anaconda"
	"github.com/dimroc/geo-twitter-listener/flagvalidator"
	"log"
	"net/url"
)

var (
	consumerKey    = flag.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret = flag.String("consumer-secret", "", "Twitter Consumer Key")
	token          = flag.String("token", "", "Twitter Consumer Key")
	tokenSecret    = flag.String("token-secret", "", "Twitter Consumer Key")
	// Defaults to NYC
	locations = flag.String("locations", "-74.3,40.462,-73.65,40.95", "Twitter geographic bounding box")
)

func main() {
	flag.Parse()
	flags := []string{"consumer-key", "consumer-secret", "token", "token-secret"}
	flagvalidator.ValidateFlags(flags)

	anaconda.SetConsumerKey(*consumerKey)
	anaconda.SetConsumerSecret(*consumerSecret)
	api := anaconda.NewTwitterApi(*token, *tokenSecret)
	api.SetLogger(anaconda.BasicLogger)

	v := url.Values{}
	v.Set("locations", *locations)
	stream := api.PublicStreamFilter(v)

	for {
		select {
		case <-stream.Quit:
			log.Println("Quitting")
		case elem := <-stream.C:
			fmt.Println(elem)
		}
	}
}
