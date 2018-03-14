package runes

// HasPrefix tests whether the rune slice s begins with prefix.
func HasPrefix(s, prefix []rune) bool {
	return len(s) >= len(prefix) && Equal(s[0:len(prefix)], prefix)
}

// HasSuffix tests whether the rune slice s ends with suffix.
func HasSuffix(s, suffix []rune) bool {
	return len(s) >= len(suffix) && Equal(s[len(s)-len(suffix):], suffix)
}

// Equal returns a boolean reporting whether a and b
// are the same length and contain the same bytes.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

// Index returns the index of the first instance of sep in s, or -1 if sep is not present in s.
func Index(s, sep []rune) int {
	size := len(s)
	sepsize := len(sep)
	if size < sepsize {
		return -1
	}

	for i := 0; i < size; i++ {
		if HasPrefix(s[i:], sep) {
			return i
		}
	}
	return -1
}

func Copy(s []rune) []rune {
	if s == nil {
		return nil
	}
	size := len(s)
	c := make([]rune, size)
	copy(c, s)
	return c
}
