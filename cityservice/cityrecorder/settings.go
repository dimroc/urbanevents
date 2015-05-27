package cityrecorder

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const (
	settingsFilename string = "conf.json"
)

type Settings struct {
	Cities []City `json:"cities"`
	lookup map[string]City
}

func (s *Settings) Save() error {
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

func LoadSettings() (Settings, error) {
	contents, err := ioutil.ReadFile(settingsFilename)
	if err != nil {
		log.Panic(err)
	}

	s := Settings{}
	s.lookup = make(map[string]City)
	err = json.Unmarshal(contents, &s)
	if err == nil {
		for _, city := range s.Cities {
			s.lookup[city.Key] = city
		}
	}

	return s, err
}

func (s *Settings) FindCity(cityKey string) City {
	return s.lookup[cityKey]
}
