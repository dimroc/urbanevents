package cityrecorder_test

import (
	"encoding/json"
	. "github.com/dimroc/urban-events/cityservice/cityrecorder"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
	cities := []City{{"nyc", "New York City", bounds}}
	geoevents := []GeoEvent{}
	loadFromFixtureFile("fixtures/geoevents.json", &geoevents)

	return &fixture{
		Cities:    cities,
		GeoEvents: geoevents,
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

func setup() {
	cmd := exec.Command("GO_ENV=test rake elasticsearch:recreate_index")
	cmd.Run()
	IndexName = "test-ntc-geoevents"
}

func teardown() {
}

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	teardown()

	os.Exit(retCode)
}
