package node

import (
	"errors"
	"github.com/brechtvm/osm/bbox"
	"github.com/brechtvm/osm/item"
	"github.com/brechtvm/osm/point"
	"github.com/brechtvm/osm/tags"
	"github.com/brechtvm/osm/user"
	"math"
	"time"
)

type Node struct {
	Id_        int64
	Position_  *point.Point
	User_      *user.User
	Tags_      *tags.Tags
	Timestamp_ time.Time
	Version_   int64
	Changeset_ int64
	Visible_   bool
	modified   bool
	deleted    bool
}

// the NodeList implements the sort.Interface
type NodeList []*Node

// part of the item.Item interface
func (self *Node) Changeset() int64 { return self.Changeset_ }

// part of the item.Item interface
func (self *Node) Id() int64              { return self.Id_ }
func (self *Node) Position() *point.Point { return self.Position_ }

// part of the item.Item interface
func (self *Node) Tags() *tags.Tags { return self.Tags_ }

// part of the item.Item interface
func (self *Node) Timestamp() time.Time { return self.Timestamp_ }

// part of the item.Item interface
func (self *Node) Type() item.ItemType { return item.TypeNode }

// part of the item.Item interface
func (self *Node) User() *user.User { return self.User_ }

// part of the item.Item interface
func (self *Node) Version() int64 { return self.Version_ }

// part of the item.Item interface
func (self *Node) Visible() bool { return self.Visible_ }

var newNodeNum int64 = 0

func newNodeId() int64 {
	newNodeNum -= 1
	return newNodeNum
}

func New(p *point.Point) *Node {
	return &Node{
		Position_:  p,
		Id_:        newNodeId(),
		Tags_:      tags.New(),
		Timestamp_: time.Now(),
		Version_:   0,
		Changeset_: 0,
		Visible_:   true,
		User_:      user.New(0, ""),
		modified:   true,
		deleted:    false,
	}
}

func (n *Node) SetTags(t *tags.Tags) {
	n.modified = true
	n.Tags_ = t
}

// set the position of node n to the position of p
func (n *Node) MoveTo(p *point.Point) error {
	n.modified = true
	n.Position_ = p
	return nil
}

// relative move, i.e.
//
//   n.RMoveTo(&Point{Lat: 0.1, Lon: 0.0})
//
// would move the node n 0.1 degrees to north
func (n *Node) RMoveTo(p *point.Point) error {
	n.modified = true
	n.Position_.Lat += p.Lat
	n.Position_.Lon += p.Lon
	return nil
}

// marks the node as deleted, note that there is no difference
// in the output as XML currently
func (n *Node) Delete() {
	n.deleted = true
	n.modified = true
}

// equal if both positions are equal
func (n *Node) Equal(o *Node) bool {
	return n.Position().Equal(o.Position())
}

func (nl NodeList) BoundingBox() (bb *bbox.BBox, err error) {
	if nl == nil {
		err = errors.New("Empty Nodelist")
		return
	}
	n := []*Node(nl)
	llat := n[0].Position_.Lat
	ulat := n[0].Position_.Lat
	llon := n[0].Position_.Lon
	ulon := n[0].Position_.Lon

	l := len(n)
	for i := 1; i < l; i++ {
		llat = math.Min(llat, n[i].Position_.Lat)
		ulat = math.Max(ulat, n[i].Position_.Lat)
		llon = math.Min(llon, n[i].Position_.Lon)
		ulon = math.Max(ulon, n[i].Position_.Lon)
	}

	return &bbox.BBox{
		LowerLeft:  point.New(llat, llon),
		UpperRight: point.New(ulat, ulon),
	}, nil
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
