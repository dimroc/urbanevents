package main

import (
	"github.com/dimroc/urbanevents/cityservice/cityrecorder"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"os"
)

func main() {
	recorder := cityrecorder.NewInstagramRecorder(
		os.Getenv("INSTAGRAM_CLIENT_ID"),
		os.Getenv("INSTAGRAM_CLIENT_SECRET"),
		cityrecorder.StdoutWriter,
		nil,
	)

	defer recorder.Close()
	recorder.DeleteAllSubscriptions()
	printExistingSubscriptions(recorder)
}

func printExistingSubscriptions(recorder *cityrecorder.InstagramRecorder) {
	Logger.Debug("Existing Subscriptions:")
	subscriptions := recorder.GetSubscriptions()
	for _, subscription := range subscriptions {
		Logger.Debug(ToJsonStringUnsafe(subscription))
	}
}
