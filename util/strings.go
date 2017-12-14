package util

func EscapeRune(s string, r rune) string {
	var escaped []rune
	for _, sr := range s {
		if sr == r {
			escaped = append(escaped, '\\')
		}
		escaped = append(escaped, sr)
	}
	return string(escaped)
}
