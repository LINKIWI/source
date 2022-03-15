package util

// MergeMaps takes one or more maps as input and merges them, in sequential order, into a new map.
// Note that later values take precedence in the merge, and that a cloned map is returned.
func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	for _, item := range maps {
		for k, v := range item {
			merged[k] = v
		}
	}

	return merged
}
