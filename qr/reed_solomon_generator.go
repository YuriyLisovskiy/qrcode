//  Copyright (c) 2018 Yuriy Lisovskiy
//  Distributed under the Apache License Version 2.0,
//  see the accompanying file LICENSE or https://opensource.org/licenses/Apache-2.0

package qr // import "github.com/YuriyLisovskiy/qrcode/qr"

// Computes the Reed-Solomon error correction codewords for a sequence of data codewords
// at a given degree. Objects are immutable, and the state only depends on the degree.
// This class exists because each data block in a QR Code shares the same the divisor polynomial.
type reedSolomonGenerator struct {
	// Coefficients of the divisor polynomial, stored from highest to lowest power, excluding the leading term which
	// is always 1. For example the polynomial x^3 + 255x^2 + 8x + 93 is stored as the uint8 array {255, 8, 93}.
	coefficients []uint8
}

// Creates a Reed-Solomon ECC generator for the given degree. This could be implemented
// as a lookup table over all possible parameter values, instead of as an algorithm.
func newReedSolomonGenerator(degree int) (reedSolomonGenerator, error) {
	newRSG := reedSolomonGenerator{}
	if degree < 1 || degree > 255 {
		return reedSolomonGenerator{}, rsgErr("newReedSolomonGenerator", "degree out of range")
	}
	newRSG.coefficients = make([]uint8, degree)
	newRSG.coefficients[degree-1] = 1
	root := uint8(1)
	for i := 0; i < degree; i++ {
		for j := 0; j < len(newRSG.coefficients); j++ {
			newRSG.coefficients[j] = newRSG.multiply(newRSG.coefficients[j], root)
			if j+1 < len(newRSG.coefficients) {
				newRSG.coefficients[j] ^= newRSG.coefficients[j+1]
			}
		}
		root = newRSG.multiply(root, 0x02)
	}
	return newRSG, nil
}

// Computes and returns the Reed-Solomon error correction codewords for the given
// sequence of data codewords. The returned object is always a new byte array.
// This method does not alter this object's state (because it is immutable).
func (rsg reedSolomonGenerator) getRemainder(data *[]uint8) []uint8 {
	result := make([]uint8, len(rsg.coefficients))
	for _, b := range *data {
		factor := b ^ result[0]
		result = result[1:]
		result = append(result, 0)
		for j := 0; j < len(result); j++ {
			result[j] ^= rsg.multiply(rsg.coefficients[j], factor)
		}
	}
	return result
}

// Returns the product of the two given field elements modulo GF(2^8/0x11D).
//
// All inputs are valid. This could be implemented as a 256*256 lookup table.
func (reedSolomonGenerator) multiply(x, y uint8) uint8 {
	z := 0
	for i := 7; i >= 0; i-- {
		z = (z << 1) ^ ((z >> 7) * 0x11D)
		z = int(uint8(z) ^ ((y>>uint8(i))&1)*x)
	}
	return uint8(z)
}
