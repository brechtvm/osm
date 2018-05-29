package way

import (
	"fmt"
	"golang.org/x/net/html"
	"time"
)

func (w *Way) String() string {
	s := fmt.Sprintf(`  <way id="%d" version="%d" timestamp="%s" changeset="%d" uid="%d" user="%s">`+"\n",
		w.Id_, w.Version_, w.Timestamp_.Format(time.RFC3339), w.Changeset_, w.User_.Id, html.EscapeString(w.User_.Name))
	for _, n := range w.NodeIDs {
		s += fmt.Sprintf(`    <nd ref="%d" />`+"\n", n)
	}
	return s + w.Tags_.String() + "  </way>\n"
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
