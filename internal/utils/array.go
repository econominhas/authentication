package utils

func InArray(s string, arr []string) bool {
	for _, a := range arr {
		if s == a {
			return true
		}
	}

	return false
}

func AllInArray(baseArr []string, maybeMissingItems []string) bool {
	for _, s := range maybeMissingItems {
		if !InArray(s, baseArr) {
			return false
		}
	}

	return true
}
