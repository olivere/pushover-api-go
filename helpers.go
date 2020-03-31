package pushover

import "unicode/utf8"

const (
	ellipsis = "…"
)

// cut a UTF-8 encoded string to a maximum length, and
// optionally add a trailing sequence.
func cut(s string, max int, trailing string) string {
	ls := utf8.RuneCountInString(s)
	lt := utf8.RuneCountInString(trailing)
	if ls > max-lt {
		r := []rune(s)[:max-lt]
		r = append(r, []rune(trailing)...)
		return string(r)
	}
	return s
}
