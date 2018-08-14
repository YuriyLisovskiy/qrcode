//  Copyright (c) 2018 Yuriy Lisovskiy
//  Distributed under the Apache License Version 2.0,
//  see the accompanying file LICENSE or https://opensource.org/licenses/Apache-2.0

package qr

import "testing"

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
		actual, _ := makeBytes(&data.input)
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
		actual, _ := makeNumeric(data.input)
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

var makeNumericErr_TestData = []struct {
	input    string
	expected error
}{
	{
		input:    "-1",
		expected: qrSegmentErr("makeNumeric", "string contains non-numeric characters"),
	},
	{
		input:    "d",
		expected: qrSegmentErr("makeNumeric", "string contains non-numeric characters"),
	},
	{
		input:    "-138456789",
		expected: qrSegmentErr("makeNumeric", "string contains non-numeric characters"),
	},
}

func Test_makeNumericErr(test *testing.T) {
	for _, data := range makeNumericErr_TestData {
		_, actual := makeNumeric(data.input)
		if actual.Error() != data.expected.Error() {
			test.Errorf(
				"qr_segment.Test_makeNumericErr:\n\tfunc does not return an error for text %s",
				data.input,
			)
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
		actual, _ := makeAlphanumeric(data.input)
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

var makeAlphanumericErr_TestData = []struct {
	input    string
	expected error
}{
	{
		input:    "SOME TEXT!",
		expected: qrSegmentErr("makeAlphanumeric", "string contains unencodable characters in alphanumeric mode"),
	},
	{
		input:    "SOME ANOTHER TEXT?",
		expected: qrSegmentErr("makeAlphanumeric", "string contains unencodable characters in alphanumeric mode"),
	},
	{
		input:    "SOME TEXT WITH NUMBERS#: 10101010101010101",
		expected: qrSegmentErr("makeAlphanumeric", "string contains unencodable characters in alphanumeric mode"),
	},
}

func Test_makeAlphanumericErr(test *testing.T) {
	for _, data := range makeAlphanumericErr_TestData {
		_, actual := makeAlphanumeric(data.input)
		if actual.Error() != data.expected.Error() {
			test.Errorf(
				"qr_segment.Test_makeAlphanumericErr:\n\tfunc does not return error for data %s",
				data.input,
			)
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
		actual, _ := makeSegments(data.input)
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

var makeEci_TestData = []struct {
	input    int64
	expected qrSegment
}{
	{
		input: 125,
		expected: qrSegment{
			Mode: modeType{
				modeBits:         7,
				numBitsCharCount: [3]int{0, 0, 0},
			},
			NumChars: 0,
			Data: []bool{
				false, true, true, true, true, true, false, true,
			},
		},
	},
	{
		input: 456789,
		expected: qrSegment{
			Mode: modeType{
				modeBits:         7,
				numBitsCharCount: [3]int{0, 0, 0},
			},
			NumChars: 0,
			Data: []bool{
				true, true, false, false, false, true, true, false, true, true, true, true, true, false, false, false,
				false, true, false, true, false, true, false, true,
			},
		},
	},
	{
		input: 453254,
		expected: qrSegment{
			Mode: modeType{
				modeBits:         7,
				numBitsCharCount: [3]int{0, 0, 0},
			},
			NumChars: 0,
			Data: []bool{
				true, true, false, false, false, true, true, false, true, true, true, false, true, false, true, false,
				true, false, false, false, false, true, true, false,
			},
		},
	},
	{
		input: 345789,
		expected: qrSegment{
			Mode: modeType{
				modeBits:         7,
				numBitsCharCount: [3]int{0, 0, 0},
			},
			NumChars: 0,
			Data: []bool{
				true, true, false, false, false, true, false, true, false, true, false, false, false, true, true, false,
				true, false, true, true, true, true, false, true,
			},
		},
	},
	{
		input: 129,
		expected: qrSegment{
			Mode: modeType{
				modeBits:         7,
				numBitsCharCount: [3]int{0, 0, 0},
			},
			NumChars: 0,
			Data: []bool{
				true, false, false, false, false, false, false, false, true, false, false, false, false, false, false, true,
			},
		},
	},
	{
		input: 16383,
		expected: qrSegment{
			Mode: modeType{
				modeBits:         7,
				numBitsCharCount: [3]int{0, 0, 0},
			},
			NumChars: 0,
			Data: []bool{
				true, false, true, true, true, true, true, true, true, true, true, true, true, true, true, true,
			},
		},
	},
	{
		input: 16385,
		expected: qrSegment{
			Mode: modeType{
				modeBits:         7,
				numBitsCharCount: [3]int{0, 0, 0},
			},
			NumChars: 0,
			Data: []bool{
				true, true, false, false, false, false, false, false, false, true, false, false, false, false, false,
				false, false, false, false, false, false, false, false, true,
			},
		},
	},
	{
		input: 999999,
		expected: qrSegment{
			Mode: modeType{
				modeBits:         7,
				numBitsCharCount: [3]int{0, 0, 0},
			},
			NumChars: 0,
			Data: []bool{
				true, true, false, false, true, true, true, true, false, true, false, false, false, false, true, false,
				false, false, true, true, true, true, true, true,
			},
		},
	},
}

func Test_makeEci(test *testing.T) {
	for _, data := range makeEci_TestData {
		actual, _ := makeEci(data.input)
		if actual.NumChars != data.expected.NumChars {
			test.Errorf(
				"qr_segment.Test_makeEci:\n\tactual numChars -> %d\n is not equal to\n\texpected numChars -> %d",
				actual.NumChars, data.expected.NumChars,
			)
		}
		if actual.Mode.modeBits != data.expected.Mode.modeBits {
			test.Errorf(
				"qr_segment.Test_makeEci:\n\tactual Mode -> %d\n is not equal to\n\texpected Mode -> %d",
				actual.Mode.modeBits, data.expected.Mode.modeBits,
			)
		}
		for i, d := range actual.Mode.numBitsCharCount {
			if d != data.expected.Mode.numBitsCharCount[i] {
				test.Errorf(
					"qr_segment.Test_makeEci:\n\tactual Mode -> %d\n is not equal to\n\texpected Mode -> %d",
					d, data.expected.Mode.numBitsCharCount[i],
				)
			}
		}
		for i, d := range actual.Data {
			if d != data.expected.Data[i] {
				test.Errorf(
					"qr_segment.Test_makeEci:\n\tactual Data -> %t\n is not equal to\n\texpected Data -> %t",
					d, data.expected.Data[i],
				)
			}
		}
	}
}

var makeEciErr_TestData = []struct {
	input    int64
	expected error
}{
	{
		input:    -1,
		expected: qrSegmentErr("makeEci", "ECI assignment value out of range"),
	},
	{
		input:    4532544532,
		expected: qrSegmentErr("makeEci", "ECI assignment value out of range"),
	},
}

func Test_makeEciErr(test *testing.T) {
	for _, data := range makeEciErr_TestData {
		_, actual := makeEci(data.input)
		if actual == nil {
			test.Errorf(
				"qr_segment.Test_makeEciErr:\n\tfunc does not return an error for assignVal %d",
				data.input,
			)
		} else {
			if actual.Error() != data.expected.Error() {
				test.Errorf(
					"qr_segment.Test_makeEciErr:\n\tfunc does not return an error for assignVal %d",
					data.input,
				)
			}
		}
	}
}

var getTotalBits_TestData = []struct {
	inputSeg []qrSegment
	inputVer int
	expected int
}{
	{
		inputSeg: []qrSegment{
			{
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
		inputVer: 34,
		expected: 111,
	},
	{
		inputSeg: []qrSegment{
			{
				Mode: modeType{
					modeBits:         7,
					numBitsCharCount: [3]int{0, 0, 0},
				},
				NumChars: 0,
				Data: []bool{
					true, true, false, false, false, true, true, false, true, true, true, false, true, false, true, false,
					true, false, false, false, false, true, true, false,
				},
			},
		},
		inputVer: 40,
		expected: 28,
	},
	{
		inputSeg: []qrSegment{
			{
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
		inputVer: 1,
		expected: 92,
	},
	{
		inputSeg: []qrSegment{
			{
				Mode: modeType{
					modeBits:         4,
					numBitsCharCount: [3]int{8, 16, 16},
				},
				NumChars: 456789,
			},
		},
		inputVer: 2,
		expected: -1,
	},
	{
		inputSeg: []qrSegment{
			{
				Mode: modeType{
					modeBits:         4,
					numBitsCharCount: [3]int{8, 16, 16},
				},
				NumChars: 456789,
			},
		},
		inputVer: 34,
		expected: -1,
	},
	{
		inputSeg: []qrSegment{
			{
				Mode: modeType{
					modeBits:         4,
					numBitsCharCount: [3]int{8, 16, 16},
				},
				NumChars: 456789,
			},
		},
		inputVer: 12,
		expected: -1,
	},
	{
		inputSeg: []qrSegment{
			{
				Mode: modeType{
					modeBits:         4,
					numBitsCharCount: [3]int{8, 16, 16},
				},
				NumChars: 456789,
			},
		},
		inputVer: 39,
		expected: -1,
	},
	{
		inputSeg: []qrSegment{
			{
				Mode: modeType{
					modeBits:         4,
					numBitsCharCount: [3]int{8, 16, 16},
				},
				NumChars: 456789,
			},
		},
		inputVer: 7,
		expected: -1,
	},
}

func Test_getTotalBits(test *testing.T) {
	for _, data := range getTotalBits_TestData {
		actual, _ := getTotalBits(&data.inputSeg, data.inputVer)
		if actual != data.expected {
			test.Errorf(
				"qr_segment.Test_getTotalBits:\n\tactual totalBits -> %d\n is not equal to\n\texpected totalBits -> %d",
				actual, data.expected,
			)
		}
	}
}

var getTotalBitsErr_TestData = []struct {
	inputSeg []qrSegment
	inputVer int
	expected error
}{
	{
		inputSeg: []qrSegment{
			{
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
		inputVer: 0,
		expected: qrSegmentErr("getTotalBits", "version number out of range"),
	},
	{
		inputSeg: []qrSegment{
			{
				Mode: modeType{
					modeBits:         7,
					numBitsCharCount: [3]int{0, 0, 0},
				},
				NumChars: 0,
				Data: []bool{
					true, true, false, false, false, true, true, false, true, true, true, false, true, false, true, false,
					true, false, false, false, false, true, true, false,
				},
			},
		},
		inputVer: 41,
		expected: qrSegmentErr("getTotalBits", "version number out of range"),
	},
}

func Test_getTotalBitsErr(test *testing.T) {
	for _, data := range getTotalBitsErr_TestData {
		_, actual := getTotalBits(&data.inputSeg, data.inputVer)
		if actual.Error() != data.expected.Error() {
			test.Errorf(
				"qr_segment.Test_getTotalBitsErr:\n\tfunc does not return error for version %d",
				data.inputVer,
			)
		}
	}
}

var isAlphanumeric_TestData = []struct {
	input    string
	expected bool
}{
	{
		input:    "ALPHANUMERIC TEXT 123:% ",
		expected: true,
	},
	{
		input:    "ALPHANUMERIC TEXT WITH NON-ALPHANUMERIC SYMBOL!",
		expected: false,
	},
	{
		input:    "Non-alphanumeric text",
		expected: false,
	},
	{
		input:    "13245876543456",
		expected: true,
	},
}

func Test_isAlphanumeric(test *testing.T) {
	for _, data := range isAlphanumeric_TestData {
		actual := isAlphanumeric(data.input)
		if actual != data.expected {
			test.Errorf(
				"qr_segment.Test_isAlphanumeric:\n\tactual totalBits -> %t\n is not equal to\n\texpected totalBits -> %t",
				actual, data.expected,
			)
		}
	}
}

var isNumeric_TestData = []struct {
	input    string
	expected bool
}{
	{
		input:    "TEXT 123:% ",
		expected: false,
	},
	{
		input:    "TEXT WITH NON-NUMERIC SYMBOLS",
		expected: false,
	},
	{
		input:    "Non-numeric text",
		expected: false,
	},
	{
		input:    "13245876543456",
		expected: true,
	},
}

func Test_isNumeric(test *testing.T) {
	for _, data := range isNumeric_TestData {
		actual := isNumeric(data.input)
		if actual != data.expected {
			test.Errorf(
				"qr_segment.Test_isNumeric:\n\tactual totalBits -> %t\n is not equal to\n\texpected totalBits -> %t",
				actual, data.expected,
			)
		}
	}
}

var getMode_TestData = []struct {
	input    qrSegment
	expected modeType
}{
	{
		input: qrSegment{
			Mode: modeType{
				modeBits:         2,
				numBitsCharCount: [3]int{9, 11, 13},
			},
		},
		expected: modeType{
			modeBits:         2,
			numBitsCharCount: [3]int{9, 11, 13},
		},
	},
	{
		input: qrSegment{
			Mode: modeType{
				modeBits:         4,
				numBitsCharCount: [3]int{8, 16, 16},
			},
		},
		expected: modeType{
			modeBits:         4,
			numBitsCharCount: [3]int{8, 16, 16},
		},
	},
	{
		input: qrSegment{
			Mode: modeType{
				modeBits:         2,
				numBitsCharCount: [3]int{9, 11, 13},
			},
		},
		expected: modeType{
			modeBits:         2,
			numBitsCharCount: [3]int{9, 11, 13},
		},
	},
}

func Test_getMode(test *testing.T) {
	for _, data := range getMode_TestData {
		actual := data.input.getMode()
		if actual.modeBits != data.expected.modeBits {
			test.Errorf(
				"qr_segment.Test_getMode:\n\tactual Mode -> %d\n is not equal to\n\texpected Mode -> %d",
				actual.modeBits, data.expected.modeBits,
			)
		}
		for i, d := range actual.numBitsCharCount {
			if d != data.expected.numBitsCharCount[i] {
				test.Errorf(
					"qr_segment.Test_getMode:\n\tactual Mode -> %d\n is not equal to\n\texpected Mode -> %d",
					d, data.expected.numBitsCharCount[i],
				)
			}
		}
	}
}

var getNumChars_TestData = []struct {
	input    qrSegment
	expected int
}{
	{
		input: qrSegment{
			NumChars: 1,
		},
		expected: 1,
	},
	{
		input: qrSegment{
			NumChars: 10,
		},
		expected: 10,
	},
	{
		input: qrSegment{
			NumChars: 0,
		},
		expected: 0,
	},
	{
		input: qrSegment{
			NumChars: 121,
		},
		expected: 121,
	},
	{
		input: qrSegment{
			NumChars: 98,
		},
		expected: 98,
	},
}

func Test_getNumChars(test *testing.T) {
	for _, data := range getNumChars_TestData {
		actual := data.input.getNumChars()
		if actual != data.expected {
			test.Errorf(
				"qr_segment.Test_getNumChars:\n\tactual numChars -> %d\n is not equal to\n\texpected numChars -> %d",
				actual, data.expected,
			)
		}
	}
}

var getData_TestData = []struct {
	input    qrSegment
	expected []bool
}{
	{
		input: qrSegment{
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
		expected: []bool{
			true, false, true, false, false, false, false, false, true, false, false, false, true, true, true, true,
			true, false, true, true, false, false, true, true, false, false, true, true, true, false, false, false,
			true, false, true, false, true, false, false, true, false, true, true, true, false, true, true, true,
			false, true,
		},
	},
	{
		input: qrSegment{
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
		expected: []bool{
			true, false, true, false, false, false, false, false, true, false, false, false, true, true, true, true,
			true, false, true, true, false, false, true, true, false, false, true, false, true, true, true, true,
			false, true, false, false, false, false, true, false, false, false, true, true, true, false, true, false,
			false, true, false, true, false, true, false, false, true, false, true, false, false, true, false, false,
			false, true, true, true, false, false, true, true, true, false, false, false, true, false, true, false,
			true, false, false, true, false, true, true, true, false, true, true, true, false, true,
		},
	},
	{
		input: qrSegment{
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
		expected: []bool{
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
}

func Test_getData(test *testing.T) {
	for _, data := range getData_TestData {
		actual := data.input.getData()
		for i, d := range *actual {
			if d != data.expected[i] {
				test.Errorf(
					"qr_segment.Test_getData:\n\tactual numChars -> %t\n is not equal to\n\texpected numChars -> %t",
					d, data.expected[i],
				)
			}
		}
	}
}
