package cityrecorder

import (
	ig "github.com/carbocation/go-instagram/instagram"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"regexp"
)

var (
	mediaIdRegex = regexp.MustCompile(`http.:\/\/.*instagram.com\/p\/(\w*)[\/]?`)
)

type MediaRetriever interface {
	Get(mediaID string) (*ig.Media, error)
}

type instagramLinkEnricher struct {
	mediaRetriever MediaRetriever
}

func NewInstagramLinkEnricher(clientId, clientSecret string) Enricher {
	if len(clientId) == 0 || len(clientSecret) == 0 {
		Logger.Panic("Instagram Link Enricher needs proper IG credentials")
	}

	client := ig.NewClient(nil)
	client.ClientID = clientId
	client.ClientSecret = clientSecret

	return NewInstagramLinkEnricherWithMediaRetriever(client.Media)
}

// Inject dependency to improve testability.
func NewInstagramLinkEnricherWithMediaRetriever(retriever MediaRetriever) Enricher {
	return &instagramLinkEnricher{
		mediaRetriever: retriever,
	}
}

func (enricher *instagramLinkEnricher) Enrich(g GeoEvent) GeoEvent {
	Logger.Debug("Enriching geoevent with instagram media from url: %s", g.ExpandedUrl)
	matches := mediaIdRegex.FindStringSubmatch(g.ExpandedUrl)

	if len(matches) == 2 {
		mediaId := matches[1]
		Logger.Debug("Retrieving IG media with id: %s", mediaId)
		media, err := enricher.mediaRetriever.Get(mediaId)
		if err != nil {
			Logger.Warning("Could not enrich geoevent with ig media from g.ExpandedUrl: %s", err)
		} else {
			newGeo := g
			newGeo.MediaType = media.Type
			newGeo.ThumbnailUrl = SafelyRetrieveInstagramThumbnail(*media)
			newGeo.MediaUrl = SafelyRetrieveInstagramMediaUrl(*media)
			return newGeo
		}
	}

	return g
}
