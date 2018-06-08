package utils

type BitBuffer []bool

func (bb BitBuffer) GetBytes() []uint8 {
	k := 1
	if len(bb) % 8 == 0 {
		k = 0
	}
	result := make([]uint8, len(bb) / 8 + k)
	for i := 0; i < len(bb); i++ {
		temp := 0
		if bb[i] {
			temp = 1 << (7 - (uint8(i) & 7))
		}
		result[i >> 3] |= uint8(temp)
	}
	return result
}

func (bb BitBuffer) AppendBits(val uint32, len int) BitBuffer {
	if len < 0 || len > 31 || val >> uint32(len) != 0 {
		panic("value out of range")
	}
	for i := len - 1; i >= 0; i-- {
		bb = append(bb, ((val >> uint32(i)) & 1) != 0)
	}
	return bb
}
