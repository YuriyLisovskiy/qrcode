package utils

import (
	"testing"
	"github.com/YuriyLisovskiy/qrcode/src/utils"
)

func TestXOR(test *testing.T) {
	for _, td := range XORTestData {
		actual := utils.XOR(td.Expr1, td.Expr2)
		if actual != td.Expected {
			test.Errorf("utils.TestXOR: expected %t, actual %t", td.Expected, actual)
		}
	}
}
