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

func geoJsonFromBoundingBox(t anaconda.Tweet) GeoJson {
	if t.Place.PlaceType == "poi" {
		return &BoundingBox{
			Coordinates: t.Place.BoundingBox.Coordinates,
			Type:        t.Place.BoundingBox.Type,
		}
	} else {
		return nil
	}
}

func pointFromTweet(t anaconda.Tweet) ([2]float64, error) {
	if t.Place.PlaceType == "poi" {
		return geoJsonFromBoundingBox(t).Center(), nil
	} else if t.Coordinates != nil {
		return t.Coordinates.Coordinates, nil
	} else {
		return [2]float64{}, errors.New("No coordinate for tweet")
	}
}

func getHashtagTexts(t anaconda.Tweet) []string {
	e := t.Entities
	texts := make([]string, len(e.Hashtags))
	for index, hashtag := range e.Hashtags {
		texts[index] = hashtag.Text
	}

	return texts
}

func getImageUrl(t anaconda.Tweet) string {
	if len(t.Entities.Media) > 0 && t.Entities.Media[0].Type == "photo" {
		return t.Entities.Media[0].Media_url
	} else {
		return ""
	}
}

func getThumbnailUrl(t anaconda.Tweet) string {
	imageUrl := getImageUrl(t)
	if len(imageUrl) > 0 {
		return imageUrl + ":thumb"
	} else {
		return ""
	}
}

func getVideoUrl(t anaconda.Tweet) string {
	return "" // Doesn't exist yet.
}

func getMediaType(t anaconda.Tweet) string {
	if len(t.Entities.Media) > 0 {
		current := t.Entities.Media[0].Type
		if current == "photo" {
			return "image"
		} else if len(current) > 0 {
			return current
		} else {
			return "text"
		}
	} else {
		return "text"
	}
}

func generateLink(t anaconda.Tweet) string {
	// https://twitter.com/thereaIbanksy/status/613791445858648064
	return fmt.Sprintf("https://twitter.com/%s/status/%s", t.User.ScreenName, t.IdStr)
}

func getLocationType(t anaconda.Tweet) string {
	if t.Place.PlaceType == "poi" {
		return "poi"
	} else {
		return "coordinate"
	}
}

func metadataFromTweet(t anaconda.Tweet) (metadata, error) {
	point, err := pointFromTweet(t)

	if err != nil {
		return metadata{}, err
	} else {
		// Leaned towards readability over speed, so some
		// methods perform redundant actions.
		return metadata{
			GeoJson:      geoJsonFromBoundingBox(t),
			Hashtags:     getHashtagTexts(t),
			ImageUrl:     getImageUrl(t),
			Link:         generateLink(t),
			LocationType: getLocationType(t),
			MediaType:    getMediaType(t),
			Point:        point,
			ThumbnailUrl: getThumbnailUrl(t),
			VideoUrl:     getVideoUrl(t),
		}, nil
	}
}

func newFromTweet(city City, t anaconda.Tweet) (GeoEvent, error) {
	metadata, err := metadataFromTweet(t)

	if err == nil {
		createdAt, _ := t.CreatedAtTime()
		return GeoEvent{
			CityKey:      city.Key,
			CreatedAt:    createdAt,
			FullName:     t.User.Name,
			GeoJson:      metadata.GeoJson,
			Hashtags:     metadata.Hashtags,
			Id:           t.IdStr,
			ImageUrl:     metadata.ImageUrl,
			Link:         metadata.Link,
			LocationType: metadata.LocationType,
			MediaType:    metadata.MediaType,
			Payload:      t.Text,
			Point:        metadata.Point,
			Service:      "twitter",
			ThumbnailUrl: metadata.ThumbnailUrl,
			Type:         "geoevent",
			Username:     t.User.ScreenName,
			VideoUrl:     metadata.VideoUrl,
		}, nil
	} else {
		return GeoEvent{}, err
	}
}

type metadata struct {
	GeoJson      GeoJson
	Hashtags     []string
	ImageUrl     string
	Link         string
	LocationType string
	MediaType    string
	Point        [2]float64
	ThumbnailUrl string
	VideoUrl     string
}
