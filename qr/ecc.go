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
