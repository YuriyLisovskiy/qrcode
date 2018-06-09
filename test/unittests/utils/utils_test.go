package utils

import (
	"testing"
	"github.com/YuriyLisovskiy/qrcode/src/utils"
	"github.com/YuriyLisovskiy/qrcode/test/unittests/utils/test_data"
)

func TestXOR(test *testing.T) {
	for _, td := range test_data.XORTestData {
		actual := utils.XOR(td.Expr1, td.Expr2)
		if actual != td.Expected {
			test.Errorf("utils.TestXOR: expected %t, actual %t", td.Expected, actual)
		}
	}
}
