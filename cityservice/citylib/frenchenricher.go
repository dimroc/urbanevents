package citylib

import (
	. "github.com/dimroc/urbanevents/cityservice/utils"
)

type FrenchEnricher struct {
}

func NewFrenchEnricher() *FrenchEnricher {
	return &FrenchEnricher{}
}

func (h *FrenchEnricher) Enrich(g GeoEvent) GeoEvent {
	if g.CityKey == "paris" {
		Logger.Debug("Enriching geoevent with french because of city %s", g.CityKey)
		g.TextFrench = g.Text
	}

	return g
}
