package set

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
