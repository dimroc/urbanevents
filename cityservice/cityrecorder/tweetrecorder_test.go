package cityrecorder_test

import (
	"github.com/dimroc/urbanevents/cityservice/cityrecorder"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewGeoEventFromTweet(t *testing.T) {
	Convey("Given a city and an anaconda.Tweet with a poi", t, func() {
		city := Fixture.GetCity()

		Convey("and an anaconda.Tweet with a poi", func() {
			tweet := Fixture.GetPoiTweet()

			Convey("it should create a geoevent", func() {
				geoevent, err := cityrecorder.NewGeoEventFromTweet(city, tweet)

				So(err, ShouldBeNil)
				So(geoevent.LocationType, ShouldEqual, "poi")
				So(geoevent.GeoJson.Type, ShouldEqual, "Polygon")
			})
		})

		Convey("and an anaconda.Tweet with a coordinate", func() {
			tweet := Fixture.GetCoordinateTweet()

			Convey("it should create a geoevent", func() {
				geoevent, err := cityrecorder.NewGeoEventFromTweet(city, tweet)

				So(err, ShouldBeNil)
				So(geoevent.LocationType, ShouldEqual, "coordinate")
				So(geoevent.GeoJson.Type, ShouldEqual, "point")
			})
		})
	})
}
