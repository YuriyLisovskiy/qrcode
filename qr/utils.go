package qr

import "fmt"

// Returns XOR result between two boolean values.
func xor(x, y bool) bool {
	return (x && !y) || (!x && y)
}

func generatorErr(method, msg string) string {
	return fmt.Sprintf("package qr: Generator.%s: %s", method, msg)
}

func rsgErr(method, msg string) string {
	return fmt.Sprintf("package qr: reedSolomonGenerator.%s: %s", method, msg)
}

func qrSegmentErr(method, msg string) string {
	return fmt.Sprintf("package qr: qrSegment.%s: %s", method, msg)
}
