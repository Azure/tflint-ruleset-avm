package attrvalue

type set[T comparable] struct {
	m map[T]struct{}
}

func (t set[T]) equals(vals set[T]) bool {
	if len(t.m) != len(vals.m) {
		return false
	}
	for k, _ := range vals.m {
		_, ok := t.m[k]
		if !ok {
			return false
		}
	}
	return true
}

func newSet[T comparable](ts []T) set[T] {
	s := set[T]{
		m: make(map[T]struct{}),
	}
	for _, v := range ts {
		s.m[v] = struct{}{}
	}
	return s
}
