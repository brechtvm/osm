package xml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"osm"
	"osm/item"
	"osm/node"
	"osm/point"
	"osm/relation"
	"osm/tags"
	"osm/user"
	"osm/way"
	"strings"
	"time"
)

type DotOSM struct {
	nextLine func() (string, error)
	data     *osm.OSM
}

// returns an osm.Parser which can be used as argument to osm.New()
func Parser(r io.Reader) osm.Parser {
	return &DotOSM{nextLine: newLineReader(r)}
}

// returns an osm.Parser which can be used as argument to osm.New(), reads
// from byte array
func ByteParser(data []byte) osm.Parser {
	return Parser(bytes.NewReader(data))
}

// returns an osm.Parser which can be used as argument to osm.New(), reads
// from the given file
func FileParser(file string) (osm.Parser, io.Closer, error) {
	fh, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	return Parser(fh), fh, nil
}

// implements the osm.Parser interface
func (p *DotOSM) Parse() (o *osm.OSM, err error) {
	// FIXME - parse "<bbox "...
	o = osm.NewOSM()
	p.data = o
	var i item.Item
	var line string
	var isEOF bool
	for {
		line, err = p.nextLine()
		if err == io.EOF {
			isEOF = true
		} else if err != nil {
			return
		}
		if line != "" {
			i, err = p.parseItem(line)
			if err != nil {
				break
			}
			switch i.Type() {
			case item.TypeNode:
				o.Nodes[i.(*node.Node).Id_] = i.(*node.Node)
			case item.TypeWay:
				o.Ways[i.(*way.Way).Id_] = i.(*way.Way)
			case item.TypeRelation:
				o.Relations[i.(*relation.Relation).Id_] = i.(*relation.Relation)
			default:
			}
		}
		if isEOF {
			break
		}
	}
	return
}

func xmlClosedItem(s string, t string) bool {
	if strings.HasSuffix(s, "/>") {
		return true
	}
	if s == "</"+t+">" {
		return true
	}
	return false
}

type xmlItem struct {
	id        int64
	t         item.ItemType
	user      *user.User
	tags      *tags.Tags
	ts        time.Time
	version   int64
	changeset int64
	visible   bool
	data      interface{}
}

func (x *xmlItem) Id() int64            { return x.id }
func (x *xmlItem) Type() item.ItemType  { return x.t }
func (x *xmlItem) User() *user.User     { return x.user }
func (x *xmlItem) Tags() *tags.Tags     { return x.tags }
func (x *xmlItem) Version() int64       { return x.version }
func (x *xmlItem) Changeset() int64     { return x.changeset }
func (x *xmlItem) Visible() bool        { return x.visible }
func (x *xmlItem) Data() interface{}    { return x.data }
func (x *xmlItem) Timestamp() time.Time { return x.ts }

func (x *DotOSM) parseItem(line string) (i item.Item, err error) {
	prefix, closed, m := x.parseLine(line)

	var ts time.Time
	if m["timestamp"] != "" {
		ts, err = time.Parse(time.RFC3339, m["timestamp"])
		if err != nil {
			err = errors.New(fmt.Sprintf("Failed to parse timestamp '%s': %s\n", m["timestamp"], err))
			return
		}
	}
	switch prefix {
	case "<node":
		n := &node.Node{
			Id_:        str2int64(m["id"]),
			User_:      user.New(str2int64(m["uid"]), m["user"]),
			Position_:  point.New(str2float64(m["lat"]), str2float64(m["lon"])),
			Timestamp_: ts,
			Version_:   str2int64(m["version"]),
			Changeset_: str2int64(m["changeset"]),
			Visible_:   str2bool(m["visible"]),
		}
		if !closed {
			n.Tags_ = x.parseTags("node")
		}
		// fmt.Printf("NODE=%q\n", n)
		return n, nil

	case "<way":
		if closed {
			err = errors.New(fmt.Sprintf("Way %s has no nodes\n", m["id"]))
			return
		}

		w := &way.Way{
			Id_:        str2int64(m["id"]),
			User_:      user.New(str2int64(m["uid"]), m["user"]),
			Timestamp_: ts,
			Version_:   str2int64(m["version"]),
			Changeset_: str2int64(m["changeset"]),
			Visible_:   str2bool(m["visible"]),
		}
		var t *tags.Tags
		var nd []*node.Node
		t, nd, err = x.parseWay(str2int64(m["id"]))
		if err != nil {
			return
		}
		w.Tags_ = t
		w.Nodes_ = nd

		// fmt.Printf("WAY=%q\n", w)
		return w, nil

	case "<relation":
		if closed {
			err = errors.New(fmt.Sprintf("Relation %s has no members\n", m["id"]))
			return
		}

		rel := &relation.Relation{
			Id_:        str2int64(m["id"]),
			User_:      user.New(str2int64(m["uid"]), m["user"]),
			Timestamp_: ts,
			Version_:   str2int64(m["version"]),
			Changeset_: str2int64(m["changeset"]),
			Visible_:   str2bool(m["visible"]),
		}
		members, tags := x.parseRelation(str2int64(m["id"]))
		rel.Members_ = members
		rel.Tags_ = tags

		// fmt.Printf("RELATION=%q\n", rel)
		return rel, nil

	case "<osm":
		// fmt.Printf("OSM=%q\n", m)
	case "</osm>":
		// fmt.Printf("/OSM\n")
	case "<bounds":
		// fmt.Printf("BOUNDS=%q\n", m)
	default:
		// fmt.Printf("OTHER=%q\n", m)
	}
	return &xmlItem{t: item.TypeUnknown}, nil
}

func (x *DotOSM) parseWay(id int64) (t *tags.Tags, n []*node.Node, err error) {
	t = tags.New()
	line, _ := x.nextLine()
	for line != "</way>" {
		prefix, _, m := x.parseLine(line)
		switch prefix {
		case "<nd":
			ref := str2int64(m["ref"])
			nd := x.data.GetNode(ref)
			if nd == nil {
				err = errors.New(fmt.Sprintf("missing node %d in way %d\n", ref, id))
				return
			}
			n = append(n, nd)
		case "<tag":
			t.Add(m["k"], decodeXML(m["v"]))
		default:
			panic("unknown element " + prefix)
		}
		line, _ = x.nextLine()
	}

	return
}

func (x *DotOSM) parseRelation(id int64) (members []*relation.Member, t *tags.Tags) {
	t = tags.New()
	line, _ := x.nextLine()
	for line != "</relation>" {
		prefix, _, m := x.parseLine(line)
		switch prefix {
		case "<member":
			ref := str2int64(m["ref"])

			member := &relation.Member{Type_: item.ItemTypeFromString(m["type"]), Role: m["role"]}
			switch member.Type() {
			case item.TypeNode:
				member.Ref = x.data.GetNode(ref)
			case item.TypeWay:
				member.Ref = x.data.GetWay(ref)
			case item.TypeRelation:
				member.Ref = x.data.GetRelation(ref)
			}
			if member.Ref == nil {
				fmt.Printf("WARNING: Missing %s id %d in relation %d\n", member.Type(), ref, id)
			}
			members = append(members, member)

		case "<tag":
			t.Add(m["k"], decodeXML(m["v"]))
		default:
			panic("unknown element " + prefix)
		}
		line, _ = x.nextLine()
	}
	return members, t
}

func (x *DotOSM) parseLine(line string) (str string, closed bool, m map[string]string) {
	s := strings.SplitN(line, " ", 2)
	str = s[0]
	if len(s) == 2 {
		m = line2map(s[1])
		closed = line[len(line)-2] == '/' // && line[len(line)-1] == '>'
	} else {
		closed = line[0] == '<' && line[1] == '/' && line[len(line)-1] == '>'
	}
	return
}

func (x *DotOSM) parseTags(s string) *tags.Tags {
	line, _ := x.nextLine()
	eot := "</" + s + ">"
	t := tags.New()
	for line != eot {
		prefix, _, m := x.parseLine(line)
		if prefix != "<tag" {
			panic("Non Tag element found: " + prefix)
		}
		t.Add(m["k"], decodeXML(m["v"]))
		line, _ = x.nextLine()
	}
	return t
}

var enc_entities = map[string][]byte{
	"amp":  []byte{'&'},
	"lt":   []byte{'<'},
	"gt":   []byte{'>'},
	"quot": []byte{'"'},
	"apos": []byte{'\''},
	"#xD":  []byte{'\n'},
	"#xA":  []byte{'\r'},
}

var dec_entities = map[byte]string{
	'&':  "&amp;",
	'"':  "&quot;",
	'\'': "&apos;",
	'<':  "&lt;",
	'>':  "&gt;",
	'\n': "&#xA;",
	'\r': "&#xD;",
}

func encodeXML(v string) string {
	s := []byte(v)
	var o []byte
	for i := 0; i < len(s); i++ {
		c, ok := dec_entities[s[i]]
		if ok {
			o = append(o, []byte(c)...)
		} else {
			o = append(o, s[i])
		}
	}
	return string(o)
}

func decodeXML(v string) string {
	s := []byte(v)
	var o []byte
	for i := 0; i < len(s); i++ {
		if s[i] == '&' {
			for j := i; j < len(s); j++ {
				if s[j] == ';' {
					c := s[i+1 : j]
					o = append(o, enc_entities[string(c)]...)
					i = j
					break
				}
			}
		} else {
			o = append(o, s[i])
		}
	}
	return string(o)
}

func newLineReader(fh io.Reader) func() (string, error) {
	var buf []byte
	fn := func() (string, error) {
		var line []byte
		var err error
		var c int
		nl := bytes.Index(buf, []byte{'\n'})
		for nl == -1 {
			var rbuf = make([]byte, 32768)
			c, err = fh.Read(rbuf)
			if c == 0 {
				break
			}
			buf = append(buf, rbuf...)
			nl = bytes.Index(buf, []byte{'\n'})
		}
		if nl != -1 {
			line = buf[0:nl]
			buf = buf[nl+1:]
		} else {
			line = buf
		}
		line = bytes.Trim(line, " \t")
		return string(line), err
	}
	return fn
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
