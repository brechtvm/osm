package node

import (
	"fmt"
	"golang.org/x/net/html"
	"time"
)

// OSM XML output of a node
func (n *Node) String() string {
	// Do not use %f but use %v!
	s := fmt.Sprintf(`  <node id="%d" lat="%v" lon="%v" version="%d" timestamp="%s" changeset="%d" uid="%d" user="%s" `,
		n.Id_,
		n.Position_.Lat,
		n.Position_.Lon,
		n.Version_,
		n.Timestamp_.Format(time.RFC3339),
		n.Changeset_,
		n.User_.Id,
		html.EscapeString(n.User_.Name))
	t := n.Tags_.String()
	if t == "" {
		return s + " />\n"
	}
	return s + ">\n" + t + "  </node>\n"
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
