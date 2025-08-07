package ability

import "fa_machine/common/sets"

type ExtendableState interface {
	GetNexts() map[byte]sets.Set[State]
}
