package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/dimroc/urban-events/cityservice/cityrecorder"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/rs/cors"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"os"
)

const (
	CtxSettingsKey          = "city.settings"
	CtxElasticConnectionKey = "city.elasticconnection"
)

func ElasticConnectionMiddleware(e *cityrecorder.ElasticConnection) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		context.Set(r, CtxElasticConnectionKey, e)
		next(w, r)
	})
}

func SettingsMiddleware(settings cityrecorder.Settings) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		context.Set(r, CtxSettingsKey, settings)
		next(w, r)
	})
}

func main() {
	settings, settingsErr := cityrecorder.LoadSettings()
	if settingsErr != nil {
		log.Panic(settingsErr)
	}

	recorder := cityrecorder.NewTweetRecorder(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_TOKEN"),
		os.Getenv("TWITTER_TOKEN_SECRET"),
	)

	pusher := cityrecorder.NewPusherFromURL(os.Getenv("PUSHER_URL"))
	elastic := cityrecorder.NewElasticConnection(os.Getenv("ELASTICSEARCH_URL"))
	broadcaster := cityrecorder.NewBroadcastWriter(pusher, elastic, cityrecorder.StdoutWriter)

	for _, city := range settings.Cities {
		fmt.Println("Configuring city:", city)
		go recorder.Record(city, broadcaster)
	}

	router := mux.NewRouter()
	apiRoutes := router.PathPrefix("/api/v1").Subrouter()
	apiRoutes.HandleFunc("/settings", SettingsHandler).Methods("GET")
	apiRoutes.HandleFunc("/cities/{city}", CityHandler).Methods("GET")

	n := negroni.Classic()
	n.Use(cors.Default())
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.Use(SettingsMiddleware(settings))
	n.Use(ElasticConnectionMiddleware(elastic))
	n.UseHandler(context.ClearHandler(router))
	n.Run(":8080")
}

func SettingsHandler(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{IndentJSON: true})
	settings := context.Get(req, CtxSettingsKey)
	r.JSON(w, http.StatusOK, settings)
}

func CityHandler(w http.ResponseWriter, req *http.Request) {
	city := GetCity(req)

	r := render.New(render.Options{IndentJSON: true})
	r.JSON(w, http.StatusOK, city.GetStats(GetElasticConnection(req)))
}

func GetCity(req *http.Request) cityrecorder.City {
	vars := mux.Vars(req)
	cityKey := vars["city"]
	settings := GetSettings(req)
	return settings.FindCity(cityKey)
}

func GetSettings(req *http.Request) cityrecorder.Settings {
	if rv := context.Get(req, CtxSettingsKey); rv != nil {
		return rv.(cityrecorder.Settings)
	}

	log.Panic("Could not retrieve Settings")
	return cityrecorder.Settings{}
}

func GetElasticConnection(req *http.Request) *cityrecorder.ElasticConnection {
	if rv := context.Get(req, CtxElasticConnectionKey); rv != nil {
		return rv.(*cityrecorder.ElasticConnection)
	}

	log.Panic("Could not retrieve Elastic Connection")
	return nil
}
