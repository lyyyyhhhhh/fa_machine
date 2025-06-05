package sets

type Set[T comparable] map[T]struct{}

var emptyStruct = struct{}{}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Add(val T) bool {
	if _, exist := s[val]; !exist {
		s[val] = emptyStruct
		return true
	}
	return false
}

func (s Set[T]) AddAll(vals []T) bool {
	var added bool
	for _, val := range vals {
		if _, exist := s[val]; !exist {
			s[val] = emptyStruct
			added = true
		}
	}
	return added
}

func (s Set[T]) Remove(val T) {
	if _, exist := s[val]; exist {
		delete(s, val)
	}
}

func (s Set[T]) ToSlices() []T {
	ret := make([]T, len(s))
	idx := 0
	for val := range s {
		ret[idx] = val
		idx++
	}
	return ret
}

func (s Set[T]) Contains(val T) bool {
	_, exist := s[val]
	return exist
}

func (s Set[T]) ContainsAll(vals []T) bool {
	for _, val := range vals {
		_, exist := s[val]
		if !exist {
			return false
		}
	}
	return true
}

func (s Set[T]) ContainsOne(vals []T) bool {
	for _, val := range vals {
		_, exist := s[val]
		if exist {
			return true
		}
	}
	return false
}
