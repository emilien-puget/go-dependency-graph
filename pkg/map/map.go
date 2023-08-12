package mymap

import "sort"

func OrderedKeys[v any](tab map[string]v) []string {
	keys := make([]string, 0, len(tab))
	for k := range tab {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
