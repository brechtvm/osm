package relation

import (
	"fmt"
	"osm/node"
	"osm/way"
	"time"
)

func (r *Relation) String() string {
	s := fmt.Sprintf("  <relation id='%d' timestamp='%s' uid='%d' user='%s' visible='%t'"+
		" version='%d' changeset='%d'>\n",
		r.Id_, r.Timestamp_.Format(time.RFC3339), r.User_.Id, r.User_.Name, r.Visible_, r.Version_, r.Changeset_)
	for _, m := range r.GetMembers() {
		var id int64
		switch m.Ref.(type) {
		case *node.Node:
			id = m.Ref.(*node.Node).Id_
		case *way.Way:
			id = m.Ref.(*way.Way).Id_
		case *Relation:
			id = m.Ref.(*Relation).Id_
		}
		s += fmt.Sprintf("    <member type='%s' ref='%d' role='%s' />\n", m.Type(), id, m.Role)
	}
	return s + r.Tags_.String() + "  </relation>\n"
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
