package way

import ()

// sort.Interface:

// sort negative ids on top of all other in ascending math.Abs(Id)
func (wl *WayList) Less(i, j int) bool {
	w := []*Way(*wl)
	if w[i].Id_ < 0 && w[j].Id_ < 0 {
		if w[i].Id_ > w[j].Id_ {
			return true
		}
		return false
	}
	if w[i].Id_ < 0 && w[j].Id_ > 0 {
		return false
	}
	if w[i].Id_ > 0 && w[j].Id_ < 0 {
		return true
	}
	if w[i].Id_ < w[j].Id_ {
		return true
	}
	return false
}

func (wl *WayList) Len() int {
	return len([]*Way(*wl))
}

func (wl *WayList) Swap(i, j int) {
	w := []*Way(*wl)
	w[i], w[j] = w[j], w[i]
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
