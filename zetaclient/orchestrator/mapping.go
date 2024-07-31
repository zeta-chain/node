package orchestrator

import "cmp"

// This is a collection of generic functions that can be used to manipulate maps BY pointer.
// This is useful for observers/signers because we want to operate with the same map and not a copy.
// It simplifies code semantics by hiding complexity of accessing map by a pointer.

// mapHas checks if the map contains the given key.
func mapHas[K cmp.Ordered, V any](m *map[K]V, key K) bool {
	_, ok := (*m)[key]
	return ok
}

// mapSet sets the value for the given key in the map
// and (optionally) runs a callback after setting the value.
func mapSet[K cmp.Ordered, V any](m *map[K]V, key K, value V, afterSet func(K, V)) {
	(*m)[key] = value

	if afterSet != nil {
		afterSet(key, value)
	}
}

// mapUnset removes the value for the given key from the map (if exists)
// and optionally runs a callback before removing the value.
func mapUnset[K cmp.Ordered, V any](m *map[K]V, key K, beforeUnset func(K, V)) bool {
	if !mapHas(m, key) {
		return false
	}

	if beforeUnset != nil {
		beforeUnset(key, (*m)[key])
	}

	delete(*m, key)

	return true
}

// mapDeleteMissingKeys removes elements from the map IF they are not in the presentKeys.
func mapDeleteMissingKeys[K cmp.Ordered, V any](m *map[K]V, presentKeys []K, beforeUnset func(K, V)) {
	set := make(map[K]struct{}, len(presentKeys))
	for _, id := range presentKeys {
		set[id] = struct{}{}
	}

	for key := range *m {
		if _, found := set[key]; !found {
			mapUnset(m, key, beforeUnset)
		}
	}
}
