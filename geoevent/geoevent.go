package geoevent

import (
	"github.com/azr/anaconda"
)

type geoEvent struct {
	GeoJson geoJson `json:"geojson"`
	Type    string  `json:"type"`
	Payload string  `json:"payload"`
}

type geoJson interface {}

type point struct {
	Coordinates [2]float64 `json:"coordinates"` // Coordinate always has to have exactly 2 values
	Type        string     `json:"type"`
}

type boundingBox struct {
	Coordinates [][][]float64 `json:"coordinates"`
	Type        string        `json:"type"`
}

func geoJsonFromPoint(t anaconda.Tweet) geoJson {
	return &point{
		Coordinates: t.Coordinates.Coordinates,
		Type:        t.Coordinates.Type,
	}
}

func geoJsonFromBoundingBox(t anaconda.Tweet) geoJson {
	return &boundingBox{
		Coordinates: t.Place.BoundingBox.Coordinates,
		Type:        t.Place.BoundingBox.Type,
	}
}

func NewFromTweet(t anaconda.Tweet) *geoEvent {
	if t.Coordinates != nil {
		return &geoEvent{
			GeoJson: geoJsonFromPoint(t),
			Type:    "tweet", Payload: t.Text,
		}
	} else {
		return &geoEvent{
			GeoJson: geoJsonFromBoundingBox(t),
			Type:    "tweet", Payload: t.Text,
		}
	}
}
