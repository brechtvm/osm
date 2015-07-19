package node

import ()

func (n *Node) MergeOther(o *Node) bool {
	return !(n.modified || n.deleted || n.Version() >= o.Version())
}
