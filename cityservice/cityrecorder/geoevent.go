package cityrecorder

import (
	"fmt"
	"time"
)

type GeoEvent struct {
	CityKey      string     `json:"city"`
	CreatedAt    time.Time  `json:"createdAt"`
	FullName     string     `json:"fullName"`
	GeoJson      GeoJson    `json:"geojson"`
	Hashtags     []string   `json:"hashtags"`
	Id           string     `json:"id"`
	ImageUrl     string     `json:"imageUrl"`
	Link         string     `json:"link"`
	LocationType string     `json:"locationType"`
	MediaType    string     `json:"mediaType"`
	Payload      string     `json:"payload"`
	Point        [2]float64 `json:"point"`
	Service      string     `json:"service"`
	ThumbnailUrl string     `json:"thumbnailUrl"`
	Type         string     `json:"type"`
	Username     string     `json:"username"`
	VideoUrl     string     `json:"videoUrl"`
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
	// TODO: Get average
	center := [2]float64{bb.Coordinates[0][0][0], bb.Coordinates[0][0][1]}
	return center
}

func (g *GeoEvent) String() string {
	return fmt.Sprintf(
		"{CreatedAt: %s, GeoJson: %s, Point: %s, Id: %s, CityKey: %s, LocationType: %s, Type: %s, Payload: %s}",
		g.CreatedAt,
		g.GeoJson,
		g.Point,
		g.Id,
		g.CityKey,
		g.LocationType,
		g.Type,
		g.Payload,
	)
}
