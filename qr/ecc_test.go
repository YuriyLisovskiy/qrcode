package qr

import "testing"

var EccType_TestData = []struct{
	actual eccType
	expected eccType
} {
	{
		actual: eccLOW,
		expected: 0,
	},
	{
		actual: eccMEDIUM,
		expected: 1,
	},
	{
		actual: eccQUARTILE,
		expected: 2,
	},
	{
		actual: eccHIGH,
		expected: 3,
	},
}

func Test_EccType(test *testing.T) {
	for _, data := range EccType_TestData {
		if data.actual != data.expected {
			test.Errorf("qr.Test_EccType: actual -> %d is not equal to expected -> %d", data.actual, data.expected)
		}
	}
}
