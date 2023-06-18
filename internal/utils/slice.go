package utils

func SliceContains[T comparable](elem T, slice []T) bool {
	for _, x := range slice {
		if x == elem {
			return true
		}
	}
	return false
}
