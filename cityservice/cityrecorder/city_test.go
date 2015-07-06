package cityrecorder_test

import (
	"github.com/dimroc/urbanevents/cityservice/cityrecorder"
	//. "github.com/dimroc/urbanevents/cityservice/utils"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"time"
)

func TestIntegrationGetDetails(t *testing.T) {
	Convey("Given a populated elasticsearch", t, func() {
		elastic := cityrecorder.NewElasticConnection(os.Getenv("ELASTICSEARCH_URL"))
		defer elastic.Connection.Close()
		for index, geoevent := range Fixture.GeoEvents {
			geoevent.CreatedAt = time.Now().AddDate(0, 0, -index)
			elastic.Write(geoevent)
		}

		elastic.Refresh()
		city := Fixture.GetCity()

		Convey("the stats should be correct", func() {
			detailed := city.GetDetails(elastic)
			So(len(detailed.Stats.TweetCounts), ShouldEqual, 7)
			So(len(detailed.Stats.InstagramCounts), ShouldEqual, 7)
			So(len(detailed.Stats.Days), ShouldEqual, 7)
			So(detailed.Stats.TweetCounts[0], ShouldEqual, 1)
			So(detailed.Stats.InstagramCounts[0], ShouldEqual, 0)
			So(detailed.Stats.Days[0].Day, ShouldEqual, time.Now().Day)
		})
	})
}
