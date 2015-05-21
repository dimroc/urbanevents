package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/azr/anaconda"
	"github.com/dimroc/urban-events/cityrecorder/flagvalidator"
	"github.com/dimroc/urban-events/cityrecorder/geoevent"
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

func main() {
	flag.Parse()
	flags := []string{"consumer-key", "consumer-secret", "token", "token-secret"}
	flagvalidator.ValidateFlags(flags)

	anaconda.SetConsumerKey(*consumerKey)
	anaconda.SetConsumerSecret(*consumerSecret)
	api := anaconda.NewTwitterApi(*token, *tokenSecret)

	outbox := tweetWriter()

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

func tweetWriter() chan<- anaconda.Tweet { // return send only channel
	outbox := make(chan anaconda.Tweet)
	go func() {
		for tweet := range outbox {
			g, err := NewFromTweet(tweet)
			if err != nil {
				continue
			}

			jsonOut, err := json.Marshal(g)
			if err == nil {
				fmt.Println(string(jsonOut))
			}
		}
	}()

	return outbox
}

func geoJsonFromPoint(t anaconda.Tweet) geoevent.GeoJson {
	return &geoevent.Point{
		Coordinates: t.Coordinates.Coordinates,
		Type:        t.Coordinates.Type,
	}
}

func geoJsonFromBoundingBox(t anaconda.Tweet) geoevent.GeoJson {
	return &geoevent.BoundingBox{
		Coordinates: t.Place.BoundingBox.Coordinates,
		Type:        t.Place.BoundingBox.Type,
	}
}

func NewFromTweet(t anaconda.Tweet) (*geoevent.GeoEvent, error) {
	if t.Coordinates != nil {
		return &geoevent.GeoEvent{
			Id:           t.Id,
			City:         "nyc",
			GeoJson:      geoJsonFromPoint(t),
			Type:         "tweet",
			Payload:      t.Text,
			LocationType: "coordinate",
		}, nil
	} else if t.Place.PlaceType == "poi" {
		return &geoevent.GeoEvent{
			Id:           t.Id,
			City:         "nyc",
			GeoJson:      geoJsonFromBoundingBox(t),
			Type:         "tweet",
			Payload:      t.Text,
			LocationType: t.Place.PlaceType,
		}, nil
	} else {
		return nil, errors.New("Tweet does not contain a coordinate or place of interest")
	}
}
