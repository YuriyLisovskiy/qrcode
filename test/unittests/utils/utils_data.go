package utils

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
