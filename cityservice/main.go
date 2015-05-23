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
	CtxSettingsKey = "city.settings"
)

func SettingsMiddleware(settings cityrecorder.Settings) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		context.Set(r, CtxSettingsKey, settings)
		next(w, r)
	})
}

func main() {
	settings, settingsErr := cityrecorder.LoadSettings()
	if settingsErr != nil {
		log.Fatal(settingsErr)
	}

	recorder := cityrecorder.NewTweetRecorder(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
		os.Getenv("TWITTER_TOKEN"),
		os.Getenv("TWITTER_TOKEN_SECRET"),
	)

	pusher := cityrecorder.NewPusherFromURL(os.Getenv("PUSHER_URL"))

	for _, city := range settings.Cities {
		fmt.Println("Configuring city:", city)
		go recorder.Record(city, pusher)
	}

	router := mux.NewRouter()
	apiRoutes := router.PathPrefix("/api/v1").Subrouter()
	apiRoutes.HandleFunc("/settings", SettingsHandler).Methods("GET")

	n := negroni.Classic()
	n.Use(cors.Default())
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.Use(SettingsMiddleware(settings))
	n.UseHandler(context.ClearHandler(router))
	n.Run(":8080")
}

func SettingsHandler(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{IndentJSON: true})
	settings := context.Get(req, CtxSettingsKey)
	r.JSON(w, http.StatusOK, settings)
}
