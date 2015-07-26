package citylib_test

import (
	"github.com/dimroc/urbanevents/cityservice/citylib"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCreateGeoEventFromInstagram(t *testing.T) {
	Convey("Given an Instagram Media Object", t, func() {
		media := Fixture.GetInstagramMedia()[0]

		Convey("the geoevent should be correct", func() {
			geoevent := citylib.CreateGeoEventFromInstagram(media)

			So(geoevent.Id, ShouldEqual, media.ID)
			So(geoevent.MediaUrl, ShouldEqual, media.Images.StandardResolution.URL)
			So(geoevent.Place, ShouldEqual, "Hill Country Chicken")
		})
	})

	Convey("Given an Instagram Media Object without a caption", t, func() {
		media := Fixture.GetInstagramMedia()[2]

		Convey("the geoevent should be correct", func() {
			geoevent := citylib.CreateGeoEventFromInstagram(media)

			So(geoevent.Id, ShouldEqual, media.ID)
			So(geoevent.MediaUrl, ShouldEqual, "https://scontent.cdninstagram.com/hphotos-xaf1/t51.2885-15/e15/11374545_116871491983909_567056437_n.jpg")
		})
	})
}
