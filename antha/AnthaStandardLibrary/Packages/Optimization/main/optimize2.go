package main

func comb(n, r int) []string {
	if r == 0 {
		s := ""
		for i := 0; i < n; i++ {
			s += "0"
		}

		return []string{s}
	}

	ret := make([]string, 0, n)

	for i := 0; i < n; i++ {
		s := ""
		for j := 0; j < i; j++ {
			s += "0"
		}
		s += "1"

		c := comb(n-i-1, r-1)

		for _, v := range c {
			ret = append(ret, s+v)
		}
	}

	return ret

}
