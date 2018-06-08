package utils

type Mode struct {
	ModeBits         int
	NumBitsCharCount [3]int
}

func NewMode(mode, cc0, cc1, cc2 int) Mode {
	newMode := Mode{ModeBits: mode}
	newMode.NumBitsCharCount[0] = cc0
	newMode.NumBitsCharCount[1] = cc1
	newMode.NumBitsCharCount[2] = cc2
	return newMode
}

func (m Mode) GetModeBits() int {
	return m.ModeBits
}

func (m Mode) NumCharCountBits(ver int) int {
	if 1 <= ver && ver <= 9 {
		return m.NumBitsCharCount[0]
	}
	if 10 <= ver && ver <= 26 {
		return m.NumBitsCharCount[1]
	}
	if 27 <= ver && ver <= 40 {
		return m.NumBitsCharCount[2]
	}
	panic("version number out of range")
}
