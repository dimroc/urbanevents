package main

import (
	"github.com/codegangsta/negroni"
	"github.com/dimroc/urbanevents/cityservice/citylib"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/unrolled/render"
	"net/http"
	"os"
)

const (
	CTX_ELASTIC_CONNECTION_KEY = "city.elasticconnection"
)

var (
	settingsFilename = GetenvOrDefault("CITYSERVICE_SETTINGS", "../config/nyc.json")
)

func main() {
	Logger.Info("Running in " + GO_ENV)
	settings, settingsErr := citylib.LoadSettings(settingsFilename)
	Check(settingsErr)

	elastic := citylib.NewElasticConnection(os.Getenv("ELASTICSEARCH_URL"))

	router := mux.NewRouter()
	apiRoutes := router.PathPrefix("/api/v1").Subrouter()
	apiRoutes.HandleFunc("/cities", CitiesHandler).Methods("GET")
	apiRoutes.HandleFunc("/cities/{city}", citylib.CityHandler).Methods("GET")

	n := negroni.Classic()
	n.Use(cors.Default())
	n.Use(citylib.SettingsMiddleware(settings))
	n.Use(citylib.ElasticMiddleware(elastic))
	n.UseHandler(context.ClearHandler(router))
	n.Run(":5000")
}

func CitiesHandler(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{IndentJSON: true})
	settings := citylib.GetSettings(req)
	r.JSON(w, http.StatusOK, settings.Cities)
}
