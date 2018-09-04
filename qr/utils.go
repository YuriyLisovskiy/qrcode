//  Copyright (c) 2018 Yuriy Lisovskiy
//  Distributed under the Apache License Version 2.0,
//  see the accompanying file LICENSE or https://opensource.org/licenses/Apache-2.0

package qr // import "github.com/YuriyLisovskiy/qrcode/qr"

// Returns XOR result between two boolean values.
func xor(x, y bool) bool {
	return (x && !y) || (!x && y)
}
