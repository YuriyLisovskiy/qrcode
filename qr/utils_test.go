//  Copyright (c) 2018 Yuriy Lisovskiy
//  Distributed under the Apache License Version 2.0,
//  see the accompanying file LICENSE or https://opensource.org/licenses/Apache-2.0

package qr

import "testing"

var xor_TestData = []struct {
		expr1    bool
		expr2    bool
		expected bool
	}{
		{false, false, false},
		{false, true, true},
		{true, false, true},
		{false, false, false},
	}

func Test_xor(test *testing.T) {
	for _, data := range xor_TestData {
		actual := xor(data.expr1, data.expr2)
		if actual != data.expected {
			test.Errorf(
				"utils.Test_xor:\n\tactual xor result -> %t\n is not equal to\n\texpected xor result len -> %t",
				actual, data.expected,
			)
		}
	}
}


