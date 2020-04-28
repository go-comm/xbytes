package xbytes

func roundUp(v int) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func roundLog2(v int) int {
	n := 0
	if v >= 65536 {
		n += 16
		v = (v + 65535) >> 16
	}

	if v >= 256 {
		n += 8
		v = (v + 255) >> 8
	}

	if v >= 16 {
		n += 4
		v = (v + 15) >> 4
	}

	if v >= 4 {
		n += 2
		v = (v + 3) >> 2
	}

	if v >= 2 {
		n++
		v = (v + 1) >> 1
	}

	if v == 2 {
		n++
	}
	return n
}
