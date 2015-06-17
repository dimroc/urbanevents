package cityrecorder_test

import (
	"encoding/json"
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
}

func (f *fixture) GetCity() City {
	return f.Cities[0]
}

func newFixture() *fixture {
	bounds := [][]float64{[]float64{-74.3, 40.462}, []float64{-73.65, 40.95}}
	cities := []City{{"nyc", "New York City", []string{}, bounds}}
	geoevents := []GeoEvent{}
	loadFromFixtureFile("fixtures/geoevents.json", &geoevents)

	media := []ig.Media{}
	loadFromFixtureFile("fixtures/media.json", &media)

	return &fixture{
		Cities:    cities,
		GeoEvents: geoevents,
		Media:     media,
	}
}

func loadFromFixtureFile(filename string, v interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	jsonErr := json.Unmarshal(data, v)
	if jsonErr != nil {
		log.Fatalln("error:", jsonErr)
	}
}

func truncateDocuments() {
	elastic := NewElasticConnection(os.Getenv("ELASTICSEARCH_URL"))
	defer elastic.Connection.Close()
	//indices []string, types []string, args map[string]interface{}, query interface{}

	indices := []string{IndexName}
	types := []string{"tweet"}
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
	IndexName = "test-geoevents"
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
