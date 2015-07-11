package cityrecorder

import (
	"errors"
	"fmt"
	"github.com/dimroc/anaconda"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"log"
	"net/url"
	"strings"
)

type TwitterRecorder struct {
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	TokenSecret    string
	Enricher       Enricher
}

type tweetEntry struct {
	Tweet anaconda.Tweet
	City  City
}

func (p *TwitterRecorder) Configured() bool {
	if len(p.ConsumerKey) == 0 || len(p.ConsumerSecret) == 0 || len(p.Token) == 0 || len(p.TokenSecret) == 0 {
		return false
	}

	return true
}

func (p *TwitterRecorder) String() string {
	return fmt.Sprintf("ConsumerKey: %s, ConsumerSecret: %s, Token: %s, TokenSecret: %s", p.ConsumerKey, p.ConsumerSecret, p.Token, p.TokenSecret)
}

func NewTwitterRecorder(consumerKey string, consumerSecret string, token string, tokenSecret string, enricher Enricher) *TwitterRecorder {
	recorder := &TwitterRecorder{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		Token:          token,
		TokenSecret:    tokenSecret,
		Enricher:       enricher,
	}

	if !recorder.Configured() {
		log.Fatal(fmt.Sprintf("Recorder configuration is invalid: %s", recorder))
	}

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	return recorder
}

func (t *TwitterRecorder) Record(city City, writer Writer) {
	api := anaconda.NewTwitterApi(t.Token, t.TokenSecret)
	outbox := t.tweetWriter(writer)

	v := url.Values{}
	v.Set("locations", city.LocationString())
	stream := api.PublicStreamFilter(v)

	Logger.Notice("Listening to tweets from " + city.Key)
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

func (t *TwitterRecorder) tweetWriter(w Writer) chan<- tweetEntry { // return send only channel
	outbox := make(chan tweetEntry)
	go func() {
		for entry := range outbox {
			tweet := entry.Tweet
			city := entry.City

			g, err := NewGeoEventFromTweet(city, tweet)
			if err != nil {
				//Logger.Debug("Unable to create geoevent for city %s from tweet. %s", city.Key, err)
				continue
			}

			err = w.Write(t.Enricher.Enrich(g))
			if err != nil {
				Logger.Warning("Failed to write geoevent: "+g.String(), err)
			}
		}
	}()

	return outbox
}

func geoJsonFromTweet(t anaconda.Tweet) GeoJson {
	if t.Place.PlaceType == "poi" {
		return GeoJsonFrom(t.Place.BoundingBox.Type, t.Place.BoundingBox.Coordinates)
	} else {
		return GeoJsonFrom("point", t.Coordinates.Coordinates)
	}
}

func pointFromTweet(t anaconda.Tweet) ([2]float64, error) {
	if t.Place.PlaceType == "poi" {
		return geoJsonFromTweet(t).Center(), nil
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

func getInstagramUrl(t anaconda.Tweet) string {
	for _, url := range t.Entities.Urls {
		if strings.Contains(url.Url, "instagram") {
			return url.Url
		}
	}

	return ""
}

func getMediaUrl(t anaconda.Tweet) string {
	switch getMediaType(t) {
	case "video":
		return t.ExtendedEntities.Media[0].VideoInfo.Variants[0].Url
	default:
		return getImageUrl(t)
	}
}

func getImageUrl(t anaconda.Tweet) string {
	if len(t.Entities.Media) > 0 && t.Entities.Media[0].Type == "photo" {
		return t.Entities.Media[0].Media_url
	} else {
		return ""
	}
}

func getThumbnailUrl(t anaconda.Tweet) string {
	switch getMediaType(t) {
	case "video":
		return t.ExtendedEntities.Media[0].Media_url
	default:
		imageUrl := getImageUrl(t)
		if len(imageUrl) > 0 {
			return imageUrl + ":thumb"
		} else {
			return ""
		}
	}
}

func getMediaType(t anaconda.Tweet) string {
	if len(t.ExtendedEntities.Media) > 0 &&
		len(t.ExtendedEntities.Media[0].VideoInfo.Variants) > 0 {
		return "video"
	}

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

func getInstagramLinkFromExpandedUrl(t anaconda.Tweet) string {
	if len(t.Entities.Urls) > 0 {
		for _, url := range t.Entities.Urls {
			if strings.Contains(url.Expanded_url, "instagram") {
				return url.Expanded_url
			}
		}
	}

	return ""
}

func NewGeoEventFromTweet(city City, t anaconda.Tweet) (GeoEvent, error) {
	point, err := pointFromTweet(t)

	if err == nil {
		createdAt, _ := t.CreatedAtTime()
		instagramUrl := getInstagramUrl(t)
		if instagramUrl != "" {
			Logger.Warning("INSTAGRAM URL IN TWEET: %s", instagramUrl)
		}

		return GeoEvent{
			CityKey:      city.Key,
			CreatedAt:    createdAt,
			FullName:     t.User.Name,
			GeoJson:      geoJsonFromTweet(t),
			Hashtags:     getHashtagTexts(t),
			Id:           t.IdStr,
			MediaUrl:     getMediaUrl(t),
			Link:         generateLink(t),
			LocationType: getLocationType(t),
			MediaType:    getMediaType(t),
			Text:         t.Text,
			Point:        point,
			Service:      "twitter",
			ThumbnailUrl: getThumbnailUrl(t),
			Type:         "geoevent",
			Username:     t.User.ScreenName,
			Place:        t.Place.Name,
			ExpandedUrl:  getInstagramLinkFromExpandedUrl(t),
		}, nil
	} else {
		return GeoEvent{}, err
	}
}
