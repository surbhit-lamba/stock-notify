package utils

func SliceContains[T comparable](elem T, slice []T) bool {
	for _, x := range slice {
		if x == elem {
			return true
		}
	}
	return false
}

func SliceReverse[T any](slice []T) []T {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
