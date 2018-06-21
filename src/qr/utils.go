package qr

func xor(x, y bool) bool {
	return (x && !y) || (!x && y)
}

func makePixel(size int, px string) []string {
	if size < 1 {
		panic("invalid pixel size")
	}
	for i := 1; i < size; i++ {
		px += px
	}
	res := []string{px}
	if size > 1 {
		for i := 1; i < size; i++ {
			res = append(res, px)
		}
	}
	return res
}
