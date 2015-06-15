package main

import (
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/dimroc/urbanevents/cityservice/cityrecorder"
	"github.com/unrolled/render"
	"log"
	"net/http"
)

func SettingsHandler(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{IndentJSON: true})
	settings := context.Get(req, CTX_SETTINGS_KEY)
	r.JSON(w, http.StatusOK, settings)
}

func CitiesHandler(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{IndentJSON: true})
	settings := GetSettings(req)

	r.JSON(w, http.StatusOK, settings.GetCityDetails(GetElasticConnection(req)))
}

func CityHandler(w http.ResponseWriter, req *http.Request) {
	city := GetCity(req)

	r := render.New(render.Options{IndentJSON: true})
	r.JSON(w, http.StatusOK, city.GetDetails(GetElasticConnection(req)))
}

func GetCity(req *http.Request) cityrecorder.City {
	vars := mux.Vars(req)
	cityKey := vars["city"]
	settings := GetSettings(req)
	return settings.FindCity(cityKey)
}

func GetSettings(req *http.Request) cityrecorder.Settings {
	if rv := context.Get(req, CTX_SETTINGS_KEY); rv != nil {
		return rv.(cityrecorder.Settings)
	}

	log.Panic("Could not retrieve Settings")
	return cityrecorder.Settings{}
}

func GetElasticConnection(req *http.Request) cityrecorder.Elastic {
	if rv := context.Get(req, CTX_ELASTIC_CONNECTION_KEY); rv != nil {
		return rv.(cityrecorder.Elastic)
	}

	log.Panic("Could not retrieve Elastic Connection")
	return nil
}
