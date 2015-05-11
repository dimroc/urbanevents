package main

import (
	"flag"
	"fmt"
	"github.com/azr/anaconda"
	"log"
	"net/url"
	//"io/ioutil"
	//"net/http"
)

var (
	consumerKey    = flag.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret = flag.String("consumer-secret", "", "Twitter Consumer Key")
	token          = flag.String("token", "", "Twitter Consumer Key")
	tokenSecret    = flag.String("token-secret", "", "Twitter Consumer Key")
)

func main() {
	flag.Parse()

	anaconda.SetConsumerKey(*consumerKey)
	anaconda.SetConsumerSecret(*consumerSecret)
	api := anaconda.NewTwitterApi(*token, *tokenSecret)
	api.SetLogger(anaconda.BasicLogger)

	searchResult, err := api.GetSearch("twitter", nil)
	if err != nil {
		log.Fatal(err)
	}

	// works
	for _, tweet := range searchResult.Statuses {
		fmt.Println(tweet.Text)
	}

	v := url.Values{}
	v.Set("track", "tweet")
	// Fails: Twitter streaming: leaving after an irremediable error: [401 Authorization Required]
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
