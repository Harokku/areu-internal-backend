package utils

// Contains check if {e} string exist in {s} slice of strings
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
