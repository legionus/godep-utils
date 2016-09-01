package utils

func InSliceString(n string, list []string) bool {
	for _, s := range list {
		if n == s {
			return true
		}
	}
	return false
}
