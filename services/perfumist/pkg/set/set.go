package set

func MakeSet[T comparable](array []T) map[T]struct{} {
	set := make(map[T]struct{})
	for _, item := range array {
		set[item] = struct{}{}
	}
	return set
}

func Intersect[T comparable](first map[T]struct{}, second map[T]struct{}) map[T]struct{} {
	intersection := make(map[T]struct{})
	for key := range first {
		if _, ok := second[key]; ok {
			intersection[key] = struct{}{}
		}
	}
	return intersection
}

func Union[T comparable](first map[T]struct{}, second map[T]struct{}) map[T]struct{} {
	un := make(map[T]struct{})
	for key := range first {
		un[key] = struct{}{}
	}
	for key := range second {
		un[key] = struct{}{}
	}
	return un
}
