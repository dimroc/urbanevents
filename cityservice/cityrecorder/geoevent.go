package cityrecorder

import (
	"fmt"
	"time"
)

type GeoEvent struct {
	CreatedAt    time.Time  `json:"createdAt"`
	GeoJson      GeoJson    `json:"geojson"`
	Point        [2]float64 `json:"point"`
	Id           string     `json:"id"`
	CityKey      string     `json:"city"`
	LocationType string     `json:"locationType"`
	Type         string     `json:"type"`
	Payload      string     `json:"payload"`
	Metadata     Metadata   `json:"metadata"`
}

type GeoJson interface {
	Center() [2]float64
}

type Point struct {
	Coordinates [2]float64 `json:"coordinates"` // Coordinate always has to have exactly 2 values
	Type        string     `json:"type"`
}

func (p *Point) Center() [2]float64 {
	return p.Coordinates
}

type BoundingBox struct {
	Coordinates [][][]float64 `json:"coordinates"`
	Type        string        `json:"type"`
}

func (bb *BoundingBox) Center() [2]float64 {
	center := [2]float64{bb.Coordinates[0][0][0], bb.Coordinates[0][0][1]}
	return center
}

type Metadata interface{}

type Tweet struct {
	ScreenName string   `json:"screenName"`
	Hashtags   []string `json:"hashtags"`
	MediaTypes []string `json:"mediaTypes"`
	MediaUrls  []string `json:"mediaUrls"`
}

func (g *GeoEvent) String() string {
	return fmt.Sprintf(
		"{CreatedAt: %s, GeoJson: %s, Id: %s, CityKey: %s, LocationType: %s, Type: %s, Payload: %s, Metadata: %s}",
		g.CreatedAt,
		g.GeoJson,
		g.Id,
		g.CityKey,
		g.LocationType,
		g.Type,
		g.Payload,
		g.Metadata,
	)
}
