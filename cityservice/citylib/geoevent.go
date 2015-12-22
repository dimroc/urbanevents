package citylib

import (
	"encoding/json"
	"fmt"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	elastigo "github.com/mattbaird/elastigo/lib"
	"log"
	"strings"
	"time"
)

type GeoEvent struct {
	CityKey       string     `json:"city"`
	CreatedAt     time.Time  `json:"createdAt"`
	FullName      string     `json:"fullName"`
	GeoJson       GeoJson    `json:"geojson"`
	Hashtags      []string   `json:"hashtags"`
	Id            string     `json:"id"`
	MediaUrl      string     `json:"mediaUrl"`
	Link          string     `json:"link"`
	LocationType  string     `json:"locationType"`
	MediaType     string     `json:"mediaType"`
	Text          string     `json:"text,omitempty"`
	TextFrench    string     `json:"text_fr,omitempty"`
	Point         [2]float64 `json:"point"`
	Service       string     `json:"service"`
	ThumbnailUrl  string     `json:"thumbnailUrl"`
	Type          string     `json:"type"`
	Username      string     `json:"username"`
	Place         string     `json:"place"`
	Neighborhoods []string   `json:"neighborhoods"`
	ExpandedUrl   string     `json:"-"`
}

type GeoJson struct {
	Type           string           `json:"type"`
	CoordinatesRaw *json.RawMessage `json:"coordinates"` // Coordinate always has to have exactly 2 values
}

func GeoEventsFromHits(hits *elastigo.Hits) []GeoEvent {
	response := []GeoEvent{}

	for _, hit := range hits.Hits {
		geoevent := GeoEvent{}
		err := json.Unmarshal(*hit.Source, &geoevent)
		Check(err)

		response = append(response, geoevent)
	}

	return response
}

func GeoEventsFromElasticSearch(result *elastigo.SearchResult) []GeoEvent {
	return GeoEventsFromHits(&result.Hits)
}

func (g GeoJson) Center() [2]float64 {
	return g.GenerateShape().Center()
}

func (g GeoJson) String() string {
	j, err := json.Marshal(&g.CoordinatesRaw)
	if err != nil {
		log.Panic(err)
	}

	return fmt.Sprintf(
		"{Type: %s, Coordinates: %s}",
		g.Type,
		string(j),
	)
}

func (g GeoJson) TryCollapseEmptyBoundingBox() GeoJson {
	if strings.ToLower(g.Type) == "polygon" {
		shape := g.GenerateShape()
		bb := shape.(*BoundingBox)

		if bb.Collapsed() {
			Logger.Debug("Collapsing %s: ", bb)
			return GeoJsonFrom("Point", [2]float64{bb.Coordinates[0][0][0], bb.Coordinates[0][0][1]})
		}
	}

	return g
}

type GeoShape interface {
	Center() [2]float64
}

type Point struct {
	GeoJson
	Coordinates [2]float64 `json:"coordinates"` // Coordinate always has to have exactly 2 values
}

func (p *Point) Center() [2]float64 {
	return p.Coordinates
}

type BoundingBox struct {
	GeoJson
	Coordinates [][][]float64 `json:"coordinates"`
}

func (bb *BoundingBox) Center() [2]float64 {
	long := (bb.Coordinates[0][0][0] + bb.Coordinates[0][1][0]) / 2.0
	lat := (bb.Coordinates[0][0][1] + bb.Coordinates[0][1][1]) / 2.0
	center := [2]float64{long, lat}
	return center
}

func (bb *BoundingBox) Collapsed() bool {
	return bb.Coordinates[0][0][0] == bb.Coordinates[0][1][0] &&
		bb.Coordinates[0][0][0] == bb.Coordinates[0][2][0] &&
		bb.Coordinates[0][0][0] == bb.Coordinates[0][3][0]
}

func GeoJsonFrom(typeValue string, v interface{}) GeoJson {
	b, err := json.Marshal(v)
	Check(err)

	geojson := GeoJson{
		Type:           typeValue,
		CoordinatesRaw: (*json.RawMessage)(&b),
	}

	return geojson
}

func (g *GeoEvent) String() string {
	return fmt.Sprintf(
		"{CreatedAt: %s, GeoJson: %s, Point: [%s,%s], Id: %s, CityKey: %s, LocationType: %s, Type: %s, Text: %s}",
		g.CreatedAt,
		g.GeoJson.String(),
		g.Point[0],
		g.Point[1],
		g.Id,
		g.CityKey,
		g.LocationType,
		g.Type,
		g.Text,
	)
}

func (geojson *GeoJson) GenerateShape() GeoShape {
	var shape GeoShape
	var coordinatesDestination interface{}

	switch strings.ToLower(geojson.Type) {
	case "point":
		point := &Point{GeoJson: *geojson}
		coordinatesDestination = &point.Coordinates
		shape = point
	case "polygon":
		box := &BoundingBox{GeoJson: *geojson}
		coordinatesDestination = &box.Coordinates
		shape = box
	}

	err := json.Unmarshal(*geojson.CoordinatesRaw, coordinatesDestination)
	if err != nil {
		Logger.Critical("%s", err)
		panic(err)
	}

	return shape
}
