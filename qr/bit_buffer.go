package qr

// An appendable sequence of bits (0's and 1's).
type bitBuffer []bool

// Packs this buffer's bits into bytes in big endian,
// padding with '0' bit values, and returns the new vector.
func (bb bitBuffer) getBytes() []uint8 {
	k := 1
	if len(bb)%8 == 0 {
		k = 0
	}
	result := make([]uint8, len(bb)/8+k)
	for i := 0; i < len(bb); i++ {
		temp := 0
		if bb[i] {
			temp = 1 << uint(7 - (i & 7))
		}
		result[i>>3] |= uint8(temp)
	}
	return result
}

// Appends the given number of low bits of the given value
// to this sequence. Requires 0 <= val < 2^len.
func (bb bitBuffer) AppendBits(val uint32, length int) bitBuffer {
	if length < 0 || length > 31 || val>>uint32(length) != 0 {
		panic("value out of range")
	}
	for i := length - 1; i >= 0; i-- {
		temp := (val >> uint(i)) & 1
		bb = append(bb, temp != 0)
	}
	return bb
}
