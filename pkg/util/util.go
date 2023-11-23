package util

func SliceToSet[T comparable](a []T) map[T]struct{} {
	m := make(map[T]struct{}, len(a))
	for _, v := range a {
		m[v] = struct{}{}
	}

	return m
}

func SliceFilter[T any](a []T, keep func(T) bool) []T {
	res := make([]T, 0)
	for _, v := range a {
		if keep(v) {
			res = append(res, v)
		}
	}

	return res
}
