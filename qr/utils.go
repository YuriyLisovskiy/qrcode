package qr

func xor(x, y bool) bool {
	return (x && !y) || (!x && y)
}
