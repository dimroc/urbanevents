package cityrecorder_test

import (
	"encoding/json"
	"github.com/dimroc/anaconda"
	ig "github.com/dimroc/go-instagram/instagram"
	. "github.com/dimroc/urbanevents/cityservice/cityrecorder"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	Fixture = newFixture()
)

type fixture struct {
	Cities    []City
	GeoEvents []GeoEvent
	Media     []ig.Media
	Tweets    []anaconda.Tweet
}

func (f *fixture) GetCity() City {
	return f.Cities[0]
}

func (f *fixture) GetPoiTweet() anaconda.Tweet {
	return f.Tweets[0]
}

func (f *fixture) GetCoordinateTweet() anaconda.Tweet {
	return f.Tweets[1]
}

func (f *fixture) GetVideoTweet() anaconda.Tweet {
	return f.Tweets[2]
}

func (f *fixture) GetInstagramTweet() anaconda.Tweet {
	return f.Tweets[3]
}

func (f *fixture) GetInstagramMedia() []ig.Media {
	return f.Media
}

func (f *fixture) GetZeroBbGeoEvent() GeoEvent {
	return f.GeoEvents[4]
}

func newFixture() *fixture {
	bounds := [4]float64{-74.3, 40.462, -73.65, 40.95}
	cities := []City{{"nyc", "New York City", []string{}, bounds, PackCircles(bounds)}}

	log.Println("Loading GeoEvents")
	geoevents := []GeoEvent{}
	loadFromFixtureFile("fixtures/geoevents.json", &geoevents)

	log.Println("Loading Tweets")
	tweets := []anaconda.Tweet{}
	loadFromFixtureFile("fixtures/tweets.json", &tweets)

	log.Println("Loading IG Media")
	media := []ig.Media{}
	loadFromFixtureFile("fixtures/media.json", &media)

	return &fixture{
		Cities:    cities,
		GeoEvents: geoevents,
		Media:     media,
		Tweets:    tweets,
	}
}

func loadFromFixtureFile(filename string, v interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Unable to loadFromFixtureFile: ", err)
	}

	jsonErr := json.Unmarshal(data, v)
	if jsonErr != nil {
		log.Fatalln("Unable to marshal in loadFromFixtureFile:", jsonErr)
	}
}

func truncateDocuments() {
	elastic := NewElasticConnection(os.Getenv("ELASTICSEARCH_URL"))
	defer elastic.Connection.Close()
	//indices []string, types []string, args map[string]interface{}, query interface{}

	indices := []string{ES_IndexName}
	types := []string{"geoevent"}
	args := map[string]interface{}{}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	_, err := elastic.Connection.DeleteByQuery(indices, types, args, query)
	Check(err)
}

func setup() {
	ES_IndexName = "test-geoevents-write"
	truncateDocuments()
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	teardown()

	os.Exit(retCode)
}
