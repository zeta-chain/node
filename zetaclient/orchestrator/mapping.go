package orchestrator

import "cmp"

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

// mapDeleteMissingKeys removes signers from the map IF they are not in the presentKeys.
func mapDeleteMissingKeys[K cmp.Ordered, V any](m *map[K]V, presentKeys []K, beforeUnset func(K, V)) {
	presentKeysSet := make(map[K]struct{}, len(presentKeys))
	for _, id := range presentKeys {
		presentKeysSet[id] = struct{}{}
	}

	for key := range *m {
		if _, isPresent := presentKeysSet[key]; !isPresent {
			mapUnset(m, key, beforeUnset)
		}
	}
}
