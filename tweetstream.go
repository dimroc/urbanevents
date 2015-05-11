package main

import (
	"flag"
	"fmt"
	"github.com/nhjk/tweetstream"
	"log"
	//"net/url"
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

	log.Println(*consumerKey)
	log.Println(*consumerSecret)
	log.Println(*token)
	log.Println(*tokenSecret)
	// Add your credentials
	tweetstream.SetCredentials(*consumerKey, *consumerSecret, *token, *tokenSecret)

	// Add optional parameters. If none are set, it defaults to a sample stream.
	tweetstream.Track = []string{"golang", "gophers", "twitter"}

	// Stream.
	tweets, err := tweetstream.Stream()
	if err != nil {
		log.Fatal(err)
	}

	for tweet := range tweets {
		fmt.Println(tweet.User.Name, "tweeted", tweet.Text)
	}
}
