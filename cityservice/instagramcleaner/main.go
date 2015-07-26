package main

import (
	"github.com/dimroc/urbanevents/cityservice/citylib"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"os"
)

func main() {
	recorder := citylib.NewInstagramRecorder(
		os.Getenv("INSTAGRAM_CLIENT_ID"),
		os.Getenv("INSTAGRAM_CLIENT_SECRET"),
		citylib.StdoutWriter,
		nil,
	)

	defer recorder.Close()
	recorder.DeleteAllSubscriptions()
	printExistingSubscriptions(recorder)
}

func printExistingSubscriptions(recorder *citylib.InstagramRecorder) {
	Logger.Debug("Existing Subscriptions:")
	subscriptions := recorder.GetSubscriptions()
	for _, subscription := range subscriptions {
		Logger.Debug(ToJsonStringUnsafe(subscription))
	}
}
