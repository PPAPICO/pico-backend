package collections

func Map[T, U any](iterable []T, mapper func(T) U) []U {
	result := make([]U, len(iterable))
	for i, item := range iterable {
		result[i] = mapper(item)
	}
	return result
}

func Filter[T any](iterable []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(iterable))
	for _, item := range iterable {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Find[T any](iterable []T, predicate func(T) bool) *T {
	for _, item := range iterable {
		if predicate(item) {
			return &item
		}
	}
	return nil
}
