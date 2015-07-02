package cityrecorder

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
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
	Counts []int       `json:"counts"`
	Days   []time.Time `json:"days"`
}

func (c *City) GetDetails(e Elastic) CityDetails {
	return c.GetDetailsFor(e, 7)
}

func (c *City) GetDetailsFor(e Elastic, nDays int) CityDetails {
	counts, days := c.retrieveStats(e, nDays)

	stats := CityCounts{
		Counts: counts,
		Days:   days,
	}

	return CityDetails{*c, stats}
}

func (c *City) retrieveStats(e Elastic, daysBack int) ([]int, []time.Time) {
	queryJson := `
{
  "size": 0,
  "aggs": {
    "tweet_count": {
      "filter": {
        "terms": {
          "city": [
            "%s"
          ]
        }
      },
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
`
	query := fmt.Sprintf(queryJson, c.Key, getDateRanges(daysBack))

	out := e.Search(query)
	response := aggregationResult{}

	err := json.Unmarshal(out.Aggregations, &response)
	if err != nil {
		log.Panic(err)
	}

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
	TweetCount struct {
		DocCount int64 `json:"doc_count"`
		Range    struct {
			Buckets []struct {
				Key      string    `json:"key"`
				DocCount int       `json:"doc_count"`
				To       time.Time `json:"to_as_string"`
			} `json:"buckets"`
		} `json:"range"`
	} `json:"tweet_count"`
}

func (a *aggregationResult) GetCountsAndDays() ([]int, []time.Time) {
	buckets := a.TweetCount.Range.Buckets
	length := len(buckets)
	counts := make([]int, length)
	days := make([]time.Time, length)

	for index, bucket := range buckets {
		// Reverse order of counts and days so it's descending
		counts[length-1-index] = bucket.DocCount
		days[length-1-index] = bucket.To
	}

	return counts, days
}
