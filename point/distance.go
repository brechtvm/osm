package point

import (
	"github.com/brechtvm/osm/distance"
	"math"
)

// the Great Circle Distance
func (p *Point) DistanceOf(q *Point) distance.Distance {
	// convert degree -> rad
	plat := p.Lat * math.Pi / 180.0
	plon := p.Lon * math.Pi / 180.0

	qlat := q.Lat * math.Pi / 180.0
	qlon := q.Lon * math.Pi / 180.0

	a := math.Cos(qlat) * math.Sin(plon-qlon)
	b := math.Cos(plat)*math.Sin(qlat) - math.Sin(plat)*math.Cos(qlat)*math.Cos(plon-qlon)
	c := math.Sin(plat)*math.Sin(qlat) + math.Cos(plat)*math.Cos(qlat)*math.Cos(plon-qlon)

	dist := math.Atan2(math.Sqrt(a*a+b*b), c)
	if dist < 0.0 {
		dist += 2 * math.Pi
	}
	return distance.Distance(dist) * distance.EarthRadius
}

// and another way to compute the great circle distance
func (p *Point) DistanceOf2(q *Point) distance.Distance {
	return distance.Distance(p.Angle(q)) * distance.EarthRadius
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
