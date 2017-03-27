package pbf

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/brechtbm/osmpbf"
	"github.com/brechtvm/osm"
	"github.com/brechtvm/osm/item"
	"github.com/brechtvm/osm/node"
	"github.com/brechtvm/osm/point"
	"github.com/brechtvm/osm/relation"
	"github.com/brechtvm/osm/tags"
	"github.com/brechtvm/osm/user"
	"github.com/brechtvm/osm/way"
	"io"
	"log"
	"os"
	"runtime"
)

// returns an osm.Parser which can be used as argument to osm.New(), reads
// from the given file
func FileParser(file string) (osm.Parser, io.Closer, error) {
	fh, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	return Parser(fh), fh, nil
}

// returns an osm.Parser which can be used as argument to osm.New(), reads
// from byte array
func ByteParser(data []byte) osm.Parser {
	return Parser(bytes.NewReader(data))
}

// returns an osm.Parser which can be used as argument to osm.New()
func Parser(r io.Reader) osm.Parser {
	return &Pbf{r}
}

type Pbf struct {
	r io.Reader
}

// implements the osm.Parser interface
func (p *Pbf) Parse() (o *osm.OSM, err error) {
	d := osmpbf.NewDecoder(p.r)

	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		return
	}

	o = osm.NewOSM()
	var v interface{}
	for {
		if v, err = d.Decode(); err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				var t tags.Tags
				if v.Tags != nil && len(v.Tags) != 0 {
					t = tags.Tags(v.Tags)
				}
				o.Nodes[v.ID] = &node.Node{
					Id_:        v.ID,
					User_:      user.New(int64(v.Info.Uid), v.Info.User),
					Position_:  point.New(v.Lat, v.Lon),
					Timestamp_: v.Info.Timestamp,
					Changeset_: v.Info.Changeset,
					Version_:   int64(v.Info.Version),
					Visible_:   v.Info.Visible,
					Tags_:      &t,
				}
			case *osmpbf.Way:
				var t tags.Tags
				if v.Tags != nil && len(v.Tags) != 0 {
					t = tags.Tags(v.Tags)
				}
				w := &way.Way{
					Id_:        v.ID,
					User_:      user.New(int64(v.Info.Uid), v.Info.User),
					Timestamp_: v.Info.Timestamp,
					Changeset_: v.Info.Changeset,
					Version_:   int64(v.Info.Version),
					Visible_:   v.Info.Visible,
					Tags_:      &t,
				}
				var nd []*node.Node
				for _, id := range v.NodeIDs {
					n := o.Nodes[id]
					if n == nil {
						err = errors.New(fmt.Sprintf("Missing node #%d in way #%d", id, v.ID))
						o = nil
						return
					}
					nd = append(nd, n)
				}
				w.Nodes_ = nd
				o.Ways[v.ID] = w
			case *osmpbf.Relation:
				var t tags.Tags
				if v.Tags != nil && len(v.Tags) != 0 {
					t = tags.Tags(v.Tags)
				}
				r := &relation.Relation{
					Id_:        v.ID,
					User_:      user.New(int64(v.Info.Uid), v.Info.User),
					Timestamp_: v.Info.Timestamp,
					Changeset_: v.Info.Changeset,
					Version_:   int64(v.Info.Version),
					Visible_:   v.Info.Visible,
					Tags_:      &t,
				}
				var members []*relation.Member
				for _, m := range v.Members {
					member := &relation.Member{Role: m.Role, Id_: m.ID}
					switch m.Type {
					case osmpbf.NodeType:
						member.Type_ = item.TypeNode
						member.Ref = o.GetNode(m.ID)
					case osmpbf.WayType:
						member.Type_ = item.TypeWay
						member.Ref = o.GetWay(m.ID)
					case osmpbf.RelationType:
						member.Type_ = item.TypeRelation
						member.Ref = o.GetRelation(m.ID)
					}
					if member.Ref == nil {
						err = errors.New(fmt.Sprintf("Missing member #%d (%s) in way #%d", m.ID, member.Type(), v.ID))
						o = nil
						return
					}
					members = append(members, member)
				}
				r.Members_ = members
				o.Relations[v.ID] = r
			default:
				log.Printf("ERROR: unknown type %T\n", v)
			}
		}
	}
	return
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
