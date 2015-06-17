package cityrecorder_test

import (
	ig "github.com/carbocation/go-instagram/instagram"
	"github.com/dimroc/urbanevents/cityservice/cityrecorder"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"time"
)

func TestCreateGeoEventFromInstagram(t *testing.T) {
	Convey("Given an Instagram Media Object", t, func() {
		media := Fixture.GetInstagramMedia()[0]

		Convey("the geoevent should be correct", func() {
			geoevent := cityrecorder.CreateGeoEventFromInstagram(media)

			So(geoevent.Id, ShouldEqual, media.ID)
		})
	})
}
