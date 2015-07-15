package osm

import (
	"github.com/vetinari/osm/node"
	"github.com/vetinari/osm/relation"
	"github.com/vetinari/osm/tags"
	"github.com/vetinari/osm/way"
)

func (o *OSM) FilterTags(tl... *tags.Tags) *OSM {
	no := NewOSM()
	for id := range o.Nodes {
		n := o.GetNode(id)
		if n == nil {
			continue
		}
		if n.Tags() == nil {
			continue
		}

		nt := map[string]string(*n.Tags())
	findNodeTags:
		for _, t := range tl {
			for k := range nt {
				if t.Has(k) && t.Get(k) == nt[k] {
					no.Nodes[id] = n
					break findNodeTags
				}
			}
		}
	}

	for id := range o.Ways {
		n := o.GetWay(id)
		if n == nil {
			continue
		}
		if n.Tags() == nil {
			continue
		}

		nt := map[string]string(*n.Tags())
	findWayTags:
		for _, t := range tl {
			for k := range nt {
				if t.Has(k) && t.Get(k) == nt[k] {
					no.Ways[id] = n
					for _, nd := range n.GetNodes() {
						no.Nodes[n.Id_] = nd
					}
					break findWayTags
				}
			}
		}
	}

	for id := range o.Relations {
		n := o.GetRelation(id)
		if n == nil {
			continue
		}
		if n.Tags() == nil {
			continue
		}

		nt := map[string]string(*n.Tags())
	findRelationTags:
		for _, t := range tl {
			for k := range nt {
				if t.Has(k) && t.Get(k) == nt[k] {
					no.Relations[id] = n
					for _, m := range n.GetMembers() {
						switch m.Ref.(type) {
						case *node.Node:
							no.Nodes[(m.Ref).(*node.Node).Id_] = (m.Ref).(*node.Node)
						case *way.Way:
							no.Ways[(m.Ref).(*way.Way).Id_] = (m.Ref).(*way.Way)
							for _, nd := range (m.Ref).(*way.Way).GetNodes() {
								node := o.GetNode(nd.Id_)
								if node != nil {
									no.Nodes[node.Id_] = node
								}
							}
						case *relation.Relation:
							no.Relations[(m.Ref).(*relation.Relation).Id_] = (m.Ref).(*relation.Relation)
							for _, nd := range (m.Ref).(*relation.Relation).GetNodes() {
								node := o.GetNode(nd.Id_)
								if node != nil {
									no.Nodes[node.Id_] = node
								}
							}
							for _, wy := range (m.Ref).(*relation.Relation).GetWays() {
								w := o.GetWay(wy.Id_)
								if w != nil {
									no.Ways[w.Id_] = w
									for _, nd := range (m.Ref).(*way.Way).GetNodes() {
										node := o.GetNode(nd.Id_)
										if node != nil {
											no.Nodes[node.Id_] = node
										}
									}
								}
							}
							// FIXME - add relations recursively
						}
					}
					break findRelationTags
				}
			}
		}
	}
	return no
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
