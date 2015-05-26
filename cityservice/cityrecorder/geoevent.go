package cityrecorder

import (
	"time"
)

type GeoEvent struct {
	CreatedAt    time.Time `json:"createdAt"`
	GeoJson      GeoJson   `json:"geojson"`
	Id           string    `json:"id"`
	CityKey      string    `json:"city"`
	LocationType string    `json:"locationType"`
	Type         string    `json:"type"`
	Payload      string    `json:"payload"`
	Metadata     Metadata  `json:"metadata"`
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

type Metadata interface{}

type Tweet struct {
	ScreenName string   `json:"screenName"`
	Hashtags   []string `json:"hashtags"`
	MediaTypes []string `json:"mediaTypes"`
	MediaUrls  []string `json:"mediaUrls"`
}
