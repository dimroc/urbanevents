package citylib

import (
	"encoding/json"
	"fmt"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	elastigo "github.com/mattbaird/elastigo/lib"
	"io/ioutil"
	"log"
	"strings"
)

type Settings struct {
	Cities []City `json:"cities"`
	lookup map[string]City
}

func (s *Settings) GetCities(e Elastic) []City {
	return s.Cities
}

func (s *Settings) QueryCities(e Elastic, term string) []CityGeoEvents {
	if len(term) == 0 {
		return []CityGeoEvents{}
	}

	queryString := generateAcrossCityQuery(term, s.Cities)
	out := e.Search(queryString)

	result := cityAggregationResult{}
	err := json.Unmarshal(out.Aggregations, &result)
	Check(err)

	cityEvents := make([]CityGeoEvents, len(result.Cities.Buckets))
	for index, bucket := range result.Cities.Buckets {
		events := GeoEventsFromHits(&bucket.TopCityHits.Hits)
		cityBucket := CityGeoEvents{
			Key:       bucket.Key,
			GeoEvents: events,
		}

		cityEvents[index] = cityBucket
	}

	return cityEvents
}

func (s *Settings) GetCityDetails(e Elastic) []CityDetails {
	cities := make([]CityDetails, len(s.Cities))
	for index, city := range s.Cities {
		cities[index] = city.GetDetails(e)
	}

	return cities
}

func (s *Settings) Save(settingsFilename string) error {
	jsonOut, err := json.Marshal(s)
	if err != nil {
		log.Panic(err)
	}
	return ioutil.WriteFile(settingsFilename, jsonOut, 0644)
}

func (s *Settings) String() string {
	jsonOut, err := json.Marshal(s)
	if err != nil {
		log.Panic(err)
	}

	return string(jsonOut)
}

func LoadSettings(settingsFilename string) (Settings, error) {
	contents, err := ioutil.ReadFile(settingsFilename)
	if err != nil {
		log.Panic(err)
	}

	s := Settings{}
	s.lookup = make(map[string]City)
	err = json.Unmarshal(contents, &s)
	if err == nil {

		// Generate circles for cities and reassign settings.Cities
		citiesWithCircles := make([]City, len(s.Cities))
		for index, city := range s.Cities {
			city.GenerateCircles()
			citiesWithCircles[index] = city
		}

		// Create lookup map
		s.Cities = citiesWithCircles
		for _, city := range s.Cities {
			s.lookup[city.Key] = city
		}
	}

	return s, err
}

func (s *Settings) FindCity(cityKey string) City {
	return s.lookup[cityKey]
}

func generateAcrossCityQuery(term string, cities []City) string {
	queryJson := `
{
  "size": 0,
  "query": {
    "filtered": {
      "query": {
        "simple_query_string": {
          "query": "%s",
          "fields": [
            "text",
            "fullName",
            "hashtags",
            "username",
            "place"
          ]
        }
      },
      "filter": {
        "and": [
          { "terms": { "city": ["%s"] } },
          { "terms": { "mediaType": ["image","video"] } },
          { "exists" : { "field" : "neighborhoods" } }
        ]
      }
    }
  },
  "aggs": {
    "cities": {
      "terms": {
        "field": "city"
      },
      "aggs": {
        "top_city_hits": {
          "top_hits": {
            "sort": [
              {
                "createdAt": {
                  "order": "desc"
                }
              }
            ]
          }
        }
      }
    }
  }
}`

	cityKeys := make([]string, len(cities))
	for index, city := range cities {
		cityKeys[index] = city.Key
	}

	return fmt.Sprintf(queryJson, term, strings.Join(cityKeys, `","`))
}

type cityAggregationResult struct {
	Cities struct {
		Buckets []struct {
			Key         string `json:"key"`
			DocCount    int    `json:"doc_count"`
			TopCityHits struct {
				Hits elastigo.Hits `json:"hits"`
			} `json:"top_city_hits"`
		} `json:"buckets"`
	} `json:"cities"`
}
