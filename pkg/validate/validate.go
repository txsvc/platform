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

// IsMemberOf returns true if claim is part of the list of claims
func IsMemberOf(claim string, claims ...string) bool {
	if len(claims) == 0 {
		return false
	}
	for _, s := range claims {
		if s == claim {
			return true
		}
	}
	return false
}
