package relation

import ()

func (r *Relation) MergeOther(o *Relation) bool {
	return !(r.modified || r.deleted || r.Version() >= o.Version())
}
