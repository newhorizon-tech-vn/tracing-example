package util

func MapFunc[T1, T2 any](s []T1, f func(T1) T2) []T2 {
	set := make([]T2, 0)
	for _, v := range s {
		set = append(set, f(v))
	}
	return set
}
