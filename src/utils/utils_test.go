package utils

import "testing"

var (
	XORTestData = []struct{
		Expr1 bool
		Expr2 bool
		Expected bool
	}{
		{false, false, false},
		{false, true, true},
		{true, false, true},
		{false, false, false},
	}
)

func TestXOR(test *testing.T) {
	for _, td := range XORTestData {
		actual := XOR(td.Expr1, td.Expr2)
		if actual != td.Expected {
			test.Errorf("utils.TestXOR: expected %t, actual %t", td.Expected, actual)
		}
	}
}
