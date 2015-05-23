package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/dimroc/urban-events/cityservice/cityrecorder"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"os"
)

func SettingsMiddleware(settings cityrecorder.Settings) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		context.Set(r, "settings", settings)
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

	mux := mux.NewRouter()
	mux.HandleFunc("/api/v1/settings", Settings).Methods("GET")

	n := negroni.Classic()
	n.Use(SettingsMiddleware(settings))
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(context.ClearHandler(mux))
	n.Run(":8080")
}

func Settings(w http.ResponseWriter, req *http.Request) {
	r := render.New(render.Options{IndentJSON: true})
	settings := context.Get(req, "settings")
	r.JSON(w, http.StatusOK, settings)
}
