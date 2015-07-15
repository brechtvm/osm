package tags

import (
	//	"fmt"
	"strings"
)

type Tags map[string]string

// allocates a new *Tags
func New() *Tags {
	tm := Tags(make(map[string]string))
	return &tm
}

// deletes the named key from the *Tags
func (t *Tags) Delete(k string) {
	if t == nil {
		return
	}
	delete(map[string]string(*t), k)
}

// adds the named key + value to the *Tags (replacing an eventually existing value)
func (t *Tags) Add(k string, v string) {
	if t == nil {
		t = New()
	}
	m := map[string]string(*t)
	m[k] = v
}

// checks if the *Tags have the given key
func (t *Tags) Has(key string) bool {
	m := map[string]string(*t)
	_, ok := m[key]
	return ok
}

// returns the given key from the *Tags - NOTE: a non
// existing tag returns an empty string, so use
//
//   var val string
//   if t.Has(key) {
//      val = t.Get(key)
//   }
func (t *Tags) Get(key string) string {
	m := map[string]string(*t)
	return m[key]
}

// see func (tags *Tags) Reverse()
var OppositeValues = map[string]map[string]string{
	// tag k=...
	"oneway": map[string]string{
		// v => reversed v
		"yes": "-1",
		"-1":  "yes",
	},
	"incline": map[string]string{
		"up":   "down",
		"down": "up",
	},
}

// see func (tags *Tags) Reverse()
var OppositeKeys = map[string]string{
	// tag k=...
	"destination:forward":       "destination:backward",
	"destination:lanes:forward": "destination:lanes:backward",
}

// Reverses the meaning of tags, i.e. changes oneway=yes -> oneway=-1
// or sometag:forward -> sometag:backward (e.g. destination:forward
// -> destination:backward)
// NOTE - the OppositeKeys and OppositeValues is probably far from
// complete ...
func (tags *Tags) Reverse() {
	t := map[string]string(*tags)
	ntags := New()
	n := map[string]string(*ntags)
	for k, v := range t {
		if strings.HasSuffix(k, ":forward") {
			k = strings.Replace(k, ":forward", ":backward", 1)
			n[k] = v
			continue
		} else if strings.HasSuffix(k, ":backward") {
			k = strings.Replace(k, ":backward", ":forward", 1)
			n[k] = v
			continue
		}

		ok, exists := OppositeKeys[k]
		if exists {
			n[ok] = v
			continue
		}
		o, exists := OppositeValues[k]
		if exists {
			for ok, ov := range o {
				if ok == v {
					n[k] = ov
					break
				}
			}
		} else {
			n[k] = v
		}
	}
	tags = ntags
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
