package set

func Intersect[T comparable](first map[T]struct{}, second map[T]struct{}) map[T]struct{} {
	intersection := make(map[T]struct{})
	for key := range first {
		if _, ok := second[key]; ok {
			intersection[key] = struct{}{}
		}
	}
	return intersection
}
