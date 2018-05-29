package node

import ()

// sort negative ids on top of all other in ascending math.Abs(Id)
func (nl *NodeList) Less(i, j int) bool {
	n := []*Node(*nl)
	if n[i].Id_ < 0 && n[j].Id_ < 0 {
		if n[i].Id_ > n[j].Id_ {
			return true
		}
		return false
	}
	if n[i].Id_ < 0 && n[j].Id_ > 0 {
		return false
	}
	if n[i].Id_ > 0 && n[j].Id_ < 0 {
		return true
	}
	if n[i].Id_ < n[j].Id_ {
		return true
	}
	return false
}

func (nl *NodeList) Len() int {
	return len([]*Node(*nl))
}

func (nl *NodeList) Swap(i, j int) {
	n := []*Node(*nl)
	n[i], n[j] = n[j], n[i]
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
