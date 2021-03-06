package bbox

import (
	"fmt"
	"github.com/brechtvm/osm/point"
)

type BBox struct {
	LowerLeft  *point.Point
	UpperRight *point.Point
}

func (b *BBox) String() string {
	return fmt.Sprintf("  <bounds minlat='%f' minlon='%f' maxlat='%f' maxlon='%f'  />\n",
		b.LowerLeft.Lat, b.LowerLeft.Lon, b.UpperRight.Lat, b.UpperRight.Lon)
}

func (b *BBox) Contains(p *point.Point) bool {
	return p.Lat > b.LowerLeft.Lat && p.Lat < b.UpperRight.Lat &&
		p.Lon > b.LowerLeft.Lon && p.Lon < b.UpperRight.Lon
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
