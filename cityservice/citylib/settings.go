package citylib

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Settings struct {
	Cities []City `json:"cities"`
	lookup map[string]City
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
