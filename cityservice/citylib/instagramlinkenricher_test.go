package citylib_test

import (
	"github.com/dimroc/urbanevents/cityservice/citylib"
	"github.com/dimroc/urbanevents/cityservice/mock_citylib"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInstagramLinkEnrich(t *testing.T) {
	Convey("Given a GeoEvent with an instagram url", t, func() {
		geoevent := Fixture.GeoEvents[0]
		geoevent.ExpandedUrl = "https://instagram.com/p/47N0xct3-P/"

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		media := &Fixture.GetInstagramMedia()[0]
		mediaRetriever := mock_citylib.NewMockMediaRetriever(mockCtrl)
		mediaRetriever.EXPECT().GetShortcode("47N0xct3-P").Return(media, nil)

		Convey("there should be returned neighborhoods", func() {
			enricher := citylib.NewInstagramLinkEnricherWithMediaRetriever(mediaRetriever)
			newGeo := enricher.Enrich(geoevent)

			So(geoevent.MediaType, ShouldEqual, "text")
			So(geoevent.MediaUrl, ShouldBeEmpty)

			So(newGeo.MediaType, ShouldEqual, "image")
			So(newGeo.MediaUrl, ShouldContainSubstring, "instagram")
		})
	})
}
