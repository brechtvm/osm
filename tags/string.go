package tags

import (
	"fmt"
	"sort"
)

func (t *Tags) String() string {
	s := ""
	if t == nil {
		return s
	}
	keys := make([]string, 0)
	for key, _ := range *t {
		//if key == "ref" || key == "name" || key == "type" {
		//	s += fmt.Sprintf(`    <tag k="%s" v="%s" />`+"\n", key, encodeXML((*t)[key]))
		//} else {
			keys = append(keys, key)
		//}
	}
	sort.Strings(keys)

	for _, key := range keys {
		s += fmt.Sprintf(`    <tag k="%s" v="%s" />`+"\n", key, encodeXML((*t)[key]))
	}
	return s
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

// vim: ts=4 sw=4 noexpandtab nolist syn=go
