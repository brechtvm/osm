package relation

import ()

// sort.Interface:

// sort negative ids on top of all other in ascending math.Abs(Id)
func (rl *RelationList) Less(i, j int) bool {
	r := []*Relation(*rl)
	if r[i].Id_ < 0 && r[j].Id_ < 0 {
		if r[i].Id_ > r[j].Id_ {
			return true
		}
		return false
	}
	if r[i].Id_ < 0 && r[j].Id_ > 0 {
		return false
	}
	if r[i].Id_ > 0 && r[j].Id_ < 0 {
		return true
	}
	if r[i].Id_ < r[j].Id_ {
		return true
	}
	return false
}

func (rl *RelationList) Len() int {
	return len([]*Relation(*rl))
}

func (rl *RelationList) Swap(i, j int) {
	r := []*Relation(*rl)
	r[i], r[j] = r[j], r[i]
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
