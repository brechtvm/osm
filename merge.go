package osm

import ()

func (self *OSM) Merge(other *OSM) {
	for _, n := range other.Nodes {
		nid := n.Id()
		if self.Nodes[nid] == nil {
			self.Nodes[nid] = n
			continue
		}
		if self.Nodes[nid].MergeOther(n) {
			self.Nodes[nid] = n
		}
	}
	for _, w := range other.Ways {
		wid := w.Id()
		if self.Ways[wid] == nil {
			self.Ways[wid] = w
			continue
		}
		if self.Ways[wid].MergeOther(w) {
			self.Ways[wid] = w
		}
	}
	for _, r := range other.Relations {
		rid := r.Id()
		if self.Relations[rid] == nil {
			self.Relations[rid] = r
			continue
		}
		if self.Relations[rid].MergeOther(r) {
			self.Relations[rid] = r
		}
	}
}
