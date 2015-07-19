package overpass

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/vetinari/osm"
	"github.com/vetinari/osm/bbox"
	"github.com/vetinari/osm/xml"
	"io/ioutil"
	"net/http"
	"os"
)

type OverpassAPI struct {
	URL     string
	Timeout int
	Data    []byte // response data
	Type    osm.DataFormat
}

func New() *OverpassAPI {
	return &OverpassAPI{
		URL:     "http://overpass-api.de/api/interpreter",
		Timeout: 180,
	}
}

func (o *OverpassAPI) Query(query string) (err error) {
	// FIXME - validate query (<print mode="meta" />)
	// FIXME - add timeout to query unless set
	var req *http.Request
	var res *http.Response

	body := bytes.NewBufferString(query)

	req, err = http.NewRequest("POST", o.URL, body)
	if err != nil {
		return
	}

	// give query timeout + some extra seconds for the clients
	// client := http.Client{Timeout: (o.Timeout + 10) * time.Second}
	client := http.Client{}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err = client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	o.Data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	switch res.StatusCode {
	case http.StatusOK:
		switch res.Header.Get("Content-Type") {
		case "application/osm3s+xml":
			o.Type = osm.FmtXML
		case "application/json":
			o.Type = osm.FmtOverpassJSON
		default:
			o.Type = osm.FmtUnknown
		}
	case http.StatusBadRequest:
		err = errors.New(fmt.Sprintf("Syntax error in request: %s", parseError(o.Data)))
	case 429: // FIXME when net/http exports StatusTooManyRequests
		err = errors.New("Too many requests")
	case http.StatusGatewayTimeout:
		err = errors.New("Server load too high")
	default:
		fmt.Fprintf(os.Stderr, "OverpassAPI Response=%s\n", o.Data)
		err = errors.New(fmt.Sprintf("Server returned unknown error %d", res.StatusCode))
	}
	return
}

func parseError(data []byte) []byte {
	off := bytes.Index(data, []byte(">Error</strong>: "))
	if off == -1 {
		return []byte("error message not found")
	}
	data = data[off+17:]
	off = bytes.Index(data, []byte("</p>"))
	if off == -1 {
		return data
	}
	return data[:off]
}

func BBoxParser(bb *bbox.BBox) (osm.Parser, error) {
	query := fmt.Sprintf(`<osm-script output="xml">
  <union into="_">
    <bbox-query s="%.7f" w="%.7f" n="%.7f" e="%.7f" />
      <recurse type="up"/>
      <recurse type="down"/>
  </union>
  <print limit="" mode="meta" order="id"/>
</osm-script>`, bb.LowerLeft.Lat,
		bb.LowerLeft.Lon,
		bb.UpperRight.Lat,
		bb.UpperRight.Lon)
	return QueryParser(query)
}

func QueryParser(query string) (osm.Parser, error) {
	o := New()
	err := o.Query(query)
	if err != nil {
		return nil, err
	}
	switch o.Type {
	case osm.FmtXML:
		return xml.ByteParser(o.Data), nil
	case osm.FmtOverpassJSON:
		return nil, errors.New("Cannot parse JSON yet")
	default:
		return nil, errors.New("Unknown Content-Type")
	}
}

// vim: ts=4 sw=4 noexpandtab nolist
