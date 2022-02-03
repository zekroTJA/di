package di

type hashSet[T comparable] map[T]struct{}

func (t hashSet[T]) Set(v T) bool {
	if t.Has(v) {
		return false
	}
	t[v] = struct{}{}
	return true
}

func (t hashSet[T]) Has(v T) bool {
	_, ok := t[v]
	return ok
}

func (t hashSet[T]) Remove(v T) {
	delete(t, v)
}

func (t hashSet[T]) Values() []T {
	vals := make([]T, 0, len(t))
	for v := range t {
		vals = append(vals, v)
	}
	return vals
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustValue[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
