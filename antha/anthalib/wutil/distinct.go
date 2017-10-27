package wutil

// remove copies in place

func SADistinct(sa []string) []string {
	m := make(map[string]bool, len(sa))

	for _, s := range sa {
		m[s] = true
	}

	r := make([]string, 0, len(m))

	for _, s := range sa {
		if m[s] {
			r = append(r, s)
			delete(m, s)
		}
	}

	return r
}
