package cityrecorder

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type City struct {
	Key     string      `json:"key"`
	Display string      `json:"display"`
	Bounds  [][]float64 `json:"bounds"` //long,lat pair defining the bounding rectangle
}

type CityDetailed struct {
	City
	Stats CityCounts
}

func (c *City) LocationString() string {
	return fmt.Sprintf("%v,%v,%v,%v",
		c.Bounds[0][0],
		c.Bounds[0][1],
		c.Bounds[1][0],
		c.Bounds[1][1],
	)
}

type CityCounts struct {
	Counts []int       `json:"counts"`
	Days   []time.Time `json:"days"`
}

func (c *City) GetStats(e *ElasticConnection) CityDetailed {
	return c.GetStatsFor(e, 7)
}

func (c *City) GetStatsFor(e *ElasticConnection, nDays int) CityDetailed {
	counts := make([]int, nDays)
	days := make([]time.Time, nDays)
	current := time.Now()

	for index := 0; index < nDays; index++ {
		counts[index] = nDays - index
		prevDay := current.AddDate(0, 0, -1)
		statsForDay(e, index)
		days[index] = current
		current = prevDay
	}

	stats := CityCounts{
		Counts: counts,
		Days:   days,
	}

	return CityDetailed{*c, stats}
}

func statsForDay(e *ElasticConnection, daysBack int) {
	queryJson := `{
		"size": 0,
		"aggs": {
			"tweet_count": {
				"filter": {
					"range": {
						"createdAt": {
							"lt": "%s",
							"gte": "%s"
						}
					}
				},
				"aggs": {
					"city": {
						"terms": {
							"field": "city"
						}
					}
				}
			}
		}
	}`

	var lt string
	if daysBack == 0 {
		lt = "now"
	} else {
		lt = time.Now().UTC().AddDate(0, 0, -daysBack+1).String()[0:10]
	}

	prevDate := time.Now().UTC().AddDate(0, 0, -daysBack).String()[0:10]
	gte := prevDate
	query := fmt.Sprintf(queryJson, lt, gte)

	out := e.Search(query)
	response := aggregationResult{}

	err := json.Unmarshal(out.Aggregations, &response)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(response)
}

type aggregationResult struct {
	TweetCount struct {
		DocCount int64 `json:"doc_count"`
		City     struct {
			Buckets []struct {
				Key      string `json:"key"`
				DocCount int64  `json:"doc_count"`
			} `json:"buckets"`
		} `json:"city"`
	} `json:"tweet_count"`
}
