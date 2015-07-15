package bbox

import (
	"fmt"
	"github.com/vetinari/osm/point"
)

type BBox struct {
	LowerLeft  *point.Point
	UpperRight *point.Point
}

func (b *BBox) String() string {
	return fmt.Sprintf("  <bounds minlat='%f' minlon='%f' maxlat='%f' maxlon='%f' origin='%s' />\n",
		b.LowerLeft.Lat, b.LowerLeft.Lon, b.UpperRight.Lon, b.UpperRight.Lat, "")
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
