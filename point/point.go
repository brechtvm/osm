package point

import (
	"math"
)

// the Name is currently only used in GPX(p *Point) where the
// Point p is a named point..
type Point struct {
	Lat  float64
	Lon  float64
	Name string
}

func New(lat, lon float64) *Point {
	return &Point{Lat: lat, Lon: lon}
}

// a point in cartesian system
type CPoint struct {
	x float64
	y float64
	z float64
}

func deg2rad(v float64) float64 {
	return v * math.Pi / 180.0
}

func rad2deg(v float64) float64 {
	return v * 180.0 / math.Pi
}

// Converts a spherical (i.e. normal OSM coordinates) to a point in cartesian
// coordinate system. Note that this point is on a unit sphere (radius == 1).
func (p *Point) Cartesian() *CPoint {
	rlat := deg2rad(p.Lat)
	rlon := deg2rad(p.Lon)
	return &CPoint{
		x: math.Cos(rlon) * math.Cos(rlat),
		y: math.Sin(rlon) * math.Cos(rlat),
		z: math.Sin(rlat),
	}
}

// Converts a point in cartesian coordinates to the normal OSM spherical
func (p *CPoint) Spherical() *Point {
	return &Point{
		Lat: rad2deg(math.Pi/2 - math.Acos(p.z/(math.Sqrt(p.x*p.x+p.y*p.y+p.z*p.z)))),
		Lon: rad2deg(math.Atan(p.y / p.x)),
	}
}

// The angle between two points (in the center of the sphere)
func (p *Point) Angle(q *Point) float64 {
	// http://www.intmath.com/vectors/7-vectors-in-3d-space.php
	cp := p.Cartesian()
	cq := q.Cartesian()

	th := math.Acos(
		(cp.x*cq.x + cp.y*cq.y + cp.z*cq.z) /
			(math.Sqrt(cp.x*cp.x+cp.y*cp.y+cp.z*cp.z) *
				math.Sqrt(cq.x*cq.x+cq.y*cq.y+cq.z*cq.z)))
	// fmt.Printf("th=%f => %f => %f\n", th, th * 180.0 / math.Pi, th * float64(EarthRadius))
	return th
}

func (p *Point) Equal(q *Point) bool {
	return p.Lat == q.Lat && p.Lon == q.Lon
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
