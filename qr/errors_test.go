package qr

import (
	"fmt"
	"errors"
	"testing"
)

var generatorErr_TestData = []struct {
	errMsg   string
	method   string
	expected error
}{
	{
		errMsg: "some error message 1",
		method: "someMethod1",
		expected: errors.New("package qr: Generator.someMethod1: some error message 1"),
	},
	{
		errMsg: fmt.Sprintf("some error message with data: %d, %s, %t", 10, "Hello", false),
		method: "genMethod",
		expected: errors.New("package qr: Generator.genMethod: some error message with data: 10, Hello, false"),
	},
}

func Test_generatorErr(test *testing.T) {
	for _, data := range generatorErr_TestData {
		actual := generatorErr(data.method, data.errMsg)
		if actual.Error() != data.expected.Error() {
			test.Errorf(
				"utils.Test_generatorErr:\n\tactual Generator err message -> %s\n is not equal to\n\texpected Generator err message -> %s",
				actual, data.expected,
			)
		}
	}
}

var rsgErr_TestData = []struct {
	errMsg   string
	method   string
	expected error
}{
	{
		errMsg: "some error message 1",
		method: "someMethod1",
		expected: errors.New("package qr: reedSolomonGenerator.someMethod1: some error message 1"),
	},
	{
		errMsg: fmt.Sprintf("some error message with data: %d, %s, %t", 10, "Hello", false),
		method: "genMethod",
		expected: errors.New("package qr: reedSolomonGenerator.genMethod: some error message with data: 10, Hello, false"),
	},
}

func Test_rsgErr(test *testing.T) {
	for _, data := range rsgErr_TestData {
		actual := rsgErr(data.method, data.errMsg)
		if actual.Error() != data.expected.Error() {
			test.Errorf(
				"utils.Test_rsgErr:\n\tactual reedSolomonGenerator err message -> %s\n is not equal to\n\texpected reedSolomonGenerator err message -> %s",
				actual, data.expected,
			)
		}
	}
}

var qrSegmentErr_TestData = []struct {
	errMsg   string
	method   string
	expected error
}{
	{
		errMsg: "some error message 1",
		method: "someMethod1",
		expected: errors.New("package qr: qrSegment.someMethod1: some error message 1"),
	},
	{
		errMsg: fmt.Sprintf("some error message with data: %d, %s, %t", 10, "Hello", false),
		method: "genMethod",
		expected: errors.New("package qr: qrSegment.genMethod: some error message with data: 10, Hello, false"),
	},
}

func Test_qrSegmentErr(test *testing.T) {
	for _, data := range qrSegmentErr_TestData {
		actual := qrSegmentErr(data.method, data.errMsg)
		if actual.Error() != data.expected.Error() {
			test.Errorf(
				"utils.Test_qrSegmentErr:\n\tactual qrSegment err message -> %s\n is not equal to\n\texpected qrSegment err message -> %s",
				actual, data.expected,
			)
		}
	}
}
