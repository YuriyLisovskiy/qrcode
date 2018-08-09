//  Copyright (c) 2018 Yuriy Lisovskiy
//  Distributed under the Apache License Version 2.0,
//  see the accompanying file LICENSE or https://opensource.org/licenses/Apache-2.0

package qr

import "fmt"

// Returns XOR result between two boolean values.
func xor(x, y bool) bool {
	return (x && !y) || (!x && y)
}

// Compose Generator error message
func generatorErr(method, msg string) string {
	return fmt.Sprintf("package qr: Generator.%s: %s", method, msg)
}

// Compose reedSolomonGenerator error message
func rsgErr(method, msg string) string {
	return fmt.Sprintf("package qr: reedSolomonGenerator.%s: %s", method, msg)
}

// Compose qrSegment error message
func qrSegmentErr(method, msg string) string {
	return fmt.Sprintf("package qr: qrSegment.%s: %s", method, msg)
}
