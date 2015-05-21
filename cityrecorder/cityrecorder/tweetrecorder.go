package cityrecorder

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/azr/anaconda"
	"log"
	"net/url"
)

type TweetRecorder struct {
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	TokenSecret    string
}

func (t *TweetRecorder) Start(locations string) {
	anaconda.SetConsumerKey(t.ConsumerKey)
	anaconda.SetConsumerSecret(t.ConsumerSecret)
	api := anaconda.NewTwitterApi(t.Token, t.TokenSecret)

	outbox := tweetWriter()

	v := url.Values{}
	v.Set("locations", locations)
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
			g, err := newFromTweet(tweet)
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

func geoJsonFromPoint(t anaconda.Tweet) GeoJson {
	return &Point{
		Coordinates: t.Coordinates.Coordinates,
		Type:        t.Coordinates.Type,
	}
}

func geoJsonFromBoundingBox(t anaconda.Tweet) GeoJson {
	return &BoundingBox{
		Coordinates: t.Place.BoundingBox.Coordinates,
		Type:        t.Place.BoundingBox.Type,
	}
}

func newFromTweet(t anaconda.Tweet) (*GeoEvent, error) {
	if t.Coordinates != nil {
		return &GeoEvent{
			Id:           t.Id,
			City:         "nyc",
			GeoJson:      geoJsonFromPoint(t),
			Type:         "tweet",
			Payload:      t.Text,
			LocationType: "coordinate",
		}, nil
	} else if t.Place.PlaceType == "poi" {
		return &GeoEvent{
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
