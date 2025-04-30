package main

func Intersection[T comparable](slices ...[]T) []T {
	counts := map[T]int{}
	var result []T

	for _, slice := range slices {
		for _, val := range slice {
			counts[val]++
		}
	}

	for val, count := range counts {
		if count == len(slices) {
			result = append(result, val)
		}
	}

	return result
}
