package relation

import (
	"errors"
	"fmt"
	"github.com/brechtvm/osm/bbox"
	"github.com/brechtvm/osm/point"
	"github.com/brechtvm/osm/way"
)

// Weighted centroid of all area parts
//
// NOTE - when a relation has node members outside of the way members, this will
//        probably result in a wrong point (bounding box is taken from all nodes,
//        not just way nodes)
func (r *Relation) Centroid() (p *point.Point, err error) {
	if !r.IsAreaRelation() {
		err = errors.New("Not an area relation")
		return
	}

	var seen []*way.Way
	var bb *bbox.BBox
	var lat float64
	var lon float64
	var den float64

	wayList, err := r.WayMembersAsWays()
	if err != nil {
		return
	}
	bb, err = r.BoundingBox()
	if err != nil {
		return
	}
	for _, w := range wayList {
		c := w.Centroid()
		if c == nil {
			err = errors.New(fmt.Sprintf("Way #%d has no centroid?!", w.Id()))
			return
		}

		a := w.Area()
		if a == -1.0 {
			err = errors.New(fmt.Sprintf("Way #%d has no area", w.Id()))
			return
		}

		count := 0
		for _, s := range seen {
			if s.Contains(w.Nodes_[0]) {
				count += 1
			}
		}
		seen = append(seen, w)

		neg := 1.0
		if count%2 == 1 {
			neg = -1.0
		}

		lon += a * neg * (c.Lon - bb.LowerLeft.Lon)
		lat += a * neg * (c.Lat - bb.LowerLeft.Lat)
		den += a * neg
	}

	return point.New(bb.LowerLeft.Lat+lat/den, bb.LowerLeft.Lon+lon/den), nil
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
