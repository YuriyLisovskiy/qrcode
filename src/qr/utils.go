package qr

func XOR(x, y bool) bool {
	return (x && !y) || (!x && y)
}
