//  Copyright (c) 2018 Yuriy Lisovskiy
//  Distributed under the Apache License Version 2.0,
//  see the accompanying file LICENSE or https://opensource.org/licenses/Apache-2.0

package qr // import "github.com/YuriyLisovskiy/qrcode/qr"

import (
	"errors"
	"fmt"
)

var (
	// bitBuffer errors.
	ErrBitBufValOutOfRange = errors.New("package qr: bitBuffer.appendBits: value out of range")

	// modeType errors.
	ErrModeTypeVerNumOutOfRange = errors.New("package qr: modeType.numCharCountBits: version number out of range")
)

// Compose Generator error message
func generatorErr(method, msg string) error {
	return errors.New(fmt.Sprintf("package qr: Generator.%s: %s", method, msg))
}

// Compose reedSolomonGenerator error message
func rsgErr(method, msg string) error {
	return errors.New(fmt.Sprintf("package qr: reedSolomonGenerator.%s: %s", method, msg))
}

// Compose qrSegment error message
func qrSegmentErr(method, msg string) error {
	return errors.New(fmt.Sprintf("package qr: qrSegment.%s: %s", method, msg))
}
