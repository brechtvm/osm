package xml

import (
	"fmt"
	"io"
	"github.com/vetinari/osm"
	"github.com/vetinari/osm/node"
	"github.com/vetinari/osm/relation"
	"github.com/vetinari/osm/way"
	"sort"
)

var xmlWriterVersion = "1.0"
// Dump() is not suitable for uploading: modified items still have the same version
func Dump(w io.Writer, o *osm.OSM) {
	w.Write([]byte("<?xml version='1.0' encoding='UTF-8'?>\n"))
	w.Write([]byte(fmt.Sprintf("<osm version='0.6' generator='osm/xml/write.go v%s'>\n", xmlWriterVersion)))

	bb, err := o.BoundingBox()
	if err == nil {
		w.Write([]byte(bb.String()))
	}

	nl := o.GetNodeList()
	sort.Sort(nl)
	for _, n := range []*node.Node(*nl) {
		w.Write([]byte(n.String()))
	}

	wl := o.GetWayList()
	sort.Sort(wl)
	for _, wy := range []*way.Way(*wl) {
		w.Write([]byte(wy.String()))
	}

	rl := o.GetRelationList()
	if len([]*relation.Relation(*rl)) != 0 {
		sort.Sort(rl)
	}
	for _, r := range []*relation.Relation(*rl) {
		w.Write([]byte(r.String()))
	}

	w.Write([]byte("</osm>\n"))
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
