package utils

import "sort"

func ReverseSlice[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func SortSlice[T any](s []T, lessThan func(i, j T) bool) {
	sort.Slice(s, func(i, j int) bool {
		return lessThan(s[i], s[j])
	})
}
