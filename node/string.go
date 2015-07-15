package node

import (
	"fmt"
	"time"
)

// OSM XML output of a node
func (n *Node) String() string {
	s := fmt.Sprintf("  <node id='%d' timestamp='%s' uid='%d' user='%s' visible='%t'"+
		" version='%d' changeset='%d' lat='%f' lon='%f'",
		n.Id_, n.Timestamp_.Format(time.RFC3339), n.User_.Id, n.User_.Name, n.Visible_,
		n.Version_, n.Changeset_, n.Position_.Lat, n.Position_.Lon)
	t := n.Tags_.String()
	if t == "" {
		return s + " />\n"
	}
	return s + ">\n" + t + "  </node>\n"
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
