package qr

// Returns XOR result between two boolean values.
func xor(x, y bool) bool {
	return (x && !y) || (!x && y)
}
