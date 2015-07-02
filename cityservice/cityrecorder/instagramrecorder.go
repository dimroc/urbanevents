package cityrecorder

import (
	"encoding/json"
	"fmt"
	ig "github.com/carbocation/go-instagram/instagram"
	"github.com/dimroc/urbanevents/cityservice/set"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

const (
	secondsUntilMediaRetrieval = 5
)

type InstagramRecorder struct {
	clientId        string
	clientSecret    string
	writer          Writer
	client          *ig.Client
	ticker          *time.Ticker
	geographyIds    *set.SetValue
	geographyMinIds map[string]string
}

type geographEntry struct {
	City string
}

func NewInstagramRecorder(clientId, clientSecret string, writer Writer) *InstagramRecorder {
	client := ig.NewClient(nil)
	client.ClientID = clientId
	client.ClientSecret = clientSecret

	recorder := &InstagramRecorder{
		clientId:        clientId,
		clientSecret:    clientSecret,
		writer:          writer,
		client:          client,
		ticker:          time.NewTicker(time.Second * secondsUntilMediaRetrieval),
		geographyIds:    set.NewSetValue(),
		geographyMinIds: make(map[string]string),
	}

	if !recorder.Configured() {
		log.Fatal(fmt.Sprintf("Recorder configuration is invalid: %s", recorder))
	}

	go recorder.startMediaFetcher()
	return recorder
}

func (p *InstagramRecorder) Configured() bool {
	if len(p.clientId) == 0 || len(p.clientSecret) == 0 {
		return false
	}

	return true
}

func (recorder *InstagramRecorder) GetSubscriptions() []ig.Realtime {
	subscriptions, err := recorder.client.Realtime.ListSubscriptions()
	Check(err)
	return subscriptions
}

// Initialize Real-Time Subscriptions with Instagram, if necessary.
func (recorder *InstagramRecorder) Subscribe(baseUrl string, cities []City) {
	//lat, lng string, radius int, callbackURL, verifyToken string
	for _, city := range cities {
		for _, circle := range city.Circles {
			// Using the circle packer generated circles, register each circle for that city via instagram.
			response, err := recorder.client.Realtime.SubscribeToGeography(
				circle.LatString(),
				circle.LngString(),
				circle.Meters(),
				baseUrl+city.Key,
				"cityservice",
			)

			Logger.Debug(ToJsonStringUnsafe(response))
			Check(err)
		}
	}
}

func (recorder *InstagramRecorder) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" || len(req.Method) == 0 {
		ig.ServeInstagramRealtimeSubscribe(rw, req)
	} else if req.Method == "POST" {
		vars := mux.Vars(req)
		cityKey := vars["city"]

		decoder := json.NewDecoder(req.Body)
		var jsonResponses []ig.RealtimeResponse
		err := decoder.Decode(&jsonResponses)

		if err != nil {
			Logger.Warning("%s", err)
		} else {
			// Hand off all responses to another goroutine to fetch RecentMedia so we free up this POST call.
			Logger.Debug(cityKey + ": " + ToJsonStringUnsafe(jsonResponses))
			for _, jsonResponse := range jsonResponses {
				recorder.geographyIds.Add(jsonResponse.ObjectID, geographEntry{City: cityKey})
			}
		}
	}
}

func (recorder *InstagramRecorder) startMediaFetcher() {
	for _ = range recorder.ticker.C {
		entries := recorder.geographyIds.ListAndClear()
		for _, entry := range entries {
			// Entry: { Key string, Value interface{} ({ City: "paris" })}
			recorder.retrieveMediaFor(entry.Key, entry.Value.(geographEntry).City)
		}
	}
}

func (recorder *InstagramRecorder) retrieveMediaFor(geographyId, cityKey string) {
	parameters := ig.Parameters{
		MinID: recorder.geographyMinIds[geographyId],
	}

	medias, _, err := recorder.client.Geographies.RecentMedia(geographyId, &parameters)
	if err != nil {
		Logger.Warning("Unable to retrieve media", err)
		return
	}

	if len(medias) > 0 {
		recorder.geographyMinIds[geographyId] = medias[0].ID
	}

	Logger.Debug("%s", ToJsonStringUnsafe(medias))
	for _, media := range medias {
		Logger.Debug("CREATING GEOEVENT %s", ToJsonStringUnsafe(media))
		geoevent := CreateGeoEventFromInstagram(media)
		geoevent.CityKey = cityKey
		recorder.writer.Write(geoevent)
	}
}

func CreateGeoEventFromInstagram(media ig.Media) GeoEvent {
	var mediaUrl string
	if media.Type == "video" {
		mediaUrl = safelyRetrieveVideo(media)
	} else {
		mediaUrl = safelyRetrieveImage(media)
	}

	return GeoEvent{
		CreatedAt:    time.Unix(media.CreatedTime, 0),
		Id:           media.ID,
		FullName:     media.User.FullName, // New to GeoEvent
		Hashtags:     media.Tags,
		Link:         media.Link, // New to GeoEvent
		LocationType: "coordinate",
		MediaType:    media.Type, // New to GeoEvent // Either image or video
		MediaUrl:     mediaUrl,   // New to GeoEvent
		Payload:      safelyRetrieveCaption(media),
		Point:        [2]float64{media.Location.Longitude, media.Location.Latitude},
		Service:      "instagram",
		ThumbnailUrl: safelyRetrieveThumbnail(media), // New to GeoEvent
		Type:         "geoevent",
		Username:     media.User.Username, // New to GeoEvent
	}
}

func (recorder *InstagramRecorder) Close() {
	recorder.ticker.Stop()
}

func (recorder *InstagramRecorder) DeleteAllSubscriptions() {
	Logger.Warning("Deleting all Instagram Real-time Subscriptions!")
	recorder.client.Realtime.DeleteAllSubscriptions()
}

func safelyRetrieveVideo(media ig.Media) string {
	if media.Videos != nil {
		return media.Videos.StandardResolution.URL
	} else {
		return ""
	}
}

func safelyRetrieveThumbnail(media ig.Media) string {
	if media.Images != nil {
		return media.Images.Thumbnail.URL
	} else {
		return ""
	}
}

func safelyRetrieveImage(media ig.Media) string {
	if media.Images != nil {
		return media.Images.StandardResolution.URL
	} else {
		return ""
	}
}

func safelyRetrieveCaption(media ig.Media) string {
	if media.Caption != nil {
		return media.Caption.Text
	} else {
		return ""
	}
}
