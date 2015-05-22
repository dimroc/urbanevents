package main

import (
	"fmt"
	"github.com/dimroc/urban-events/cityrecorder/cityrecorder"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	recorder := cityrecorder.TweetRecorder{
		ConsumerKey:    *consumerKey,
		ConsumerSecret: *consumerSecret,
		Token:          *token,
		TokenSecret:    *tokenSecret,
	}

	pusher := cityrecorder.Pusher{
		AppId:  *pusherAppId,
		Key:    *pusherKey,
		Secret: *pusherSecret,
	}

	pusher.Start(os.Stdin)

	settings, err := cityrecorder.LoadSettings()
	if err != nil {
		log.Fatal(err)
	}

	for _, city := range settings.Cities {
		recorder.Start(city, writer)
		pusher.Start(reader)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
