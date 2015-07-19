package relation

import (
	// "fmt"
	"github.com/vetinari/osm/bbox"
	"github.com/vetinari/osm/item"
	"github.com/vetinari/osm/node"
	"github.com/vetinari/osm/tags"
	"github.com/vetinari/osm/user"
	"github.com/vetinari/osm/way"
	"time"
)

func (self *Relation) Type() item.ItemType  { return item.TypeRelation }
func (self *Relation) Id() int64            { return self.Id_ }
func (self *Relation) Members() []*Member   { return self.Members_ }
func (self *Relation) User() *user.User     { return self.User_ }
func (self *Relation) Tags() *tags.Tags     { return self.Tags_ }
func (self *Relation) Timestamp() time.Time { return self.Timestamp_ }
func (self *Relation) Version() int64       { return self.Version_ }
func (self *Relation) Changeset() int64     { return self.Changeset_ }
func (self *Relation) Visible() bool        { return self.Visible_ }

var newRelationNum int64 = 0

func newRelationId() int64 {
	newRelationNum -= 1
	return newRelationNum
}

type Member struct {
	Type_ item.ItemType
	Role  string
	Ref   item.Item
	Id_   int64
}

type Relation struct {
	Id_        int64
	Members_   []*Member
	User_      *user.User
	Tags_      *tags.Tags
	Timestamp_ time.Time
	Version_   int64
	Changeset_ int64
	Visible_   bool
	modified   bool
	deleted    bool
}

type RelationList []*Relation

func NewRelation(m *Member) *Relation {
	return &Relation{
		Members_:   []*Member{m},
		Id_:        newRelationId(),
		Tags_:      tags.New(),
		Timestamp_: time.Now(),
		Version_:   0,
		Changeset_: 0,
		Visible_:   true,
		User_:      &user.User{Id: 0, Name: ""},
		modified:   true,
		deleted:    false,
	}
}

func (m *Member) Type() item.ItemType {
	return m.Type_
}

func (m *Member) Id() int64 {
	return m.Id_
}

func NewMember(role string, i item.Item) *Member {
	switch i.Type() {
	case item.TypeNode:
		return &Member{Type_: i.Type(), Role: role, Ref: i, Id_: i.(*node.Node).Id()}
	case item.TypeWay:
		return &Member{Type_: i.Type(), Role: role, Ref: i, Id_: i.(*way.Way).Id()}
	case item.TypeRelation:
		return &Member{Type_: i.Type(), Role: role, Ref: i, Id_: i.(*Relation).Id()}
	default:
		panic("invalid member type")
	}
}

func (r *Relation) GetMembers() []*Member {
	return r.Members_
}

func (r *Relation) AddMember(i item.Item, role string) {
	r.Members_ = append(r.Members_, NewMember(role, i))
}

func (r *Relation) GetNodes() []*node.Node {
	var n []*node.Node

	for _, m := range r.GetMembers() {
		switch (m.Ref).(type) {
		case *node.Node:
			n = append(n, (m.Ref).(*node.Node))
		case *way.Way:
			n = append(n, (m.Ref).(*way.Way).GetNodes()...)
		case *Relation:
			n = append(n, (m.Ref).(*Relation).GetNodes()...)
		}
	}
	return n
}

func (r *Relation) GetWays() []*way.Way {
	var w []*way.Way

	for _, m := range r.GetMembers() {
		switch (m.Ref).(type) {
		case *way.Way:
			w = append(w, (m.Ref).(*way.Way))
		case *Relation:
			w = append(w, (m.Ref).(*Relation).GetWays()...)
		default:
			continue
		}
	}
	return w
}

func (r *Relation) BoundingBox() (*bbox.BBox, error) {
	return node.NodeList(r.GetNodes()).BoundingBox()
}

// Recursively collects all members of type "way"
func (r *Relation) WayMembers() []*Member {
	var wm []*Member
	for _, m := range r.GetMembers() {
		switch (m.Ref).(type) {
		case *way.Way:
			wm = append(wm, m)
		case *Relation:
			wm = append(wm, (m.Ref).(*Relation).WayMembers()...)
		default:
			continue
		}
	}
	return wm
}

func (r *Relation) IsMultipolygon() bool {
	if r.Tags == nil {
		return false
	}
	t := map[string]string(*r.Tags_)
	return t["type"] == "multipolygon"
}

func (r *Relation) IsAreaRelation() bool {
	if r.Tags == nil {
		return false
	}
	switch r.Tags_.Get("type") {
	case "multipolygon", "boundary":
		return true
	default:
		return false
	}
	return false
}

// If a way in the r.WayMembers() output is connected to the next one in
// the list (and they must have the same role) they're joined, otherwise
// a new way is started.
//
// For area relations like "type=multipolygon" or "type=boundary" you can
// check if the member ways build a closed ring (or multiple rings) by
// running Cosed() for each way returned by this func.
func (r *Relation) WayMembersAsWays() ([]*way.Way, error) {
	var ways []*way.Way
	var err error
	all := r.WayMembers()
	if len(all) == 0 {
		return way.EmptyWays(), nil
	}
	if len(all) == 1 {
		// in doubt we get an empty list back
		n, err := way.New(all[0].Ref.(*way.Way).Nodes())
		if err != nil {
			return way.EmptyWays(), err
		}
		return []*way.Way{n}, nil
	}

	cur, all := all[0], all[1:]
	// prev := cur
	role := cur.Role
	cur_w, err := way.New(cur.Ref.(*way.Way).Nodes())
	if err != nil {
		return way.EmptyWays(), err
	}
	for _, wr := range all {
		if wr.Role != role {
			ways = append(ways, cur_w)
			cur_w, err = way.New(wr.Ref.(*way.Way).Nodes())
			if err != nil {
				return way.EmptyWays(), err
			}
			role = wr.Role
			// prev = wr
			continue
		}
		conn := cur_w.Connected(wr.Ref.(*way.Way))
		// fmt.Printf("Connected: #%d - #%d => %s\n", prev.Ref.(*way.Way).Id(), wr.Ref.(*way.Way).Id(), conn)
		var conn_way *way.Way
		switch conn {
		case way.NotConnected:
			ways = append(ways, cur_w)
			cur_w, err = way.New(wr.Ref.(*way.Way).Nodes())
			if err != nil {
				return way.EmptyWays(), err
			}
			role = wr.Role
			// prev = wr
			continue
		case way.ConnectedNormal:
			conn_way = wr.Ref.(*way.Way)
		case way.ConnectedReversed1st:
			cur_w.Reverse()
			conn_way, err = way.New(wr.Ref.(*way.Way).Nodes())
			if err != nil {
				return way.EmptyWays(), err
			}
		case way.ConnectedReversed2nd:
			conn_way, err = way.New(wr.Ref.(*way.Way).Nodes())
			if err != nil {
				return way.EmptyWays(), err
			}
			conn_way.Reverse()
		case way.ConnectedReversedBoth:
			cur_w.Reverse()
			conn_way, err = way.New(wr.Ref.(*way.Way).Nodes())
			if err != nil {
				return way.EmptyWays(), err
			}
			conn_way.Reverse()
		}
		_ = cur_w.Join(conn_way) // we already checked they're connected
		// prev = wr
	}
	return append(ways, cur_w), nil
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
