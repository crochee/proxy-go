package internal

// EnsureLeadingSlash makes sure str has lead slash
func EnsureLeadingSlash(str string) string {
	if str == "" {
		return str
	}

	if str[0] == '/' {
		return str
	}

	return "/" + str
}
