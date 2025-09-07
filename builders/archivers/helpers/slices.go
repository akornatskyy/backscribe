package helpers

func FilterOut(slice []string, keep func(string) bool) []string {
	var result []string
	for _, s := range slice {
		if keep(s) {
			result = append(result, s)
		}
	}
	return result
}
