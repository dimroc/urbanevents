package cityrecorder_test

import (
	cr "github.com/dimroc/urbanevents/cityservice/cityrecorder"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPackCircles(t *testing.T) {
	Convey("Given a bounding box", t, func() {
		city := Fixture.GetCity()
		boundingBox := city.BoundingBox

		Convey("there should be many hexagonal circles", func() {
			circles := cr.PackCircles(boundingBox)
			So(len(circles), ShouldEqual, 50)
		})
	})
}
