package citylib

import (
	ig "github.com/dimroc/go-instagram/instagram"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"regexp"
)

var (
	shortcodeRegex = regexp.MustCompile(`http.:\/\/.*instagram.com\/p\/([^\/]+)[\/]?`)
)

type MediaRetriever interface {
	GetShortcode(shortcode string) (*ig.Media, error)
}

type instagramLinkEnricher struct {
	mediaRetriever MediaRetriever
}

func NewInstagramLinkEnricher(clientId, clientSecret, accessToken string) Enricher {
	if len(clientId) == 0 || len(clientSecret) == 0 {
		Logger.Panic("Instagram Link Enricher needs proper IG credentials")
	}

	client := ig.NewClient(nil)
	client.ClientID = clientId
	client.ClientSecret = clientSecret
	client.AccessToken = accessToken

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
	matches := shortcodeRegex.FindStringSubmatch(g.ExpandedUrl)

	if len(matches) == 2 {
		shortcode := matches[1]
		Logger.Debug("Retrieving IG media with shortcode: %s", shortcode)
		media, err := enricher.mediaRetriever.GetShortcode(shortcode)
		if err != nil {
			Logger.Warning("Could not enrich geoevent with ig media from %s: %s", g.ExpandedUrl, err.Error())
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
