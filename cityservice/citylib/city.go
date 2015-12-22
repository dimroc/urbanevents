package citylib

import (
	"encoding/json"
	"fmt"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	elastigo "github.com/mattbaird/elastigo/lib"
	"strings"
	"time"
)

const (
	CITY_QUERY_SIZE = "25"
)

type City struct {
	Key         string      `json:"key"`
	Display     string      `json:"display"`
	Aliases     []string    `json:"aliases"`
	BoundingBox [4]float64  `json:"bbox"` //long,lat pair defining the bounding rectangle
	Circles     []GeoCircle `json:"circles"`
}

type CityDetails struct {
	City
	Stats CityCounts `json:"stats"`
}

func (c *City) String() string {
	return fmt.Sprintf("%s %s %s", c.Key, c.Display, c.LocationString())
}

func (c *City) LocationString() string {
	return fmt.Sprintf("%f,%f,%f,%f",
		c.BoundingBox[0],
		c.BoundingBox[1],
		c.BoundingBox[2],
		c.BoundingBox[3],
	)
}

func (c *City) GenerateCircles() {
	c.Circles = PackCircles(c.BoundingBox)
}

type CityCounts struct {
	TweetCounts     []int       `json:"tweetCounts"`
	InstagramCounts []int       `json:"instagramCounts"`
	Days            []time.Time `json:"days"`
}

func (c *City) Query(e Elastic, term string) []GeoEvent {
	if len(term) == 0 {
		return []GeoEvent{}
	}

	dsl := elastigo.Search(ES_IndexName).Type(ES_TypeName).Size(CITY_QUERY_SIZE).Pretty().Filter(
		elastigo.Filter().And(
			elastigo.Filter().Term("city", c.Key),
			elastigo.Filter().Terms("mediaType", elastigo.TEMPlain, "image", "video"),
		),
	).Query(
		elastigo.Query().Search(term),
	).Sort(
		elastigo.Sort("createdAt").Desc(),
	)

	out := e.SearchDsl(*dsl)

	return GeoEventsFromElasticSearch(out)
}

func (c *City) GetDetails(e Elastic) CityDetails {
	return c.GetDetailsFor(e, 7)
}

func (c *City) GetDetailsFor(e Elastic, nDays int) CityDetails {
	tweets, instagrams, days := c.retrieveStats(e, nDays)

	stats := CityCounts{
		TweetCounts:     tweets,
		InstagramCounts: instagrams,
		Days:            days,
	}

	return CityDetails{*c, stats}
}

func (c *City) retrieveStats(e Elastic, daysBack int) ([]int, []int, []time.Time) {
	queryJson := `
{
  "size": 0,
  "aggs": {
    "counts": {
      "filter": { "term": { "city": "%s" } },
      "aggs":
        { "twitter" : {
          "filter": { "term": { "service": "twitter" }},
          "aggs": {
            "range": {
              "date_range": {
                "field": "createdAt",
                "ranges": [%s]
              }
            }
          }
        },
         "instagram" : {
          "filter": { "term": { "service": "instagram" }},
          "aggs": {
            "range": {
              "date_range": {
                "field": "createdAt",
                "ranges": [%s]
              }
            }
          }
        }
			}
		}
	}
}
`
	dateRanges := getDateRanges(daysBack)
	query := fmt.Sprintf(queryJson, c.Key, dateRanges, dateRanges)

	out := e.Search(query)
	response := aggregationResult{}

	err := json.Unmarshal(out.Aggregations, &response)
	Check(err)
	return response.GetCountsAndDays()
}

func getDateRanges(nDays int) string {
	entries := make([]string, nDays)

	for index := 0; index < nDays; index++ {
		entries[index] = getDateRangeFor(index)
	}

	return strings.Join(entries, ",\n")
}

// Help w date range aggregations
//
//{
//"key": "testkey",
//"to": "now",
//"from": "now/d"
//}
func getDateRangeFor(daysBack int) string {
	var key, lt, gte string
	if daysBack == 0 {
		lt = "now"
		key = "now"
		gte = "now/d"
	} else {
		lt = time.Now().UTC().AddDate(0, 0, -daysBack+1).String()[0:10] + "||/d"
		key = lt[0:10]
		gte = time.Now().UTC().AddDate(0, 0, -daysBack).String()[0:10] + "||/d"
	}

	return fmt.Sprintf(`{"key": "%s", "to": "%s", "from": "%s"}`, key, lt, gte)
}

type aggregationResult struct {
	Counts struct {
		Twitter struct {
			DocCount int64 `json:"doc_count"`
			Range    struct {
				Buckets []struct {
					Key      string    `json:"key"`
					DocCount int       `json:"doc_count"`
					To       time.Time `json:"to_as_string"`
				} `json:"buckets"`
			} `json:"range"`
		} `json:"twitter"`
		Instagram struct {
			DocCount int64 `json:"doc_count"`
			Range    struct {
				Buckets []struct {
					Key      string    `json:"key"`
					DocCount int       `json:"doc_count"`
					To       time.Time `json:"to_as_string"`
				} `json:"buckets"`
			} `json:"range"`
		} `json:"instagram"`
	} `json:"counts"`
}

func (a *aggregationResult) GetCountsAndDays() ([]int, []int, []time.Time) {
	buckets := a.Counts.Twitter.Range.Buckets

	length := len(buckets)
	days := make([]time.Time, length)
	tweets := make([]int, length)

	for index, bucket := range buckets {
		// Reverse order of tweets and days so it's descending
		tweets[length-1-index] = bucket.DocCount
		days[length-1-index] = bucket.To
	}

	buckets = a.Counts.Instagram.Range.Buckets
	instagrams := make([]int, length)

	for index, bucket := range buckets {
		// Reverse order of instagram and days so it's descending
		instagrams[length-1-index] = bucket.DocCount
	}

	return tweets, instagrams, days
}
