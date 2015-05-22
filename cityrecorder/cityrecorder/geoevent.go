package cityrecorder

import ()

type GeoEvent struct {
	GeoJson      GeoJson `json:"geojson"`
	Id           int64   `json:"id"`
	CityKey      string  `json:"city"`
	LocationType string  `json:"locationType"`
	Type         string  `json:"type"`
	Payload      string  `json:"payload"`
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
