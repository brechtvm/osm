package item

import (
	"github.com/vetinari/osm/tags"
	"github.com/vetinari/osm/user"
	"time"
)

type ItemType int

const (
	TypeUnknown ItemType = 1 * iota
	TypeNode
	TypeWay
	TypeRelation
)

func (t ItemType) String() string {
	switch t {
	case TypeNode:
		return "node"
	case TypeWay:
		return "way"
	case TypeRelation:
		return "relation"
	default:
		panic("unknown item type")
	}
}

func ItemTypeFromString(s string) ItemType {
	switch s {
	case "node":
		return TypeNode
	case "way":
		return TypeWay
	case "relation":
		return TypeRelation
	default:
		return TypeUnknown
	}
}

// the Item interface is implemented by all real OSM objects (node.Node,
// way.Way, relation.Relation).
type Item interface {
	Id() int64
	Type() ItemType
	User() *user.User
	Tags() *tags.Tags
	Timestamp() time.Time
	Version() int64
	Changeset() int64
	Visible() bool
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
