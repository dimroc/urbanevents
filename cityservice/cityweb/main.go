package citysearch

import (
	"github.com/codegangsta/negroni"
	"github.com/dimroc/urbanevents/cityservice/cityrecorder"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/unrolled/render"
)

const (
	CTX_ELASTIC_CONNECTION_KEY = "city.elasticconnection"
)

func main() {
	Logger.Info("Running in " + GO_ENV)

	elastic := cityrecorder.NewElasticConnection(os.Getenv("ELASTICSEARCH_URL"))

	router := mux.NewRouter()
	apiRoutes := router.PathPrefix("/api/v1").Subrouter()
	apiRoutes.HandleFunc("/cities", CitiesHandler).Methods("GET")
	apiRoutes.HandleFunc("/cities/{city}", CityHandler).Methods("GET")

	n := negroni.Classic()
	n.Use(cors.Default())
	n.Use(cityrecorder.SettingsMiddleware(settings))
	n.Use(cityrecorder.ElasticMiddleware(elastic))
	n.UseHandler(context.ClearHandler(router))
	n.Run()
}

func CitiesHandler(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{IndentJSON: true})
	settings := GetSettings(req)
	r.JSON(w, http.StatusOK, settings.Cities)
}
