package osmorg

import (
	"errors"
	"fmt"
	"github.com/brechtvm/osm"
	"github.com/brechtvm/osm/bbox"
	"github.com/brechtvm/osm/xml"
	"io/ioutil"
	"net/http"
	"os"
)

type OSMAPI struct {
	URL     string
	Timeout int
	Data    []byte // response data
	Type    osm.DataFormat
}

func New() *OSMAPI {
	return &OSMAPI{
		URL:     "http://www.openstreetmap.org/api/0.6",
		Timeout: 300,
	}
}

func (o *OSMAPI) Query(query string) (err error) {
	var req *http.Request
	var res *http.Response

	req, err = http.NewRequest("GET", o.URL+query, nil)
	if err != nil {
		return
	}

	// give query timeout + some extra seconds for the clients
	// client := http.Client{Timeout: (o.Timeout + 10) * time.Second}
	client := http.Client{}

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
		o.Type = osm.FmtXML
	case http.StatusBadRequest:
		err = errors.New("Node/Way/Relation limit hit")
	case 509: // Bandwidth Limit Exceeded
		err = errors.New("Error: You have downloaded too much data. Please try again later.")
	default:
		fmt.Fprintf(os.Stderr, "OSM API Response=%s\n", o.Data)
		err = errors.New(fmt.Sprintf("Server returned unknown error %d", res.StatusCode))
	}
	return
}

func BBoxParser(bb *bbox.BBox) (osm.Parser, error) {
	query := fmt.Sprintf("/map?bbox=%.7f,%.7f,%.7f,%.7f",
		bb.LowerLeft.Lon,
		bb.LowerLeft.Lat,
		bb.UpperRight.Lon,
		bb.UpperRight.Lat,
	)
	return QueryParser(query)
}

func QueryParser(query string) (osm.Parser, error) {
	o := New()
	err := o.Query(query)
	if err != nil {
		return nil, err
	}

	if o.Type == osm.FmtXML {
		return xml.ByteParser(o.Data), nil
	}

	return nil, errors.New("Unknown Content-Type")
}

// vim: ts=4 sw=4 noexpandtab nolist
