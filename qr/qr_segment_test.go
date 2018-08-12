//  Copyright (c) 2018 Yuriy Lisovskiy
//  Distributed under the Apache License Version 2.0,
//  see the accompanying file LICENSE or https://opensource.org/licenses/Apache-2.0

package qr

import (
	"testing"
	)

var makeBytes_TestData = []struct {
	input    []uint8
	expected qrSegment
}{
	{
		input: []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		expected: qrSegment{
			Mode: modeType{
				modeBits:         4,
				numBitsCharCount: [3]int{8, 16, 16},
			},
			NumChars: 10,
			Data: []bool{
				false, false, false, false, false, false, false, true, false, false, false, false, false, false, true,
				false, false, false, false, false, false, false, true, true, false, false, false, false, false, true,
				false, false, false, false, false, false, false, true, false, true, false, false, false, false, false,
				true, true, false, false, false, false, false, false, true, true, true, false, false, false, false, true,
				false, false, false, false, false, false, false, true, false, false, true, false, false, false, false,
				true, false, true, false,
			},
		},
	},
	{
		input: []uint8{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		expected: qrSegment{
			Mode: modeType{
				modeBits:         4,
				numBitsCharCount: [3]int{8, 16, 16},
			},
			NumChars: 10,
			Data: []bool{
				false, false, false, false, false, false, false, true, false, false, false, false, false, false, false,
				true, false, false, false, false, false, false, false, true, false, false, false, false, false, false,
				false, true, false, false, false, false, false, false, false, true, false, false, false, false, false,
				false, false, true, false, false, false, false, false, false, false, true, false, false, false, false,
				false, false, false, true, false, false, false, false, false, false, false, true, false, false, false,
				false, false, false, false, true,
			},
		},
	},
	{
		input: []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		expected: qrSegment{
			Mode: modeType{
				modeBits:         4,
				numBitsCharCount: [3]int{8, 16, 16},
			},
			NumChars: 10,
			Data: []bool{
				false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false, false, false, false, false, false, false, false,
				false, false, false, false, false,
			},
		},
	},
}

func Test_makeBytes(test *testing.T) {
	for _, data := range makeBytes_TestData {
		actual := makeBytes(&data.input)
		if actual.NumChars != data.expected.NumChars {
			test.Errorf(
				"qr_segment.Test_makeBytes:\n\tactual numChars -> %d\n is not equal to\n\texpected numChars -> %d",
				actual.NumChars, data.expected.NumChars,
			)
		}
		if actual.Mode.modeBits != data.expected.Mode.modeBits {
			test.Errorf(
				"qr_segment.Test_makeBytes:\n\tactual Mode -> %d\n is not equal to\n\texpected Mode -> %d",
				actual.Mode.modeBits, data.expected.Mode.modeBits,
			)
		}
		for i, d := range actual.Mode.numBitsCharCount {
			if d != data.expected.Mode.numBitsCharCount[i] {
				test.Errorf(
					"qr_segment.Test_makeBytes:\n\tactual Mode -> %d\n is not equal to\n\texpected Mode -> %d",
					d, data.expected.Mode.numBitsCharCount[i],
				)
			}
		}
		for i, d := range actual.Data {
			if d != data.expected.Data[i] {
				test.Errorf(
					"qr_segment.Test_makeBytes:\n\tactual Data -> %t\n is not equal to\n\texpected Data -> %t",
					d, data.expected.Data[i],
				)
			}
		}
	}
}

var makeNumeric_TestData = []struct {
	input    string
	expected qrSegment
}{
	{
		input: "1234567890",
		expected: qrSegment{
			Mode: modeType{
				modeBits:         1,
				numBitsCharCount: [3]int{10, 12, 14},
			},
			NumChars: 10,
			Data: []bool{
				false, false, false, true, true, true, true, false, true, true, false, true, true, true, false, false,
				true, false, false, false, true, true, false, false, false, true, false, true, false, true, false, false,
				false, false,
			},
		},
	},
	{
		input: "0987654321",
		expected: qrSegment{
			Mode: modeType{
				modeBits:         1,
				numBitsCharCount: [3]int{10, 12, 14},
			},
			NumChars: 10,
			Data: []bool{
				false, false, false, true, true, false, false, false, true, false, true, false, true, true, true, true,
				true, true, false, true, false, true, true, false, true, true, false, false, false, false, false, false,
				false, true,
			},
		},
	},
	{
		input: "10101010101010101",
		expected: qrSegment{
			Mode: modeType{
				modeBits:         1,
				numBitsCharCount: [3]int{10, 12, 14},
			},
			NumChars: 17,
			Data: []bool{
				false, false, false, true, true, false, false, true, false, true, false, false, false, false, false, false,
				true, false, true, false, false, false, false, true, true, false, false, true, false, true, false, false,
				false, false, false, false, true, false, true, false, false, false, false, true, true, false, false, true,
				false, true, false, false, false, false, false, false, true,
			},
		},
	},
}

func Test_makeNumeric(test *testing.T) {
	for _, data := range makeNumeric_TestData {
		actual := makeNumeric(data.input)
		if actual.NumChars != data.expected.NumChars {
			test.Errorf(
				"qr_segment.Test_makeNumeric:\n\tactual numChars -> %d\n is not equal to\n\texpected numChars -> %d",
				actual.NumChars, data.expected.NumChars,
			)
		}
		if actual.Mode.modeBits != data.expected.Mode.modeBits {
			test.Errorf(
				"qr_segment.Test_makeNumeric:\n\tactual Mode -> %d\n is not equal to\n\texpected Mode -> %d",
				actual.Mode.modeBits, data.expected.Mode.modeBits,
			)
		}
		for i, d := range actual.Mode.numBitsCharCount {
			if d != data.expected.Mode.numBitsCharCount[i] {
				test.Errorf(
					"qr_segment.Test_makeNumeric:\n\tactual Mode -> %d\n is not equal to\n\texpected Mode -> %d",
					d, data.expected.Mode.numBitsCharCount[i],
				)
			}
		}
		for i, d := range actual.Data {
			if d != data.expected.Data[i] {
				test.Errorf(
					"qr_segment.Test_makeNumeric:\n\tactual Data -> %t\n is not equal to\n\texpected Data -> %t",
					d, data.expected.Data[i],
				)
			}
		}
	}
}

var makeAlphanumeric_TestData = []struct {
	input    string
	expected qrSegment
}{
	{
		input: "SOME TEXT",
		expected: qrSegment{
			Mode: modeType{
				modeBits:         2,
				numBitsCharCount: [3]int{9, 11, 13},
			},
			NumChars: 9,
			Data: []bool{
				true, false, true, false, false, false, false, false, true, false, false, false, true, true, true, true,
				true, false, true, true, false, false, true, true, false, false, true, true, true, false, false, false,
				true, false, true, false, true, false, false, true, false, true, true, true, false, true, true, true,
				false, true,
			},
		},
	},
	{
		input: "SOME ANOTHER TEXT",
		expected: qrSegment{
			Mode: modeType{
				modeBits:         2,
				numBitsCharCount: [3]int{9, 11, 13},
			},
			NumChars: 17,
			Data: []bool{
				true, false, true, false, false, false, false, false, true, false, false, false, true, true, true, true,
				true, false, true, true, false, false, true, true, false, false, true, false, true, true, true, true,
				false, true, false, false, false, false, true, false, false, false, true, true, true, false, true, false,
				false, true, false, true, false, true, false, false, true, false, true, false, false, true, false, false,
				false, true, true, true, false, false, true, true, true, false, false, false, true, false, true, false,
				true, false, false, true, false, true, true, true, false, true, true, true, false, true,
			},
		},
	},
	{
		input: "SOME TEXT WITH NUMBERS: 10101010101010101",
		expected: qrSegment{
			Mode: modeType{
				modeBits:         2,
				numBitsCharCount: [3]int{9, 11, 13},
			},
			NumChars: 41,
			Data: []bool{
				true, false, true, false, false, false, false, false, true, false, false, false, true, true, true, true,
				true, false, true, true, false, false, true, true, false, false, true, true, true, false, false, false,
				true, false, true, false, true, false, false, true, false, true, true, true, true, false, true, false,
				false, true, true, true, true, false, true, true, false, true, true, false, true, true, false, false,
				true, false, true, false, true, false, false, true, false, true, false, true, false, true, true, false,
				false, true, true, false, true, false, true, true, true, false, true, false, true, false, true, true,
				true, false, false, false, false, true, true, true, true, true, true, true, false, true, true, false,
				false, true, true, false, true, true, false, true, true, true, true, true, true, true, true, false,
				false, false, false, false, false, false, false, false, false, true, false, true, true, false, true,
				false, false, false, false, false, true, false, true, true, false, true, false, false, false, false,
				false, true, false, true, true, false, true, false, false, false, false, false, true, false, true,
				true, false, true, false, false, false, false, false, true, false, true, true, false, true, false,
				false, false, false, false, true, false, true, true, false, true, false, false, false, false, false,
				true, false, true, true, false, true, false, false, false, false, false, true, false, true, true,
				false, true, false, false, false, false, false, true,
			},
		},
	},
}

func Test_makeAlphanumeric(test *testing.T) {
	for _, data := range makeAlphanumeric_TestData {
		actual := makeAlphanumeric(data.input)
		if actual.NumChars != data.expected.NumChars {
			test.Errorf(
				"qr_segment.Test_makeAlphanumeric:\n\tactual numChars -> %d\n is not equal to\n\texpected numChars -> %d",
				actual.NumChars, data.expected.NumChars,
			)
		}
		if actual.Mode.modeBits != data.expected.Mode.modeBits {
			test.Errorf(
				"qr_segment.Test_makeAlphanumeric:\n\tactual Mode -> %d\n is not equal to\n\texpected Mode -> %d",
				actual.Mode.modeBits, data.expected.Mode.modeBits,
			)
		}
		for i, d := range actual.Mode.numBitsCharCount {
			if d != data.expected.Mode.numBitsCharCount[i] {
				test.Errorf(
					"qr_segment.Test_makeAlphanumeric:\n\tactual Mode -> %d\n is not equal to\n\texpected Mode -> %d",
					d, data.expected.Mode.numBitsCharCount[i],
				)
			}
		}
		for i, d := range actual.Data {
			if d != data.expected.Data[i] {
				test.Errorf(
					"qr_segment.Test_makeAlphanumeric:\n\tactual Data -> %t\n is not equal to\n\texpected Data -> %t",
					d, data.expected.Data[i],
				)
			}
		}
	}
}

var makeSegments_TestData = []struct {
	input    string
	expected []qrSegment
}{
	{
		input: "SOME TEXT WITH NUMBERS: 10101010101010101",
		expected: []qrSegment{
			{
				Mode: modeType{
					modeBits:         2,
					numBitsCharCount: [3]int{9, 11, 13},
				},
				NumChars: 41,
				Data: []bool{
					true, false, true, false, false, false, false, false, true, false, false, false, true, true, true, true,
					true, false, true, true, false, false, true, true, false, false, true, true, true, false, false, false,
					true, false, true, false, true, false, false, true, false, true, true, true, true, false, true, false,
					false, true, true, true, true, false, true, true, false, true, true, false, true, true, false, false,
					true, false, true, false, true, false, false, true, false, true, false, true, false, true, true, false,
					false, true, true, false, true, false, true, true, true, false, true, false, true, false, true, true,
					true, false, false, false, false, true, true, true, true, true, true, true, false, true, true, false,
					false, true, true, false, true, true, false, true, true, true, true, true, true, true, true, false,
					false, false, false, false, false, false, false, false, false, true, false, true, true, false, true,
					false, false, false, false, false, true, false, true, true, false, true, false, false, false, false,
					false, true, false, true, true, false, true, false, false, false, false, false, true, false, true,
					true, false, true, false, false, false, false, false, true, false, true, true, false, true, false,
					false, false, false, false, true, false, true, true, false, true, false, false, false, false, false,
					true, false, true, true, false, true, false, false, false, false, false, true, false, true, true,
					false, true, false, false, false, false, false, true,
				},
			},
		},
	},
	{
		input: "1234567890",
		expected: []qrSegment{
			{
				Mode: modeType{
					modeBits:         1,
					numBitsCharCount: [3]int{10, 12, 14},
				},
				NumChars: 10,
				Data: []bool{
					false, false, false, true, true, true, true, false, true, true, false, true, true, true, false, false,
					true, false, false, false, true, true, false, false, false, true, false, true, false, true, false, false,
					false, false,
				},
			},
		},
	},
	{
		input: "Some text to make",
		expected: []qrSegment{
			{
				Mode: modeType{
					modeBits:         4,
					numBitsCharCount: [3]int{8, 16, 16},
				},
				NumChars: 17,
				Data: []bool{
					false, true, false, true, false, false, true, true, false, true, true, false, true, true, true, true,
					false, true, true, false, true, true, false, true, false, true, true, false, false, true, false, true,
					false, false, true, false, false, false, false, false, false, true, true, true, false, true, false,
					false, false, true, true, false, false, true, false, true, false, true, true, true, true, false,
					false, false, false, true, true, true, false, true, false, false, false, false, true, false, false,
					false, false, false, false, true, true, true, false, true, false, false, false, true, true, false,
					true, true, true, true, false, false, true, false, false, false, false, false, false, true, true,
					false, true, true, false, true, false, true, true, false, false, false, false, true, false, true,
					true, false, true, false, true, true, false, true, true, false, false, true, false, true,
				},
			},
		},
	},
}

func Test_makeSegments(test *testing.T) {
	for _, data := range makeSegments_TestData {
		actual := makeSegments(data.input)
		if actual[0].NumChars != data.expected[0].NumChars {
			test.Errorf(
				"qr_segment.Test_makeSegments:\n\tactual numChars -> %d\n is not equal to\n\texpected numChars -> %d",
				actual[0].NumChars, data.expected[0].NumChars,
			)
		}
		if actual[0].Mode.modeBits != data.expected[0].Mode.modeBits {
			test.Errorf(
				"qr_segment.Test_makeSegments:\n\tactual Mode -> %d\n is not equal to\n\texpected Mode -> %d",
				actual[0].Mode.modeBits, data.expected[0].Mode.modeBits,
			)
		}
		for i, d := range actual[0].Mode.numBitsCharCount {
			if d != data.expected[0].Mode.numBitsCharCount[i] {
				test.Errorf(
					"qr_segment.Test_makeSegments:\n\tactual Mode -> %d\n is not equal to\n\texpected Mode -> %d",
					d, data.expected[0].Mode.numBitsCharCount[i],
				)
			}
		}
		for i, d := range actual[0].Data {
			if d != data.expected[0].Data[i] {
				test.Errorf(
					"qr_segment.Test_makeSegments:\n\tactual Data -> %t\n is not equal to\n\texpected Data -> %t",
					d, data.expected[0].Data[i],
				)
			}
		}
	}
}
