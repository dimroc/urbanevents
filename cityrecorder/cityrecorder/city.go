package cityrecorder

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const (
	filename string = "conf.json"
)

type City struct {
	Key       string `json:"key"`
	Display   string `json:"display"`
	Locations string `json:"locations"`
}

type Settings struct {
	Cities []City
}

func (s *Settings) Save() error {
	jsonOut, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return ioutil.WriteFile(filename, jsonOut, 0644)
}

func LoadSettings() (Settings, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	s := Settings{}
	err = json.Unmarshal(contents, &s)
	return s, err
}
