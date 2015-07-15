package way

import (
	"fmt"
	"github.com/vetinari/osm/point"
)

func (w *Way) GPX(center *point.Point) string {
	gpx := gpxHeader()
	for _, n := range w.Nodes() {
		gpx += fmt.Sprintf("   <trkpt lat='%.7f' lon='%.7f' />\n", n.Position().Lat, n.Position().Lon)
	}
	gpx += "  </trkseg>\n </trk>\n"
	if center != nil {
		gpx += fmt.Sprintf(" <wpt lat='%.7f' lon='%.7f'>\n", center.Lat, center.Lon)
		gpx += fmt.Sprintf("  <name>%s</name>\n", center.Name)
		gpx += " </wpt>\n"
	}
	gpx += "</gpx>\n"
	return gpx
}

func gpxHeader() string {
	hdr := `<?xml version='1.0' encoding='UTF-8'?>
<gpx version="1.1" generator="osm.gpx.go v%s"
    xmlns="http://www.topografix.com/GPX/1/1">.
 <metadata>
  <name></name>
  <desc></desc>
  <author>
   <name></name>
   <email domain="" id=""/>
  </author>
  <copyright author="">
  <year></year>
  <license></license>
  </copyright>
 </metadata>
 <trk>
  <trkseg>
`
	return fmt.Sprintf(hdr, "0.1")
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
