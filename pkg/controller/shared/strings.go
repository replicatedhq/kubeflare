package shared

import "sort"

func StringSlicesMatch(a []string, b []string) bool {
	// sort on a copy to preserve original order
	aa := make([]string, len(a))
	bb := make([]string, len(b))

	copy(aa, a)
	copy(bb, b)

	sort.Strings(aa)
	sort.Strings(bb)

	hasChanged := len(aa) != len(bb)

	if !hasChanged {
		for i, v := range aa {
			if v != bb[i] {
				hasChanged = true
			}
		}
	}

	return !hasChanged
}

func InterfaceArrayToStringArray(in []interface{}) []string {
	out := make([]string, len(in))

	for k, v := range in {
		out[k] = v.(string)
	}

	return out
}

func StringArrayToInterfaceArray(in []string) []interface{} {
	out := make([]interface{}, len(in))

	for k, v := range in {
		out[k] = v
	}

	return out
}

func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
