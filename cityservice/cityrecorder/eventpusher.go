package cityrecorder

import (
	eventsource "gopkg.in/antage/eventsource.v1"
	"net/http"
)

type EventPusher struct {
	EventSource eventsource.EventSource
}

func NewEventPusher() *EventPusher {
	es := eventsource.New(
		eventsource.DefaultSettings(),
		func(req *http.Request) [][]byte {
			return [][]byte{
				[]byte("X-Accel-Buffering: no"),
				[]byte("Access-Control-Allow-Origin: *"),
			}
		},
	)

	return &EventPusher{
		EventSource: es,
	}
}

func (ep *EventPusher) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ep.EventSource.ServeHTTP(rw, req)
}

func (ep *EventPusher) Write(g GeoEvent) error {
	ep.EventSource.SendEventMessage(g.ToJsonString(), "event", g.Id)
	return nil
}

func (ep *EventPusher) Close() {
	ep.EventSource.Close()
}
