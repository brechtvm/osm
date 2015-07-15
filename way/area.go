package way

import (
	"math"
	"github.com/vetinari/osm/node"
)

func (w *Way) area() float64 {
	i := len(w.Nodes_) - 2
	var a float64
	for n := 0; n < i; n++ {
		// a += (deg2rad(w.Nodes[n].Position_.Lon)*deg2rad(w.Nodes[n+1].Position_.Lat) -
		// 		deg2rad(w.Nodes[n].Position_.Lat)*deg2rad(w.Nodes[n+1].Position_.Lon))
		a += (w.Nodes_[n].Position_.Lon*w.Nodes_[n+1].Position_.Lat -
			w.Nodes_[n].Position_.Lat*w.Nodes_[n+1].Position_.Lon)
	}
	return a / 2
}

// Returns the area of the (closed) way. Note, that this area is in
// square degrees
func (w *Way) Area() float64 {
	if !w.Closed() {
		return -1.0
	}
	return math.Abs(w.area())
}

// http://alienryderflex.com/polygon/
func (w *Way) Contains(n *node.Node) bool {
	odd := false
	j := len(w.Nodes_) - 2
	for i := 0; i < len(w.Nodes_)-2; i++ {
		if (w.Nodes_[i].Position_.Lon < n.Position_.Lon && w.Nodes_[j].Position_.Lon >= n.Position_.Lon) ||
			(w.Nodes_[j].Position_.Lon < n.Position_.Lon && w.Nodes_[i].Position_.Lon >= n.Position_.Lon) {
			if n.Position_.Lat >
				w.Nodes_[i].Position_.Lat+
					(n.Position_.Lon-w.Nodes_[i].Position_.Lon)/
						(w.Nodes_[j].Position_.Lon-w.Nodes_[i].Position_.Lon)*
							(w.Nodes_[j].Position_.Lat-w.Nodes_[i].Position_.Lat) {
				odd = !odd
			}
		}
		j = i
	}
	return odd
}

/*
func (w *Way) Contains(n *Node) bool {
	if !w.Closed() {
		return false
	}
	inside := false
	p0 := w.Nodes_[0]
	l := len(w.Nodes_) - 1
	for i := 1; i < l; i++ {
		pi := w.Nodes_[i]
		if n.Position_.Lon == p0.Position_.Lon &&
			p0.Position_.Lon == pi.Position_.Lon &&
			(n.Position_.Lat >= n.Position_.Lat || n.Position_.Lat >= pi.Position_.Lat) &&
			(n.Position_.Lat <= n.Position_.Lat || n.Position_.Lat <= pi.Position_.Lat) {

			return true
		}
		if p0.Position_.Lat == pi.Position_.Lat ||
			(n.Position_.Lat <= p0.Position_.Lat && n.Position_.Lat <= pi.Position_.Lat) ||
			(n.Position_.Lat > p0.Position_.Lat && n.Position_.Lat > pi.Position_.Lat) ||
			(n.Position_.Lon > p0.Position_.Lon && n.Position_.Lon > pi.Position_.Lon) {
			p0 = pi
			continue
		}
		if p0.Position_.Lon == pi.Position_.Lon ||
			n.Position_.Lon <= ((n.Position_.Lat-p0.Position_.Lat)*(pi.Position_.Lon-p0.Position_.Lon)/
				(pi.Position_.Lat-p0.Position_.Lat)+p0.Position_.Lon) {
			inside = !inside
		}
		p0 = pi
	}
	return inside
}
*/

// vim: ts=4 sw=4 noexpandtab nolist syn=go
