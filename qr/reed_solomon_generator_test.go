//  Copyright (c) 2018 Yuriy Lisovskiy
//  Distributed under the Apache License Version 2.0,
//  see the accompanying file LICENSE or https://opensource.org/licenses/Apache-2.0

package qr

import (
	"testing"
		)

var newReedSolomonGenerator_TestData = []struct {
	degree   int
	expected reedSolomonGenerator
}{
	{
		degree: 5,
		expected: reedSolomonGenerator{
			coefficients: []uint8{
				31, 198, 63, 147, 116,
			},
		},
	},
	{
		degree: 1,
		expected: reedSolomonGenerator{
			coefficients: []uint8{
				1,
			},
		},
	},
	{
		degree: 26,
		expected: reedSolomonGenerator{
			coefficients: []uint8{
				246, 51, 183, 4, 136, 98, 199, 152, 77, 56, 206, 24, 145, 40, 209, 117, 233, 42, 135, 68, 70, 144, 146, 77, 43, 94,
			},
		},
	},
	{
		degree: 2,
		expected: reedSolomonGenerator{
			coefficients: []uint8{
				3, 2,
			},
		},
	},
}

func Test_newReedSolomonGenerator(test *testing.T) {
	for _, data := range newReedSolomonGenerator_TestData {
		actual, _ := newReedSolomonGenerator(data.degree)
		if len(actual.coefficients) != data.degree {
			test.Errorf(
				"reed_solomon_generator.Test_newReedSolomonGenerator:\n\tlen of reedSolomonGenerator's coefficients -> %d\n is not equal to\n\treedSolomonGenerator degree -> %d",
				len(actual.coefficients), data.degree,
			)
		}
		if len(actual.coefficients) != len(data.expected.coefficients) {
			test.Errorf(
				"reed_solomon_generator.Test_newReedSolomonGenerator:\n\tactual len of coefficients -> %d\n is not equal to\n\texpected len of coefficients -> %d",
				len(actual.coefficients), len(data.expected.coefficients),
			)
		}
		for i, coefficient := range actual.coefficients {
			if coefficient != data.expected.coefficients[i] {
				test.Errorf(
					"reed_solomon_generator.Test_newReedSolomonGenerator:\n\tactual coefficient -> %d\n is not equal to\n\texpected coefficient -> %d",
					coefficient, data.expected.coefficients[i],
				)
			}
		}
	}
}

var newReedSolomonGeneratorErr_TestData = []struct {
	degree   int
	expected error
}{
	{
		degree: 0,
		expected: rsgErr("newReedSolomonGenerator", "degree out of range"),
	},
	{
		degree: 256,
		expected: rsgErr("newReedSolomonGenerator", "degree out of range"),
	},
}

func Test_newReedSolomonGeneratorErr(test *testing.T) {
	for _, data := range newReedSolomonGeneratorErr_TestData {
		_, actual := newReedSolomonGenerator(data.degree)
		if actual.Error() != data.expected.Error() {
			test.Errorf(
				"reed_solomon_generator.Test_newReedSolomonGeneratorErr:\n\tfunc does not return an error for degree %d",
				data.degree,
			)
		}
	}
}

var getRemainder_TestData = []struct {
	rsg      reedSolomonGenerator
	data     *[]uint8
	expected []uint8
}{
	{
		rsg: reedSolomonGenerator{
			coefficients: []uint8{
				246, 51, 183, 4, 136, 98, 199, 152, 77, 56, 206, 24, 145, 40, 209, 117, 233, 42, 135, 68, 70, 144, 146, 77, 43, 94,
			},
		},
		data: &[]uint8{
			1, 2, 3, 4, 5, 6, 7,
		},
		expected: []uint8{
			174, 125, 216, 144, 60, 116, 108, 34, 93, 29, 125, 19, 207, 189, 182, 59, 255, 118, 194, 105, 101, 26, 236, 196, 230, 223,
		},
	},
	{
		rsg: reedSolomonGenerator{
			coefficients: []uint8{
				3, 2,
			},
		},
		data: &[]uint8{
			7, 6, 5, 4, 3, 2, 1,
		},
		expected: []uint8{
			123, 123,
		},
	},
	{
		rsg: reedSolomonGenerator{
			coefficients: []uint8{
				31, 198, 63, 147, 116,
			},
		},
		data: &[]uint8{
			0,
		},
		expected: []uint8{
			0, 0, 0, 0, 0,
		},
	},
	{
		rsg: reedSolomonGenerator{
			coefficients: []uint8{
				1,
			},
		},
		data: &[]uint8{
			1,
		},
		expected: []uint8{
			1,
		},
	},
}

func Test_getRemainder(test *testing.T) {
	for _, data := range getRemainder_TestData {
		actual := data.rsg.getRemainder(data.data)
		if len(actual) != len(data.expected) {
			test.Errorf(
				"reed_solomon_generator.Test_getRemainder:\n\tactual len of remainder array -> %d\n is not equal to\n\texpected len of remainder array -> %d",
				len(actual), len(data.expected),
			)
		}
		for i, coefficient := range actual {
			if coefficient != data.expected[i] {
				test.Errorf(
					"reed_solomon_generator.Test_getRemainder:\n\tactual remainder -> %d\n is not equal to\n\texpected remainder -> %d",
					coefficient, data.expected[i],
				)
			}
		}
	}
}

var multiply_TestData = []struct {
	x        uint8
	y        uint8
	expected uint8
}{
	{
		x: 1, y: 2, expected: 2,
	},
	{
		x: 3, y: 0, expected: 0,
	},
	{
		x: 10, y: 22, expected: 156,
	},
	{
		x: 3, y: 255, expected: 28,
	},
	{
		x: 0, y: 0, expected: 0,
	},
	{
		x: 1, y: 1, expected: 1,
	},
	{
		x: 255, y: 255, expected: 226,
	},
}

func Test_multiply(test *testing.T) {
	rsg := reedSolomonGenerator{}
	for _, data := range multiply_TestData {
		actual := rsg.multiply(data.x, data.y)
		if actual != data.expected {
			test.Errorf(
				"reed_solomon_generator.Test_multiply:\n\tactual product -> %d\n is not equal to\n\texpected product -> %d",
				actual, data.expected,
			)
		}
	}
}
