//  Copyright (c) 2018 Yuriy Lisovskiy
//  Distributed under the Apache License Version 2.0,
//  see the accompanying file LICENSE or https://opensource.org/licenses/Apache-2.0

package qr

import "testing"

var appendBits_TestData = []struct {
	inputVal []uint32
	inputLen []int
	expected bitBuffer
}{
	{
		inputVal: []uint32{1, 2, 3, 5, 7, 1, 4},
		inputLen: []int{2, 4, 8, 21, 16, 1, 9},
		expected: bitBuffer{
			false, true, false, false, true, false, false, false, false, false, false, false, true, true, false, false,
			false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
			false, true, false, true, false, false, false, false, false, false, false, false, false, false, false,
			false, false, true, true, true, true, false, false, false, false, false, false, true, false, false,
		},
	},
	{
		inputVal: []uint32{1, 1, 1, 1, 1, 1, 1},
		inputLen: []int{1, 1, 1, 1, 1, 1, 1},
		expected: bitBuffer{
			true, true, true, true, true, true, true,
		},
	},
	{
		inputVal: []uint32{0, 0, 0, 0, 0, 0, 0},
		inputLen: []int{1, 1, 1, 1, 1, 1, 1},
		expected: bitBuffer{
			false, false, false, false, false, false, false,
		},
	},
}

func Test_AppendBits(test *testing.T) {
	for _, data := range appendBits_TestData {
		actual := bitBuffer{}
		for i := range data.inputVal {
			actual, _ = actual.appendBits(data.inputVal[i], data.inputLen[i])
		}
		if len(actual) != len(data.expected) {
			test.Errorf(
				"bit_buffer.Test_AppendBits:\n\tactual bitBuffer len -> %d\nis not equal to\n\texpected bitBuffer len -> %d",
				len(actual), len(data.expected),
			)
		}
		for i, bit := range actual {
			if bit != data.expected[i] {
				test.Errorf(
					"bit_buffer.Test_AppendBits:\n\tactual bit of bitBuffer-> %t\n is not equal to\n\texpected bit of bitBuffer -> %t",
					bit, data.expected[i],
				)
			}
		}
	}
}

var getBytes_TestData = []struct {
	inputVal []uint32
	inputLen []int
	expected []uint8
}{
	{
		inputVal: []uint32{1, 2, 3, 5, 7, 1, 4, 5, 6},
		inputLen: []int{2, 4, 8, 21, 16, 1, 9, 5, 6},
		expected: []uint8{72, 12, 0, 0, 160, 0, 240, 33, 70},
	},
	{
		inputVal: []uint32{1, 1, 1, 1, 1, 1, 1},
		inputLen: []int{1, 1, 1, 1, 1, 1, 1},
		expected: []uint8{254},
	},
	{
		inputVal: []uint32{0, 0, 0, 0, 0, 0, 0},
		inputLen: []int{1, 1, 1, 1, 1, 1, 1},
		expected: []uint8{0},
	},
}

func Test_getBytes(test *testing.T) {
	for _, data := range getBytes_TestData {
		bitBuf := bitBuffer{}
		for i := range data.inputVal {
			bitBuf, _ = bitBuf.appendBits(data.inputVal[i], data.inputLen[i])
		}
		actual := bitBuf.getBytes()
		if len(actual) != len(data.expected) {
			test.Errorf(
				"bit_buffer.Test_getBytes:\n\tactual uint8 array len -> %d\n is not equal to\n\texpected uint8 array len -> %d",
				len(bitBuf), len(data.expected),
			)
		}
		for i, bit := range actual {
			if bit != data.expected[i] {
				test.Errorf(
					"bit_buffer.Test_getBytes:\n\tactual bit of bitBuffer-> %d\n is not equal to\n\texpected bit of bitBuffer -> %d",
					bit, data.expected[i],
				)
			}
		}
	}
}

var appendBitsErrBitBufValOutOfRange_TestData = []struct {
	inputVal []uint32
	inputLen []int
	expected error
}{
	{
		inputVal: []uint32{1, 2, 3, 5, 7},
		inputLen: []int{-3, -2, -100, -12, -5},
		expected: ErrBitBufValOutOfRange,
	},
	{
		inputVal: []uint32{1, 1, 1, 1, 1, 1, 1},
		inputLen: []int{32, 33, 34, 35, 36, 37, 80},
		expected: ErrBitBufValOutOfRange,
	},
	{
		inputVal: []uint32{12, 3},
		inputLen: []int{16384, 1000000},
		expected: ErrBitBufValOutOfRange,
	},
}

func Test_appendBitsErrBitBufValOutOfRange(test *testing.T) {
	var err error
	for _, data := range appendBitsErrBitBufValOutOfRange_TestData {
		actual := bitBuffer{}
		for i := range data.inputVal {
			_, err = actual.appendBits(data.inputVal[i], data.inputLen[i])
			if err != data.expected {
				test.Errorf(
					"bit_buffer.Test_appendBitsErrBitBufValOutOfRange:\n\t function does not return an error for len %d",
					data.inputLen[i],
				)
			}
		}
	}
}
