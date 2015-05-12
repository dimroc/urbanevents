package geoevent

import (
	"github.com/azr/anaconda"
)

type GeoEvent struct {
	GeoJson GeoJson `json:"geojson"`
	Type    string  `json:"type"`
	Payload string  `json:"payload"`
}

type GeoJson interface{}

type Point struct {
	Coordinates [2]float64 `json:"coordinates"` // Coordinate always has to have exactly 2 values
	Type        string     `json:"type"`
}

type BoundingBox struct {
	Coordinates [][][]float64 `json:"coordinates"`
	Type        string        `json:"type"`
}

func geoJsonFromPoint(t anaconda.Tweet) GeoJson {
	return &Point{
		Coordinates: t.Coordinates.Coordinates,
		Type:        t.Coordinates.Type,
	}
}

func geoJsonFromBoundingBox(t anaconda.Tweet) GeoJson {
	return &BoundingBox{
		Coordinates: t.Place.BoundingBox.Coordinates,
		Type:        t.Place.BoundingBox.Type,
	}
}

func NewFromTweet(t anaconda.Tweet) *GeoEvent {
	if t.Coordinates != nil {
		return &GeoEvent{
			GeoJson: geoJsonFromPoint(t),
			Type:    "tweet",
			Payload: t.Text,
		}
	} else {
		return &GeoEvent{
			GeoJson: geoJsonFromBoundingBox(t),
			Type:    "tweet",
			Payload: t.Text,
		}
	}
}