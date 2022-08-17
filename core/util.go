package core

func IndexOf(slice []string, find string) int {
	for index, item := range slice {
		if item == find {
			return index
		}
	}
	return -1
}
