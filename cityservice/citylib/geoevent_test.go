package citylib_test

import (
	"github.com/dimroc/urbanevents/cityservice/citylib"
	//. "github.com/dimroc/urbanevents/cityservice/utils"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTryCollapseEmptyBoundingBox(t *testing.T) {
	Convey("Given a GeoEvent with a 0 area bb", t, func() {
		geoevent := Fixture.GetZeroBbGeoEvent()
		So(geoevent.GeoJson.Type, ShouldEqual, "Polygon")

		Convey("it should be converted to a point", func() {
			newgeojson := geoevent.GeoJson.TryCollapseEmptyBoundingBox()
			So(newgeojson.Type, ShouldEqual, "Point")

			shape := newgeojson.GenerateShape()
			point := shape.(*citylib.Point)
			expectation := [2]float64{-74.00432554645808, 40.74185267627071}
			So(point.Coordinates, ShouldEqual, expectation)
		})
	})
}
