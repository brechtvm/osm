package way

import (
	"github.com/vetinari/osm/distance"
)

// sum of all node distances along the way
func (w *Way) Length() distance.Distance {
	var l distance.Distance
	n := len(w.Nodes_) - 1
	for i := 0; i < n; i++ {
		l += w.Nodes_[i].Position_.DistanceOf(w.Nodes_[i+1].Position_)
	}
	return l
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
