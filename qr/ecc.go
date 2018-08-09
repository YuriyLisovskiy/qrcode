//  Copyright (c) 2018 Yuriy Lisovskiy
//  Distributed under the Apache License Version 2.0,
//  see the accompanying file LICENSE or https://opensource.org/licenses/Apache-2.0

package qr

// Represents the error correction level used in a QR Code symbol.
type eccType uint

// Constants declared in ascending order of error protection.
const (
	eccLOW      eccType = iota
	eccMEDIUM
	eccQUARTILE
	eccHIGH
)
