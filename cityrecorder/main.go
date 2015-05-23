package main

import (
	"fmt"
	"github.com/dimroc/urban-events/cityrecorder/cityrecorder"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"
	"log"
	"os"
)

var (
	settings, settingsErr = cityrecorder.LoadSettings()
)

func main() {
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

	m := martini.Classic()
	m.Use(gzip.All())
	m.Use(render.Renderer())
	m.Get("/api/v1/settings", Settings)
	m.RunOnAddr(":8080")
}

func Settings(r render.Render) {
	r.JSON(200, settings)
}
