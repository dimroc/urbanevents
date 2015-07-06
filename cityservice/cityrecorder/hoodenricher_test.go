package cityrecorder_test

import (
	"github.com/dimroc/urbanevents/cityservice/cityrecorder"
	//. "github.com/dimroc/urbanevents/cityservice/utils"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestEnrich(t *testing.T) {
	Convey("Given a GeoEvent and an Elastic Connection", t, func() {
		geoevent := Fixture.GeoEvents[0]
		elastic := cityrecorder.NewElasticConnection(os.Getenv("ELASTICSEARCH_URL"))
		//elastic.SetRequestTracer(RequestTracer)
		defer elastic.Connection.Close()
		hoodEnricher := cityrecorder.NewHoodEnricher(elastic)

		Convey("there should be returned neighborhoods", func() {
			newGeo := hoodEnricher.Enrich(geoevent)

			So(geoevent.Neighborhoods, ShouldBeEmpty)
			So(newGeo.Neighborhoods, ShouldContain, "Manhattan")
		})
	})
}
