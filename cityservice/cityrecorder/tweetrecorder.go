package cityrecorder

import (
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

type tweetEntry struct {
	Tweet anaconda.Tweet
	City  City
}

func (p *TweetRecorder) Configured() bool {
	if len(p.ConsumerKey) == 0 || len(p.ConsumerSecret) == 0 || len(p.Token) == 0 || len(p.TokenSecret) == 0 {
		return false
	}

	return true
}

func (p *TweetRecorder) String() string {
	return fmt.Sprintf("ConsumerKey: %s, ConsumerSecret: %s, Token: %s, TokenSecret: %s", p.ConsumerKey, p.ConsumerSecret, p.Token, p.TokenSecret)
}

func NewTweetRecorder(consumerKey string, consumerSecret string, token string, tokenSecret string) *TweetRecorder {
	recorder := &TweetRecorder{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		Token:          token,
		TokenSecret:    tokenSecret,
	}

	if !recorder.Configured() {
		log.Fatal(fmt.Sprintf("Recorder configuration is invalid: %s", recorder))
	}

	return recorder
}

func (t *TweetRecorder) Record(city City, writer Writer) {
	anaconda.SetConsumerKey(t.ConsumerKey)
	anaconda.SetConsumerSecret(t.ConsumerSecret)
	api := anaconda.NewTwitterApi(t.Token, t.TokenSecret)

	outbox := tweetWriter(writer)

	v := url.Values{}
	v.Set("locations", city.LocationString())
	stream := api.PublicStreamFilter(v)

	for {
		select {
		case <-stream.Quit:
			log.Println("Quitting")
		case elem := <-stream.C:
			t, ok := elem.(anaconda.Tweet)
			if ok {
				outbox <- tweetEntry{Tweet: t, City: city}
			}
		}
	}
}

func tweetWriter(w Writer) chan<- tweetEntry { // return send only channel
	outbox := make(chan tweetEntry)
	go func() {
		for entry := range outbox {
			tweet := entry.Tweet
			city := entry.City

			g, err := newFromTweet(city, tweet)
			if err != nil {
				continue
			}

			err = w.Write(g)
			if err != nil {
				log.Fatal(err)
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

func newFromTweet(city City, t anaconda.Tweet) (GeoEvent, error) {
	if t.Coordinates != nil {
		return GeoEvent{
			Id:           t.Id,
			CityKey:      city.Key,
			GeoJson:      geoJsonFromPoint(t),
			Type:         "tweet",
			Payload:      t.Text,
			LocationType: "coordinate",
		}, nil
	} else if t.Place.PlaceType == "poi" {
		return GeoEvent{
			Id:           t.Id,
			CityKey:      city.Key,
			GeoJson:      geoJsonFromBoundingBox(t),
			Type:         "tweet",
			Payload:      t.Text,
			LocationType: t.Place.PlaceType,
		}, nil
	} else {
		return GeoEvent{}, errors.New("Tweet does not contain a coordinate or place of interest")
	}
}
