package relation

import (
	"fmt"
	"github.com/brechtvm/osm/node"
	"github.com/brechtvm/osm/way"
	"golang.org/x/net/html"
	"log"
	"time"
)

func (r *Relation) String() string {
	s := fmt.Sprintf(`  <relation id="%d" version="%d" timestamp="%s" changeset="%d" uid="%d" user="%s">`+"\n",
		r.Id_, r.Version_, r.Timestamp_.Format(time.RFC3339), r.Changeset_, r.User_.Id, html.EscapeString(r.User_.Name))
	for _, m := range r.GetMembers() {
		var id int64
		if m.Ref == nil {
			log.Printf("ERROR: %s member ref #%d is nil", m.Type(), m.Id())
			id = 0
		} else {
			switch m.Ref.(type) {
			case *node.Node:
				if m.Ref.(*node.Node) == nil {
					//log.Printf("Missing node #%d in relation #%d\n", m.Id(), r.Id())
					id = m.Id()
				} else {
					id = m.Ref.(*node.Node).Id_
				}
			case *way.Way:
				if m.Ref.(*way.Way) == nil {
					//log.Printf("Missing way #%d in relation #%d\n", m.Id(), r.Id())
					id = m.Id()
				} else {
					id = m.Ref.(*way.Way).Id_
				}
			case *Relation:
				if m.Ref.(*Relation) == nil {
					//log.Printf("Missing relation #%d in relation #%d\n", m.Id(), r.Id())
					id = m.Id()
				} else {
					id = m.Ref.(*Relation).Id_
				}
			}
		}
		s += fmt.Sprintf(`    <member type="%s" ref="%d" role="%s" />`+"\n", m.Type(), id, m.Role)
	}
	return s + r.Tags_.String() + "  </relation>\n"
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
