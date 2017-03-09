package way

import (
	"github.com/brechtvm/osm/point"
)

func (w *Way) Centroid() *point.Point {
	if !w.Closed() {
		return nil
	}

	var cx, cy, a float64
	l := len(w.Nodes_) - 1
	for i := 0; i < l; i++ {
		ap := w.Nodes_[i].Position_.Lon*w.Nodes_[i+1].Position_.Lat - w.Nodes_[i+1].Position_.Lon*w.Nodes_[i].Position_.Lat
		cx += (w.Nodes_[i].Position_.Lon + w.Nodes_[i+1].Position_.Lon) * ap
		cy += (w.Nodes_[i].Position_.Lat + w.Nodes_[i+1].Position_.Lat) * ap
		a += ap
	}
	cx = 2 * cx / (6 * a)
	cy = 2 * cy / (6 * a)
	return point.New(cy, cx)
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
