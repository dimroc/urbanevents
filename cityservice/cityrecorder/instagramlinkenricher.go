package cityrecorder

import (
	ig "github.com/carbocation/go-instagram/instagram"
	. "github.com/dimroc/urbanevents/cityservice/utils"
)

type instagramLinkEnricher struct {
	client *ig.Client
}

func NewInstagramLinkEnricher(clientId, clientSecret string) Enricher {
	if len(clientId) == 0 || len(clientSecret) == 0 {
		Logger.Panic("Instagram Link Enricher needs proper IG credentials")
	}

	client := ig.NewClient(nil)
	client.ClientID = clientId
	client.ClientSecret = clientSecret

	return &instagramLinkEnricher{
		client: client,
	}
}

func (enricher *instagramLinkEnricher) Enrich(g GeoEvent) GeoEvent {
	Logger.Debug("Enriching geoevent with instagram media from text: %s", g.Payload)
	return g
}
