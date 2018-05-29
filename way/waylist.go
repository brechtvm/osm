package way

import (
	"errors"
	"fmt"
)

func (wl WayList) Connected(w *Way) Connection {
	wlw, err := wl.AsWay()
	if err != nil {
		return NotConnected
	}
	return wlw.Connected(w)
}

func (wl WayList) AsWay() (w *Way, err error) {
	lst := []*Way(wl)
	if len(lst) == 0 {
		return nil, errors.New("Empty WayList")
	} else if len(lst) == 1 {
		return lst[0], nil
	}
	w, err = New(lst[0].Nodes_)
	if err != nil {
		return
	}
	wl = lst[1:]
	var n *Way
	for _, cur := range wl {
		switch w.Connected(cur) {
		case ConnectedNormal:
			n, err = New(cur.Nodes_)
			if err != nil {
				return
			}
			w.Join(n)
		case ConnectedReversed2nd:
			n, err = New(cur.Nodes_)
			if err != nil {
				return
			}
			n.Reverse()
			w.Join(n)
		case ConnectedReversed1st:
			w.Reverse()
			n, err = New(cur.Nodes_)
			if err != nil {
				return
			}
			w.Join(n)
		case ConnectedReversedBoth:
			w.Reverse()
			n, err = New(cur.Nodes_)
			if err != nil {
				return
			}
			n.Reverse()
			w.Join(n)
		case NotConnected:
			err = errors.New(fmt.Sprintf("Not connected to way #%d\n", cur.Id))
			return
		}
	}
	return w, nil
}

func (wl WayList) Add(w *Way) {
	wl = WayList(append([]*Way(wl), w))
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go
