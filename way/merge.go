package way

import (

)

func (w *Way) MergeOther(o *Way) bool {
    return !(w.modified || w.deleted || w.Version() >= o.Version())
}
