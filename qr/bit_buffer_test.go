package qr

import (
	"testing"
	)

var (
	AppendBits_TestData = []struct {
		inputVal []uint32
		inputLen []int
		expected bitBuffer
	}{
		{
			inputVal: []uint32{1, 2, 3, 5, 7, 1, 4},
			inputLen: []int{2, 4, 8, 21, 16, 1, 9},
			expected: bitBuffer{
				false, true, false, false, true, false, false, false, false, false, false, false, true, true, false,
				false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
				false, false, true, false, true, false, false, false, false, false, false, false, false, false, false,
				false, false, false, true, true, true, true, false, false, false, false, false, false, true, false, false,
			},
		},
		{
			inputVal: []uint32{1, 1, 1, 1, 1, 1, 1},
			inputLen: []int{1, 1, 1, 1, 1, 1, 1},
			expected: bitBuffer{
				false, true, false, false, true, false, false, false, false, false, false, false, true, true, false,
				false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
				false, false, true, false, true, false, false, false, false, false, false, false, false, false, false,
				false, false, false, true, true, true, true, false, false, false, false, false, false, true, false, false,
				true, true, true, true, true, true, true,
			},
		},
		{
			inputVal: []uint32{0, 0, 0, 0, 0, 0, 0},
			inputLen: []int{1, 1, 1, 1, 1, 1, 1},
			expected: bitBuffer{
				false, true, false, false, true, false, false, false, false, false, false, false, true, true, false, false,
				false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
				false, true, false, true, false, false, false, false, false, false, false, false, false, false, false, false,
				false, true, true, true, true, false, false, false, false, false, false, true, false, false, true, true,
				true, true, true, true, true, false, false, false, false, false, false, false,
			},
		},
	}
)

func Test_AppendBits(test *testing.T) {
	actual := bitBuffer{}
	for _, data := range AppendBits_TestData {
		for i := range data.inputVal {
			actual = actual.AppendBits(data.inputVal[i], data.inputLen[i])
		}
		if len(actual) != len(data.expected) {
			test.Errorf(
				"qr.Test_AppendBits: actual bitBuffer len -> %d is not equal to expected bitBuffer len -> %d",
				len(actual), len(data.expected),
			)
		}
		for i, bit := range actual {
			if bit != data.expected[i] {
				test.Errorf(
					"qr.Test_AppendBits: actual bit of bitBuffer-> %t is not equal to expected bit of bitBuffer -> %t",
					bit, data.expected[i],
				)
			}
		}
	}
}

var (
	getBytes_TestData = []struct {
		inputVal []uint32
		inputLen []int
		expected []uint8
	}{
		{
			inputVal: []uint32{1, 2, 3, 5, 7, 1, 4},
			inputLen: []int{2, 4, 8, 21, 16, 1, 9},
			expected: []uint8{72, 12, 0, 0, 160, 0, 240, 32},
		},
		{
			inputVal: []uint32{1, 1, 1, 1, 1, 1, 1},
			inputLen: []int{1, 1, 1, 1, 1, 1, 1},
			expected: []uint8{72, 12, 0, 0, 160, 0, 240, 39, 240},
		},
		{
			inputVal: []uint32{0, 0, 0, 0, 0, 0, 0},
			inputLen: []int{1, 1, 1, 1, 1, 1, 1},
			expected: []uint8{72, 12, 0, 0, 160, 0, 240, 39, 240, 0},
		},
	}
)

func Test_getBytes(test *testing.T) {
	bitBuf := bitBuffer{}
	for _, data := range getBytes_TestData {
		for i := range data.inputVal {
			bitBuf = bitBuf.AppendBits(data.inputVal[i], data.inputLen[i])
		}
		actual := bitBuf.getBytes()
		if len(actual) != len(data.expected) {
			test.Errorf(
				"qr.Test_getBytes: actual uint8 array len -> %d is not equal to expected uint8 array len -> %d",
				len(bitBuf), len(data.expected),
			)
		}
		for i, bit := range actual {
			if bit != data.expected[i] {
				test.Errorf(
					"qr.Test_getBytes: actual bit of bitBuffer-> %d is not equal to expected bit of bitBuffer -> %d",
					bit, data.expected[i],
				)
			}
		}
	}
}
