package utils

import "sort"

func SortedMap[T interface{}](m map[string]T, f func(k string, v T)) {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		f(k, m[k])
	}
}
