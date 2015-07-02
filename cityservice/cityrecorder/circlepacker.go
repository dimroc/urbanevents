package cityrecorder

import (
	geo "github.com/kellydunn/golang-geo"
)

const (
	diagonalDistance float64 = 11.18033988749895
)

type GeoCircle struct {
	Point  [2]float64 `json:"point"`
	Radius int        `json:"radius"`
}

/*
* We'll be doing Hexagonal Packing:
* https://en.wikipedia.org/wiki/Close-packing_of_equal_spheres
* Now, form the next row of spheres. Again, the centers will all lie on a straight line with x-coordinate differences of 2r,
* but there will be a shift of distance r in the x-direction so that the center of every sphere in this row aligns with
* the x-coordinate of where two spheres touch in the first row. This allows the spheres of the new row to slide in closer to
* the first row until all spheres in the new row are touching two spheres of the first row. Since the new spheres touch two
* spheres, their centers form an equilateral triangle with those two neighbors' centers. The side lengths are all 2r, so the
* height or y-coordinate difference between the rows is \scriptstyle\sqrt{3}r.
 */
func PackCircles(boundingBox [4]float64) []GeoCircle {
	circles := []GeoCircle{}
	// Create first circle in bottom left of bb.
	bottomLeft := geo.NewPoint(boundingBox[1], boundingBox[0]) // lat then long unlike geoJSON
	hypotenuse := 7.0710678118                                 // 5km east and 5km north = √(25+25) = 7.0710678118

	start := bottomLeft.PointAtDistanceAndBearing(hypotenuse, 45)
	rowOrigin := start
	cursor := rowOrigin

	top := boundingBox[3]
	right := boundingBox[2]

	// Large circles
	for cursor.Lat() < top {
		// Rows
		for cursor.Lng() < right {
			// Columns
			circles = append(circles, circleFromPoint(cursor, 5)) // 5km because of instagram limit
			cursor = moveRight(cursor)
		}

		rowOrigin = moveUp(rowOrigin)
		cursor = rowOrigin
	}

	// Fill in the gaps
	rowOrigin = start.PointAtDistanceAndBearing(hypotenuse, 45)
	cursor = rowOrigin
	for cursor.Lat() < top {
		// Rows
		for cursor.Lng() < right {
			// Columns
			circles = append(circles, circleFromPoint(cursor, 5))
			cursor = moveRight(cursor)
		}

		rowOrigin = moveUp(rowOrigin)
		cursor = rowOrigin
	}

	return circles
}

func circleFromPoint(p *geo.Point, radius int) GeoCircle {
	return GeoCircle{
		Point:  [2]float64{p.Lng(), p.Lat()},
		Radius: radius,
	}
}

func moveDiagonalLeft(point *geo.Point) *geo.Point {
	return point.PointAtDistanceAndBearing(10, -26.56494984)
}

func moveDiagonalRight(point *geo.Point) *geo.Point {
	return point.PointAtDistanceAndBearing(10, 26.56494984)
}

func moveUp(point *geo.Point) *geo.Point {
	return point.PointAtDistanceAndBearing(10, 0)
}

func moveRight(point *geo.Point) *geo.Point {
	return point.PointAtDistanceAndBearing(10, 90)
}
