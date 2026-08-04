package collections

func Map[T, U any](iterable []T, mapper func(T) U) []U {
	result := make([]U, len(iterable))
	for i, item := range iterable {
		result[i] = mapper(item)
	}
	return result
}
