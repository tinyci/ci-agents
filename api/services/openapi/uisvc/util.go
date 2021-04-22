package uisvc

func stringDeref(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
