package way

import (
	"fmt"
	"time"
)

func (w *Way) String() string {
	s := fmt.Sprintf(`  <way id="%d" timestamp="%s" uid="%d" user="%s" visible="%t"  version="%d" changeset="%d">`+"\n",
		w.Id_, w.Timestamp_.Format(time.RFC3339), w.User_.Id, w.User_.Name, w.Visible_, w.Version_, w.Changeset_)
	for _, n := range w.GetNodes() {
		s += fmt.Sprintf(`    <nd ref="%d" />`+"\n", n.Id())
	}
	return s + w.Tags_.String() + "  </way>\n"
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
