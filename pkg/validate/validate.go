package validate

// NotEmpty test all provided values to be not empty
func NotEmpty(claims ...string) bool {
	if len(claims) == 0 {
		return false
	}
	for _, s := range claims {
		if s == "" {
			return false
		}
	}
	return true
}
