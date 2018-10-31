package acl

import "sort"

// determineIndex returns the index of the given value, otherwise -1
// note, that this function only works with an alphabetically sorted haystack since it uses a binary search
// to resolve the given value.
func determineIndex(haystack []string, value string) int {
	index := sort.SearchStrings(haystack, value)

	if index < len(haystack) && haystack[index] == value {
		return index
	}

	return -1
}
