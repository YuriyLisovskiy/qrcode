package qr

import "testing"

var newMode_TestData = []struct {
	mode     int
	cc0      int
	cc1      int
	cc2      int
	expected modeType
}{
	{
		mode: 0x1, cc0: 10, cc1: 12, cc2: 14,
		expected: modeType{
			modeBits: 0x1,
			numBitsCharCount: [3]int{10, 12, 14},
		},
	},
	{
		mode: 0x2, cc0: 9, cc1: 11, cc2: 13,
		expected: modeType{
			modeBits: 0x2,
			numBitsCharCount: [3]int{9, 11, 13},
		},
	},
	{
		mode: 0x4, cc0: 8, cc1: 16, cc2: 16,
		expected: modeType{
			modeBits: 0x4,
			numBitsCharCount: [3]int{8, 16, 16},
		},
	},
	{
		mode: 0x8, cc0: 8, cc1: 10, cc2: 12,
		expected: modeType{
			modeBits: 0x8,
			numBitsCharCount: [3]int{8, 10, 12},
		},
	},
	{
		mode: 0x7, cc0: 0, cc1: 0, cc2: 0,
		expected: modeType{
			modeBits: 0x7,
			numBitsCharCount: [3]int{0, 0, 0},
		},
	},
}

func Test_newMode(test *testing.T) {
	for _, data := range newMode_TestData {
		actual := newMode(data.mode, data.cc0, data.cc1, data.cc2)
		if actual.modeBits != data.expected.modeBits {
			test.Errorf(
				"mode.Test_newMode:\n\tactual modeBits -> %d\n is not equal to\n\texpected modeBits -> %d",
				actual.modeBits, data.expected.modeBits,
			)
		}
		for i, item := range actual.numBitsCharCount {
			if item != data.expected.numBitsCharCount[i] {
				test.Errorf(
					"mode.Test_newMode:\n\tactual numBitsCharCount -> %d\n is not equal to\n\texpected numBitsCharCount -> %d",
					actual.numBitsCharCount, data.expected.numBitsCharCount,
				)
			}
		}
	}
}

var getModeBits_TestData = []struct {
	mode     modeType
	expected int
}{
	{expected: 0x1, mode: modeType{modeBits: 0x1}},
	{expected: 0x2, mode: modeType{modeBits: 0x2}},
	{expected: 0x4, mode: modeType{modeBits: 0x4}},
	{expected: 0x8, mode: modeType{modeBits: 0x8}},
	{expected: 0x7, mode: modeType{modeBits: 0x7}},
}

func Test_getModeBits(test *testing.T) {
	for _, data := range getModeBits_TestData {
		actual := data.mode.getModeBits()
		if actual != data.expected {
			test.Errorf(
				"mode.Test_getModeBits:\n\tactual modeBits -> %d\n is not equal to\n\texpected modeBits -> %d",
				actual, data.expected,
			)
		}
	}
}

var numCharCountBits_TestData = []struct {
	mode modeType
	version int
	expected int
}{
	{
		version: 1,
		mode: modeType{
			numBitsCharCount: [3]int{10, 12, 14},
		},
		expected: 10,
	},
	{
		version: 9,
		mode: modeType{
			numBitsCharCount: [3]int{10, 12, 14},
		},
		expected: 10,
	},
	{
		version: 10,
		mode: modeType{
			numBitsCharCount: [3]int{10, 12, 14},
		},
		expected: 12,
	},
	{
		version: 26,
		mode: modeType{
			numBitsCharCount: [3]int{10, 12, 14},
		},
		expected: 12,
	},
	{
		version: 27,
		mode: modeType{
			numBitsCharCount: [3]int{10, 12, 14},
		},
		expected: 14,
	},
	{
		version: 40,
		mode: modeType{
			numBitsCharCount: [3]int{10, 12, 14},
		},
		expected: 14,
	},

	{
		version: 7,
		mode: modeType{
			numBitsCharCount: [3]int{10, 12, 14},
		},
		expected: 10,
	},
	{
		version: 17,
		mode: modeType{
			numBitsCharCount: [3]int{10, 12, 14},
		},
		expected: 12,
	},
	{
		version: 37,
		mode: modeType{
			numBitsCharCount: [3]int{10, 12, 14},
		},
		expected: 14,
	},
}

func Test_numCharCountBits(test *testing.T) {
	for _, data := range numCharCountBits_TestData {
		actual := data.mode.numCharCountBits(data.version)
		if actual != data.expected {
			test.Errorf(
				"mode.Test_numCharCountBits:\n\tactual Test_numCharCountBits -> %d\n is not equal to\n\texpected Test_numCharCountBits -> %d",
				actual, data.expected,
			)
		}
	}
}
