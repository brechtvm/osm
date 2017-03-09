package way

import (
	"errors"
	"fmt"
	"github.com/brechtvm/osm/bbox"
	"github.com/brechtvm/osm/item"
	"github.com/brechtvm/osm/node"
	"github.com/brechtvm/osm/point"
	"github.com/brechtvm/osm/tags"
	"github.com/brechtvm/osm/user"
	"time"
)

// part of the item.Item interface
func (self *Way) Type() item.ItemType { return item.TypeWay }

// part of the item.Item interface
func (self *Way) Id() int64           { return self.Id_ }
func (self *Way) Nodes() []*node.Node { return self.Nodes_ }

// part of the item.Item interface
func (self *Way) User() *user.User { return self.User_ }

// part of the item.Item interface
func (self *Way) Tags() *tags.Tags { return self.Tags_ }

// part of the item.Item interface
func (self *Way) Timestamp() time.Time { return self.Timestamp_ }

// part of the item.Item interface
func (self *Way) Version() int64 { return self.Version_ }

// part of the item.Item interface
func (self *Way) Changeset() int64 { return self.Changeset_ }

// part of the item.Item interface
func (self *Way) Visible() bool { return self.Visible_ }

type Way struct {
	Id_        int64
	Nodes_     []*node.Node
	User_      *user.User
	Tags_      *tags.Tags
	Timestamp_ time.Time
	Version_   int64
	Changeset_ int64
	Visible_   bool
	modified   bool
	deleted    bool
}

type WayList []*Way

var newWayNum int64 = 0

func newWayId() int64 {
	newWayNum -= 1
	return newWayNum
}

func EmptyWays() []*Way { return []*Way{} }

// returns a new way, the node slice must have at least two nodes
func New(nl []*node.Node) (w *Way, err error) {
	if len(nl) < 2 {
		return nil, errors.New("Too few nodes for way")
	}
	return &Way{
		Nodes_:     nl,
		Id_:        newWayId(),
		Tags_:      tags.New(),
		Timestamp_: time.Now(),
		Version_:   0,
		Changeset_: 0,
		Visible_:   true,
		User_:      &user.User{Id: 0, Name: ""},
		modified:   true,
		deleted:    false,
	}, nil
}

func (w *Way) SetTags(t *tags.Tags) {
	w.modified = true
	w.Tags_ = t
}

// move all points of w, use the node with id as reference to do a relative move
func (w *Way) MoveTo(id int64, p *point.Point) error {
	ref := w.Nodes_[id]
	if ref == nil {
		return errors.New(fmt.Sprintf("Node #%d not part of Way #%d\n", id, w.Id()))
	}
	return w.RMoveTo(point.New(ref.Position().Lat-p.Lat, ref.Position().Lon-p.Lon))
}

// relative move of all points of w, see Node.RMoveTo()
func (w *Way) RMoveTo(p *point.Point) error {
	for n := range w.Nodes() {
		w.Nodes_[n].RMoveTo(p)
	}
	w.modified = true
	return nil
}

// deletes a way (or more correctly marks as deleted so it will not shown in output).
// Note that the output is currenty unaffected ...
func (w *Way) Delete() {
	w.deleted = true
	w.modified = true
}

// as w.Delete() but also deletes all nodes
func (w *Way) DeleteAll() {
	w.Delete()
	for _, n := range w.Nodes_ {
		n.Delete()
	}
}

func (w *Way) Append(n *node.Node) {
	w.Nodes_ = append(w.Nodes_, n)
	w.modified = true
}

func (w *Way) Prepend(n *node.Node) {
	w.InsertAt(-1, n)
}

func (w *Way) InsertAt(pos int, n *node.Node) {
	w.modified = true
	switch {
	case pos >= len(w.Nodes_):
		w.Append(n)
	case pos == len(w.Nodes_)-1:
		l := w.Nodes_[len(w.Nodes_)-1]
		w.Nodes_[len(w.Nodes_)-1] = n
		w.Nodes_ = append(w.Nodes_, l)
	case pos <= 0:
		nl := []*node.Node{n}
		nl = append(nl, w.Nodes_...)
		w.Nodes_ = nl
	default:
		nl := []*node.Node{}
		for _, n := range w.Nodes_[0:pos] {
			nl = append(nl, n)
		}
		nl = append(nl, n)
		nl = append(nl, w.Nodes_[pos+1:]...)
		w.Nodes_ = nl
	}
}

func (w *Way) InsertAfer(n *node.Node, newNode *node.Node) error {
	pos, err := w.NodePos(n)
	if err != nil {
		return err
	}
	w.InsertAt(pos+1, newNode)
	// w.modified = true
	return nil
}

// returns the index where a given node is in the way
func (w *Way) NodePos(n *node.Node) (int, error) {
	for i, nd := range w.Nodes_ {
		if nd.Id_ == n.Id_ {
			return i, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("Node #%d not a member of way #%d", n.Id(), w.Id()))
}

// Splits a way and returns all resulting ways. For a closed way at least two
// node ids must be given. A resulting way must contain at least two nodes.
func (w *Way) Split(ids ...int64) (ws []*Way, err error) {
	var orig []*node.Node
	for _, n := range w.Nodes_ {
		orig = append(orig, n)
	}

	if w.Closed() {
		if len(ids) < 2 {
			err = errors.New("Not enough arguments to split a closed way")
			return
		}
		var id0 int64
		id0, ids = ids[0], ids[1:]
		_, err = w.OpenAt(id0)
		if err != nil {
			w.Nodes_ = orig
			return
		}
	}

	min_len := 2*len(ids) + 1
	if min_len > len(w.Nodes_) {
		w.Nodes_ = orig
		err = errors.New("Way has too few nodes to split")
		return
	}

	ws = []*Way{w}
	cur_w := w
	for _, id := range ids {
		// the node with the id "id" will be the last node of the current
		// way, a new node with the same position is inserted at the
		// beginning of the new way
		var nid int64
		nid, err = cur_w.NextNodeId(id)
		if err != nil {
			w.Nodes_ = orig
			ws = EmptyWays()
			err = errors.New(fmt.Sprintf("Node #%d not a member of way #%d", nid, w.Id()))
			return
		}
		for i, n := range cur_w.Nodes_ {
			if n.Id_ == nid {
				if i == 0 || i == len(cur_w.Nodes_)-1 {
					w.Nodes_ = orig
					ws = EmptyWays()
					err = errors.New("Cannot split way at first or last node")
					return
				}
				if i == 1 {
					w.Nodes_ = orig
					ws = EmptyWays()
					err = errors.New("Cannot split into a way with one node")
					return
				}

				nn := node.New(n.Position_)
				nn.Tags_ = n.Tags_
				nl := []*node.Node{nn}
				for _, nd := range cur_w.Nodes_[i:] {
					nl = append(nl, nd)
				}
				cur_w.Nodes_ = cur_w.Nodes_[0:i]
				var new_w *Way
				new_w, err = New(nl)
				if err != nil {
					w.Nodes_ = orig
					ws = EmptyWays()
					return
				}
				new_w.Tags_ = cur_w.Tags_
				ws = append(ws, new_w)
				cur_w = new_w
			}
		}
	}
	w.modified = true
	return ws, nil
}

// Opens a closed way (i.e. a way where the last node has the
// same id as the first one) at the node with the given id. A
// new node is inserted where the way is open. This new node
// will be the end of the current way -> the first node will
// not necessarily stay the first node
func (w *Way) OpenAt(id int64) (nn *node.Node, err error) {
	if !w.Closed() {
		return nil, errors.New("Cannot open non-closed way")
	}
	var i int
	var nd *node.Node
	last := len(w.Nodes_) - 1
	for i, nd = range w.Nodes_ {
		if nd.Id_ != id {
			continue
		}

		w.modified = true
		nn = node.New(nd.Position_)
		nn.Tags_ = nd.Tags_

		if i == 0 {
			w.Nodes_[last] = nn
			return nn, nil
		}

		nb, nl := w.Nodes_[:i], w.Nodes_[i:]
		nl = append(nl, nb...)
		w.Nodes_ = append(nl, nn)

		return nn, nil
	}

	return nil, errors.New(fmt.Sprintf("Node #%d not a member of way #%d", id, w.Id()))
}

// Closes the way. If the first and the last node have the same position
// but not the same id, the last node is deleted. The first node is
// appended to the way.
func (w *Way) Close() error {
	if w.Closed() {
		return nil
	}
	last := len(w.Nodes_) - 1
	if last < 2 {
		return errors.New("Too few nodes")
	}

	// w.Closed() already checked that they're not the same id:
	if w.Nodes_[0].Position_.Equal(w.Nodes_[last].Position_) {
		w.Nodes_[last].Delete()
		w.Nodes_[last] = w.Nodes_[0]
	} else {
		w.Nodes_ = append(w.Nodes_, w.Nodes_[0])
	}
	return nil
}

// Returns the id of the node before the node with the given id.
func (w *Way) PrevNodeId(id int64) (int64, error) {
	if w.Nodes_[0].Id_ == id {
		return 0, errors.New("No previous node")
	}
	for i, nd := range w.Nodes_ {
		if nd.Id_ == id {
			return w.Nodes_[i-1].Id_, nil
		}
	}
	return 0, errors.New("Node not member of way")
}

// Returns the id of the node following the node with the given id.
func (w *Way) NextNodeId(id int64) (int64, error) {
	if w.Nodes_[len(w.Nodes_)-1].Id_ == id {
		return 0, errors.New("No next node")
	}
	for i := 0; i < len(w.Nodes_)-1; i++ {
		if w.Nodes_[i].Id_ == id {
			return w.Nodes_[i+1].Id_, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("Node %d not member of way %d", id, w.Id_))
}

func (w *Way) GetNodes() []*node.Node {
	return w.Nodes_
}

// Returns true if the first and the last node of the way
// have the same id (and the way must be longer than two
// nodes)
func (w *Way) Closed() bool {
	num := len(w.Nodes_)
	if num < 3 {
		return false
	}
	return w.Nodes_[0].Id_ == w.Nodes_[num-1].Id_
}

func (w *Way) Clone() (nw *Way, err error) {
	nd := []*node.Node{}
	for _, n := range w.Nodes_ {
		nd = append(nd, n)
	}
	t := tags.New()
	for key := range map[string]string(*w.Tags_) {
		t.Add(key, w.Tags_.Get(key))
	}
	nw, err = New(nd)
	if err != nil {
		nw.Tags_ = t
	}
	return
}

func (w *Way) ConnectedTo(y *Way) bool {
	if w.Nodes_[len(w.Nodes_)-1].Id_ == y.Nodes_[0].Id_ {
		return true
	}
	if w.Nodes_[len(w.Nodes_)-1].Position_.Equal(y.Nodes_[0].Position_) {
		return true
	}
	return false
}

func (w *Way) Connect(ways ...*Way) {
	for _, y := range ways {
		w.Nodes_ = append(w.Nodes_, y.Nodes_...)
		y.Delete()
		// FIXME - merge tags
	}
	w.Cleanup()
	w.modified = true
}

func (w *Way) FirstNode() *node.Node {
	return w.Nodes_[0]
}

func (w *Way) LastNode() *node.Node {
	return w.Nodes_[len(w.Nodes_)-1]
}

type Connection int

const (
	NotConnected          Connection = 0
	ConnectedNormal       Connection = 1
	ConnectedReversed2nd  Connection = 2
	ConnectedReversed1st  Connection = 3
	ConnectedReversedBoth Connection = 4
)

func (c Connection) String() string {
	switch c {
	case NotConnected:
		return "NotConnected"
	case ConnectedNormal:
		return "ConnectedNormal"
	case ConnectedReversed2nd:
		return "ConnectedReversed2nd"
	case ConnectedReversed1st:
		return "ConnectedReversed1st"
	case ConnectedReversedBoth:
		return "ConnectedReversedBoth"
	}
	return "Unknown"
}

func (w *Way) Connected(y *Way) Connection {
	if w.LastNode().Id_ == y.FirstNode().Id_ ||
		w.LastNode().Position_.Equal(y.FirstNode().Position_) {
		return ConnectedNormal
	}
	if w.LastNode().Id_ == y.LastNode().Id_ ||
		w.LastNode().Position_.Equal(y.LastNode().Position_) {
		return ConnectedReversed2nd
	}
	if w.FirstNode().Id_ == y.FirstNode().Id_ ||
		w.FirstNode().Position_.Equal(y.FirstNode().Position_) {
		return ConnectedReversed1st
	}
	if w.FirstNode().Id_ == y.LastNode().Id_ ||
		w.FirstNode().Position_.Equal(y.LastNode().Position_) {
		return ConnectedReversedBoth
	}
	return NotConnected
}

func (w *Way) Join(ways ...*Way) error {
	for _, wy := range ways {
		if !w.ConnectedTo(wy) {
			return errors.New("Cannot join non connected ways")
		}
		w.Connect(wy)
	}
	return nil
}

func (w *Way) Cleanup() *Way {
	nd := []*node.Node{w.Nodes_[0]}
	for i := 1; i < len(w.Nodes_); i++ {
		if w.Nodes_[i].Id_ != w.Nodes_[i-1].Id_ {
			nd = append(nd, w.Nodes_[i])
		}
	}
	w.Nodes_ = nd
	return w
}

func (w *Way) Equal(x *Way) bool {
	if len(w.Nodes_) != len(x.Nodes_) {
		return false
	}
	for i, n := range w.Nodes_ {
		if !n.Equal(x.Nodes_[i]) {
			return false
		}
	}
	return true
}

func (w *Way) IsClockwise() bool {
	if !w.Closed() {
		return false
	}
	return w.area() < 0
}

func (w *Way) Clockwise() {
	if !w.Closed() {
		return
	}
	if !w.IsClockwise() {
		w.Reverse()
	}
}

func (w *Way) Reverse() {
	var n []*node.Node
	for i := len(w.Nodes_) - 1; i >= 0; i-- {
		n = append(n, w.Nodes_[i])
	}
	w.Nodes_ = n
	// FIXME - reverse also the tag meanings!
}

func (w *Way) IsCounterClockwise() bool {
	if !w.Closed() {
		return false
	}
	return w.area() > 0
}

func (w *Way) BoundingBox() (*bbox.BBox, error) {
	return node.NodeList(w.Nodes_).BoundingBox()
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
