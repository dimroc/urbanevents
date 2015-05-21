package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
