package cityrecorder

import ()

//"github.com/kellydunn/golang-geo"

type GeoCircle struct {
	Point  []float64
	Radius int
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
	return []GeoCircle{}
}
