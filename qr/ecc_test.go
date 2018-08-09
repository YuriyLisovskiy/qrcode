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
			test.Errorf("ecc.Test_EccType:\n\tactual -> %d\n is not equal to\n\texpected -> %d", data.actual, data.expected)
		}
	}
}
