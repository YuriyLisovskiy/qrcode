//  Copyright (c) 2018 Yuriy Lisovskiy
//  Distributed under the Apache License Version 2.0,
//  see the accompanying file LICENSE or https://opensource.org/licenses/Apache-2.0

package qr

// The mode field of a segment.
// Provides methods to retrieve closely related values.
type modeType struct {
	modeBits         int
	numBitsCharCount [3]int
}

// Constructor
func newMode(mode, cc0, cc1, cc2 int) modeType {
	newMode := modeType{modeBits: mode}
	newMode.numBitsCharCount[0] = cc0
	newMode.numBitsCharCount[1] = cc1
	newMode.numBitsCharCount[2] = cc2
	return newMode
}

// Returns the mode indicator bits, which is an unsigned 4-bit value (range 0 to 15).
func (m modeType) getModeBits() int {
	return m.modeBits
}

// Returns the bit width of the segment character count field for this mode object at the given version number.
func (m modeType) numCharCountBits(ver int) int {
	if 1 <= ver && ver <= 9 {
		return m.numBitsCharCount[0]
	}
	if 10 <= ver && ver <= 26 {
		return m.numBitsCharCount[1]
	}
	if 27 <= ver && ver <= 40 {
		return m.numBitsCharCount[2]
	}
	panic("package qr: modeType.numCharCountBits: version number out of range")
}
