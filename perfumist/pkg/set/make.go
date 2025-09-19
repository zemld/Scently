package set

func MakeSet[T comparable](array []T) map[T]struct{} {
	set := make(map[T]struct{})
	for _, item := range array {
		set[item] = struct{}{}
	}
	return set
}
