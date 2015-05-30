package main

import (
	"github.com/dimroc/urbanevents/cityservice/cityrecorder"
	"os"
)

func main() {
	pusher := cityrecorder.NewPusherFromURL(os.Getenv("PUSHER_URL"))
	pusher.Listen(os.Stdin)
}
