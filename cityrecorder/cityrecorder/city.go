package cityrecorder

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type City struct {
	Name      string `json:"name"`
	Locations string `json:"locations"`
}

type Settings struct {
	Cities []City
}

func (s *Settings) save() error {
	filename := "conf.json"
	jsonOut, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return ioutil.WriteFile(filename, jsonOut, 0600)
}
