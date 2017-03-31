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
	"runtime/debug"
	"time"
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
	tagsToKeep := []string{"maxspeed", "highway"}

	d := osmpbf.NewDecoder(p.r)

	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		return
	}

	o = osm.NewOSM()
	o.Users = make(map[uint32]*user.User)
	o.Timestamps = make(map[string]time.Time)
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
				counter := len(o.Nodes)
				if counter%1000000 == 0 {
					log.Printf("Processed %d nodes ", counter)
					log.Printf("[%d]users", len(o.Users))
				}
				if counter%50000000 == 0 {
					debug.FreeOSMemory()
					log.Printf("Garbage collected!")
				}
				var t tags.Tags
				if v.Tags != nil && len(v.Tags) != 0 {
					subsetTags := make(map[string]string)
					for tagN, tagV := range v.Tags {
						for _, tag := range tagsToKeep {
							if tagN == tag {
								subsetTags[tagN] = tagV
							}
						}
					}
					t = tags.Tags(subsetTags)
					subsetTags = nil
				}

				// UserInfo
				if _, ok := o.Users[uint32(v.Info.Uid)]; !ok {
					o.Users[uint32(v.Info.Uid)] = user.New(uint32(v.Info.Uid), v.Info.User)
				}

				// Timestamps
				sTimestamp := fmt.Sprintf("%s", v.Info.Timestamp)
				if _, ok := o.Timestamps[sTimestamp]; !ok {
					o.Timestamps[sTimestamp] = v.Info.Timestamp
				}

				o.Nodes[v.ID] = &node.Node{
					Id_:        v.ID,
					User_:      o.Users[uint32(v.Info.Uid)],
					Position_:  point.New(v.Lat, v.Lon),
					Timestamp_: o.Timestamps[sTimestamp],
					Changeset_: v.Info.Changeset,
					Version_:   uint16(v.Info.Version),
					Visible_:   v.Info.Visible,
					Tags_:      &t,
				}
			case *osmpbf.Way:
				counter := len(o.Ways)
				if counter == 0 {
					log.Println("Processing ways")
				}
				var t tags.Tags
				if v.Tags != nil && len(v.Tags) != 0 {
					subsetTags := make(map[string]string)
					keep := false
					for tagN, _ := range v.Tags {
						if tagN == "highway" {
							keep = true
						}
						/*
							for _, tag := range tagsToKeep {
								if tagN == tag {
									subsetTags[tagN] = tagV
								}
							}
						*/
					}
					if !keep {
						fmt.Printf("Skipping way[%d]\n", v.ID)
						continue
					}
					t = tags.Tags(subsetTags)
				}
				w := &way.Way{
					Id_:        v.ID,
					User_:      user.New(uint32(v.Info.Uid), v.Info.User),
					Timestamp_: v.Info.Timestamp,
					Changeset_: v.Info.Changeset,
					Version_:   uint16(v.Info.Version),
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
				counter := len(o.Relations)
				if counter == 0 {
					log.Println("Processing relations")
				}
				var t tags.Tags
				if v.Tags != nil && len(v.Tags) != 0 {
					subsetTags := make(map[string]string)
					for tagN, tagV := range v.Tags {
						for _, tag := range tagsToKeep {
							if tagN == tag {
								subsetTags[tagN] = tagV
							}
						}
					}
					t = tags.Tags(subsetTags)
				}
				r := &relation.Relation{
					Id_:        v.ID,
					User_:      user.New(uint32(v.Info.Uid), v.Info.User),
					Timestamp_: v.Info.Timestamp,
					Changeset_: v.Info.Changeset,
					Version_:   uint16(v.Info.Version),
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
// vim: ts=4 sw=4 noexpandtab nolist syn=go
// vim: ts=4 sw=4 noexpandtab nolist syn=go
// vim: ts=4 sw=4 noexpandtab nolist syn=go
// vim: ts=4 sw=4 noexpandtab nolist syn=go
// vim: ts=4 sw=4 noexpandtab nolist syn=go
