package cityrecorder

import (
	"errors"
	"fmt"
	"github.com/azr/anaconda"
	. "github.com/dimroc/urbanevents/cityservice/utils"
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

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	return recorder
}

func (t *TweetRecorder) Record(city City, writer Writer) {
	api := anaconda.NewTwitterApi(t.Token, t.TokenSecret)
	outbox := tweetWriter(writer)

	v := url.Values{}
	v.Set("locations", city.LocationString())
	stream := api.PublicStreamFilter(v)

	Logger.Debug("Listening to tweets from " + city.Key)
	for {
		select {
		case <-stream.Quit:
			Logger.Debug("%s Stream Quit", city.Key)
		case elem := <-stream.C:
			t, ok := elem.(anaconda.Tweet)
			if ok {
				outbox <- tweetEntry{Tweet: t, City: city}
			} else {
				log.Panic("Unable to type cast tweet")
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
				//Logger.Debug("Unable to create geoevent for city %s from tweet. %s", city.Key, err)
				continue
			}

			err = w.Write(g)
			if err != nil {
				Logger.Warning("Failed to write geoevent: "+g.String(), err)
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

type localEntities anaconda.Entities

func (e localEntities) GetHashtagTexts() []string {
	texts := make([]string, len(e.Hashtags))
	for index, hashtag := range e.Hashtags {
		texts[index] = hashtag.Text
	}

	return texts
}

func (e localEntities) GetMedias() ([]string, []string, []string) {
	types := make([]string, len(e.Media))
	mediaUrls := make([]string, len(e.Media))
	expandedUrls := make([]string, len(e.Media))

	for index, media := range e.Media {
		types[index] = media.Type
		mediaUrls[index] = media.Media_url
		expandedUrls[index] = media.Expanded_url
	}

	return types, mediaUrls, expandedUrls
}

func metadataFromTweet(t anaconda.Tweet) Metadata {
	entities := localEntities(t.Entities)
	types, mediaUrls, expandedUrls := entities.GetMedias()

	return Tweet{
		ScreenName:   t.User.ScreenName,
		Hashtags:     entities.GetHashtagTexts(),
		MediaTypes:   types,
		MediaUrls:    mediaUrls,
		ExpandedUrls: expandedUrls,
	}
}

func newFromTweet(city City, t anaconda.Tweet) (GeoEvent, error) {
	if t.Coordinates != nil {
		createdAt, _ := t.CreatedAtTime()
		return GeoEvent{
			Id:           t.IdStr,
			CreatedAt:    createdAt,
			CityKey:      city.Key,
			GeoJson:      nil,
			Point:        t.Coordinates.Coordinates,
			Type:         "geoevent",
			Service:      "twitter",
			Payload:      t.Text,
			Metadata:     metadataFromTweet(t),
			LocationType: "coordinate",
		}, nil
	} else if t.Place.PlaceType == "poi" {
		createdAt, _ := t.CreatedAtTime()
		geoJson := geoJsonFromBoundingBox(t)
		return GeoEvent{
			Id:           t.IdStr,
			CreatedAt:    createdAt,
			CityKey:      city.Key,
			GeoJson:      geoJson,
			Point:        geoJson.Center(),
			Type:         "geoevent",
			Service:      "twitter",
			Payload:      t.Text,
			Metadata:     metadataFromTweet(t),
			LocationType: t.Place.PlaceType,
		}, nil
	} else {
		return GeoEvent{}, errors.New("Tweet does not contain a coordinate or place of interest")
	}
}
