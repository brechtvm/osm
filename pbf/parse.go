package pbf

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/brechtbm/osmpbf"
	"github.com/brechtvm/osm"
	"github.com/brechtvm/osm/bbox"
	"github.com/brechtvm/osm/item"
	"github.com/brechtvm/osm/node"
	"github.com/brechtvm/osm/point"
	"github.com/brechtvm/osm/relation"
	"github.com/brechtvm/osm/tags"
	"github.com/brechtvm/osm/user"
	"github.com/brechtvm/osm/way"
	"io"
	"log"
	"math"
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
func (p *Pbf) Parse(handler osm.OSMReader) (o *osm.OSM, err error) {
	//tagsToKeep := []string{"maxspeed", "highway"}

	d := osmpbf.NewDecoder(p.r)

	//err = d.Start(runtime.GOMAXPROCS(-1))
	err = d.Start(runtime.GOMAXPROCS(1))
	if err != nil {
		return
	}

	o = osm.NewOSM(handler)
	o.Users = make(map[uint32]*user.User)
	o.Timestamps = make(map[string]time.Time)
	var v interface{}

	lowerlat := math.MaxFloat32
	lowerlon := math.MaxFloat32
	upperlat := -math.MaxFloat32
	upperlon := -math.MaxFloat32
	counter := 0
	prevtime := time.Now()
	for {
		if v, err = d.Decode(); err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		} else {

			counter++
			if counter%1000000 == 0 {
				debug.FreeOSMemory()
				os.Stderr.WriteString(fmt.Sprintf("freeosmem  %5.1f s\n", time.Now().Sub(prevtime).Seconds()))
				prevtime = time.Now()
			}

			switch v := v.(type) {
			case *osmpbf.Node:

				t := tags.Tags(v.Tags)
				// UserInfo

				newnode := &node.Node{
					Id_:        v.ID,
					User_:      user.New(uint32(v.Info.Uid), v.Info.User), //o.Users[uint32(v.Info.Uid)],
					Position_:  point.New(v.Lat, v.Lon),
					Timestamp_: v.Info.Timestamp, //o.Timestamps[sTimestamp],
					Changeset_: v.Info.Changeset,
					Version_:   uint16(v.Info.Version),
					Visible_:   v.Info.Visible,
					Tags_:      &t,
				}

				lowerlat = math.Min(lowerlat, newnode.Position_.Lat)
				upperlat = math.Max(upperlat, newnode.Position_.Lat)
				lowerlon = math.Min(lowerlon, newnode.Position_.Lon)
				upperlon = math.Max(upperlon, newnode.Position_.Lon)

				if o.Handler != nil {
					if o.Handler.ReadNode(newnode) == false {
						return
					}
				} else {

					if _, ok := o.Users[uint32(v.Info.Uid)]; !ok {
						o.Users[uint32(v.Info.Uid)] = user.New(uint32(v.Info.Uid), v.Info.User)
					}

					//Timestamps
					sTimestamp := fmt.Sprintf("%s", v.Info.Timestamp)
					if _, ok := o.Timestamps[sTimestamp]; !ok {
						o.Timestamps[sTimestamp] = v.Info.Timestamp
					}
					o.Nodes[v.ID] = newnode
				}
			case *osmpbf.Way:
				t := tags.Tags(v.Tags)
				w := &way.Way{
					Id_:        v.ID,
					User_:      user.New(uint32(v.Info.Uid), v.Info.User),
					Timestamp_: v.Info.Timestamp,
					Changeset_: v.Info.Changeset,
					Version_:   uint16(v.Info.Version),
					Visible_:   v.Info.Visible,
					Tags_:      &t,
					NodeIDs:    v.NodeIDs,
				}

				if o.Handler != nil {
					if o.Handler.ReadWay(w) == false {
						return
					}
				} else {
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
				}
			case *osmpbf.Relation:
				t := tags.Tags(v.Tags)

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
						panic("mannekes toch, dit mag niet he!")
						return
					}
					members = append(members, member)
				}
				r.Members_ = members
				if o.Handler != nil {
					if o.Handler.ReadRelation(r) == false {
						return
					}
				} else {
					o.Relations[v.ID] = r
				}

			default:
				log.Printf("ERROR: unknown type %T\n", v)
			}
		}
	}

	if o.Handler != nil {
		if o.Handler.ReadBounds(&bbox.BBox{
			LowerLeft:  point.New(lowerlat, lowerlon),
			UpperRight: point.New(upperlat, upperlon),
		}) == false {
			return
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
