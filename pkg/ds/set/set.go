package set

type Set[K comparable] map[K]struct{}

func NewSet[K comparable]() Set[K] {
	return make(Set[K])
}

func (s Set[K]) Contains(k K) bool {
	_, ok := s[k]
	return ok
}

func (s Set[K]) Add(k K) Set[K] {
	s[k] = struct{}{}
	return s
}

func (s Set[K]) AddAll(keys ...K) Set[K] {
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}

func (s Set[K]) Remove(k K) Set[K] {
	delete(s, k)
	return s
}

func (s Set[K]) RemoveAll(keys ...K) Set[K] {
	for _, k := range keys {
		delete(s, k)
	}
	return s
}

func (s Set[K]) Clear() Set[K] {
	for k := range s {
		delete(s, k)
	}
	return s
}

func (s Set[K]) Size() int {
	return len(s)
}

func (s Set[K]) ToSlice() []K {
	keys := make([]K, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	return keys
}

func (s Set[K]) IsEmpty() bool {
	return len(s) == 0
}
