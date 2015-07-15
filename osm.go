package osm

import (
	"fmt"
	"github.com/vetinari/osm/bbox"
	"github.com/vetinari/osm/node"
	"github.com/vetinari/osm/relation"
	"github.com/vetinari/osm/way"
	"sort"
)

// the osm.Parser interface, currently implemented by the xml and pbf sub modules
type Parser interface {
	Parse() (*OSM, error)
}

// returns a new OSM, which is filled by the passed osm.Parser
func New(p Parser) (*OSM, error) {
	return p.Parse()
}

// returns a new and empty OSM
func NewOSM() *OSM {
	return &OSM{
		Version:   "0.6",
		Nodes:     make(map[int64]*node.Node),
		Ways:      make(map[int64]*way.Way),
		Relations: make(map[int64]*relation.Relation),
	}
}

// the main entry point for OSM data
type OSM struct {
	Version   string
	BBox      bbox.BBox
	Origin    string
	Nodes     map[int64]*node.Node
	Ways      map[int64]*way.Way
	Relations map[int64]*relation.Relation
}

func (o *OSM) BoundingBox() (*bbox.BBox, error) {
	return o.GetNodeList().BoundingBox()
}

func (o *OSM) GetNode(id int64) *node.Node {
	return o.Nodes[id]
}

func (o *OSM) GetWay(id int64) *way.Way {
	return o.Ways[id]
}

func (o *OSM) GetRelation(id int64) *relation.Relation {
	return o.Relations[id]
}

func (o *OSM) GetNodeList() *node.NodeList {
	var nl []*node.Node
	for _, n := range o.Nodes {
		nl = append(nl, n)
	}
	nlist := node.NodeList(nl)
	return &nlist
}

func (o *OSM) GetWayList() *way.WayList {
	var wl []*way.Way
	for _, w := range o.Ways {
		wl = append(wl, w)
	}
	wlist := way.WayList(wl)
	return &wlist
}

func (o *OSM) GetRelationList() *relation.RelationList {
	var rl []*relation.Relation
	for _, r := range o.Relations {
		rl = append(rl, r)
	}
	rlist := relation.RelationList(rl)
	return &rlist
}

var osmStringVersion = "0.1"

// returns a stringified (in osm XML format) version of the OSM. Do not 
// use for huge OSM data, better dump to some file via xml.Dump()
func (o *OSM) String() string {
	xml := "<?xml version='1.0' encoding='UTF-8'?>\n" +
		fmt.Sprintf("<osm version='0.6' upload='true' generator='osm.String v%s'>\n", osmStringVersion)
	bb, err := o.BoundingBox()
	if err == nil {
		xml += bb.String()
	}

	nl := o.GetNodeList()
	sort.Sort(nl)
	for _, n := range []*node.Node(*nl) {
		xml += n.String()
	}

	wl := o.GetWayList()
	sort.Sort(wl)
	for _, w := range []*way.Way(*wl) {
		xml += w.String()
	}

	rl := o.GetRelationList()
	if len([]*relation.Relation(*rl)) != 0 {
		sort.Sort(rl)
	}
	for _, r := range []*relation.Relation(*rl) {
		xml += r.String()
	}

	xml += "</osm>\n"
	return xml
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
