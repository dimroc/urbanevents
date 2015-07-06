package cityrecorder

import (
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"strings"
)

type HoodEnricher struct {
	Elastic Elastic
}

func NewHoodEnricher(e Elastic) *HoodEnricher {
	return &HoodEnricher{
		Elastic: e,
	}
}

func (h *HoodEnricher) Enrich(g GeoEvent) GeoEvent {
	hoodKeys := h.Elastic.Percolate(g.GeoJson)
	Logger.Debug("Enriching geoevent with neighborhoods %s", hoodKeys)
	hoods := []string{}
	for _, key := range hoodKeys {
		hoodStart := strings.Index(key, ",")
		if hoodStart == -1 {
			Logger.Warning("Unable to retrieve hood from %s", key)
		} else {
			hoods = append(hoods, key[hoodStart+1:len(key)])
		}
	}

	g.Neighborhoods = hoods
	return g
}
