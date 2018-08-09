package qr

import (
	"os"
	"fmt"
	"math"
	"image"
	"image/png"
	"image/draw"
	"image/color"

	"github.com/nfnt/resize"
)

// Represents an immutable square grid of black and white cells for a QR Code symbol, and
// provides static functions to create a QR Code from user-supplied textual or binary data.
// This class covers the QR Code model 2 specification, supporting all versions (sizes)
// from 1 to 40, all 4 error correction levels, and only 3 character encoding modes.
type Generator struct {
	// This QR Code symbol's version number, which is always between 1 and 40 (inclusive).
	version int

	// The width and height of this QR Code symbol, measured in modules.
	// Always equal to version &times; 4 + 17, in the range 21 to 177.
	size int

	// The error correction level used in this QR Code symbol.
	errorCorrectionLevel eccType

	// The mask pattern used in this QR Code symbol, in the range 0 to 7 (i.e. unsigned 3-bit integer).
	// Note that even if a constructor was called with automatic masking requested
	// (mask = -1), the resulting object will still have a mask value between 0 and 7.
	mask int

	// The modules of this QR Code symbol (false = white, true = black)
	modules [][]bool

	// Indicates function modules that are not subjected to masking
	isFunction [][]bool
}

// Returns a QR Code symbol representing the specified Unicode text string at the specified error correction level.
//
// As a conservative upper bound, this function is guaranteed to succeed for strings that have 2953 or fewer
// UTF-8 code units (not Unicode code points) if the low error correction level is used. The smallest possible
// QR Code version is automatically chosen for the output. The ECC level of the result may be higher than
// the ecl argument if it can be done without increasing the version.
func (qrc *Generator) EncodeText(text string) Generator {
	segments := makeSegments(text)
	return qrc.encodeSegments(&segments, eccLOW, 1, 40, -1, true)
}

// Returns a QR Code symbol representing the given binary data string at the given error correction level.
//
// This function always encodes using the binary segment mode, not any text mode. The maximum number of
// bytes allowed is 2953. The smallest possible QR Code version is automatically chosen for the output.
// The ECC level of the result may be higher than the ecl argument if it can be done without increasing the version.
func (qrc *Generator) EncodeBinary(data *[]uint8) Generator {
	return qrc.encodeSegments(&[]qrSegment{makeBytes(data)}, eccLOW, 1, 40, -1, true)
}

// Draws generated QR Code to the terminal window.
func (qrc *Generator) Draw(margin int) {
	for y := -margin; y < qrc.getSize()+margin; y++ {
		for x := -margin; x < qrc.getSize()+margin; x++ {
			if qrc.getModule(x, y) {
				print("  ")
			} else {
				print("\u2588\u2588")
			}
		}
		println()
	}
	println()
}

// Draws generated QR Code and save it to a file.
func (qrc *Generator) DrawImage(path string, margin, pictureSize uint) {
	if margin < 0 {
		panic("margin is negative")
	}
	size := qrc.getSize() + int(margin)*2
	if int(pictureSize) < size {
		panic(fmt.Sprintf("size of code is less than minimum size > %dx%d", size, size))
	}
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{R: 255, G: 255, B: 255, A: 255}}, image.ZP, draw.Src)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if qrc.getModule(x, y) {
				img.Set(x+int(margin), y+int(margin), color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	if int(pictureSize) > size {
		png.Encode(file, resize.Resize(pictureSize, pictureSize, img, resize.NearestNeighbor))
	} else {
		png.Encode(file, img)
	}
}

// Creates a new QR Code symbol with the given version number, error correction level, binary data array,
// and mask number. This is a cumbersome low-level constructor that should not be invoked directly by the user.
// To go one level up, see the encodeSegments() function.
func newQrCode(ver int, ecl eccType, dataCodewords []uint8, mask int) Generator {
	if ver < minVersion || ver > maxVersion || mask < -1 || mask > 7 {
		panic("value out of range")
	}
	newQrCode := Generator{version: ver, errorCorrectionLevel: ecl}
	if minVersion <= ver && ver <= maxVersion {
		newQrCode.size = ver*4 + 17
	} else {
		newQrCode.size = -1
	}
	newQrCode.modules = make([][]bool, newQrCode.size)
	for i := range newQrCode.modules {
		newQrCode.modules[i] = make([]bool, newQrCode.size)
	}
	newQrCode.isFunction = make([][]bool, newQrCode.size)
	for i := range newQrCode.isFunction {
		newQrCode.isFunction[i] = make([]bool, newQrCode.size)
	}
	newQrCode.drawFunctionPatterns()
	allCodewords := newQrCode.appendErrorCorrection(dataCodewords)
	newQrCode.drawCodewords(&allCodewords)
	newQrCode.mask = newQrCode.handleConstructorMasking(mask)
	return newQrCode
}

// Returns a value in the range 0 to 3 (unsigned 2-bit integer).
func (qrc *Generator) getFormatBits(ecl eccType) int {
	switch ecl {
	case eccLOW:
		return 1
	case eccMEDIUM:
		return 0
	case eccQUARTILE:
		return 3
	case eccHIGH:
		return 2
	default:
		panic("assertion error")
	}
}

// Returns a QR Code symbol representing the given data segments with the given encoding parameters.
//
// The smallest possible QR Code version within the given range is automatically chosen for the output.
// This function allows the user to create a custom sequence of segments that switches
// between modes (such as alphanumeric and binary) to encode text more efficiently.
// This function is considered to be lower level than simply encoding text or binary data.
func (qrc *Generator) encodeSegments(segs *[]qrSegment, ecl eccType, minVersion, maxVersion, mask int, boostEcl bool) Generator {
	if !(minVersion <= minVersion && minVersion <= maxVersion && maxVersion <= maxVersion) || mask < -1 || mask > 7 {
		panic("invalid value")
	}
	var version, dataUsedBits int
	for version = minVersion; ; version++ {
		dataCapacityBits := qrc.getNumDataCodewords(version, ecl) * 8
		dataUsedBits = getTotalBits(segs, version)
		if dataUsedBits != -1 && dataUsedBits <= dataCapacityBits {
			break
		}
		if version >= maxVersion {
			panic("data too long")
		}
	}
	if dataUsedBits == -1 {
		panic("assertion error")
	}
	for _, newEcl := range []eccType{eccMEDIUM, eccQUARTILE, eccHIGH} {
		if boostEcl && dataUsedBits <= qrc.getNumDataCodewords(version, newEcl)*8 {
			ecl = newEcl
		}
	}
	dataCapacityBits := uint(qrc.getNumDataCodewords(version, ecl) * 8)
	var bitBuf bitBuffer
	for _, seg := range *segs {
		bitBuf = bitBuf.AppendBits(uint32(seg.getMode().GetModeBits()), 4)
		bitBuf = bitBuf.AppendBits(uint32(seg.getNumChars()), seg.getMode().NumCharCountBits(version))
		bitBuf = append(bitBuf, *seg.getData()...)
	}
	bitBuf = bitBuf.AppendBits(0, int(math.Min(float64(4), float64(dataCapacityBits-uint(len(bitBuf))))))
	bitBuf = bitBuf.AppendBits(0, (8-len(bitBuf)%8)%8)
	for padByte := 0xEC; uint(len(bitBuf)) < dataCapacityBits; padByte ^= 0xEC ^ 0x11 {
		bitBuf = bitBuf.AppendBits(uint32(padByte), 8)
	}
	if len(bitBuf)%8 != 0 {
		panic("assertion error")
	}
	return newQrCode(version, ecl, bitBuf.getBytes(), mask)
}

// Instance method.
//
// Returns QR Code generator's version.
func (qrc *Generator) getVersion() int {
	return qrc.version
}

// Instance method.
//
// Returns the size of QR Code.
func (qrc *Generator) getSize() int {
	return qrc.size
}

// Instance method.
//
// Returns an error correction level of QR Code.
func (qrc *Generator) getErrorCorrectionLevel() eccType {
	return qrc.errorCorrectionLevel
}

// Instance method.
//
// Returns QR Code generator's mask.
func (qrc *Generator) getMask() int {
	return qrc.mask
}

// Returns the color of the module (pixel) at the given coordinates, which is either
// false for white or true for black. The top left corner has the coordinates (x=0, y=0).
//
// If the given coordinates are out of bounds, then false (white) is returned.
func (qrc *Generator) getModule(x, y int) bool {
	return 0 <= x && x < qrc.size && 0 <= y && y < qrc.size && qrc.module(x, y)
}

// Returns the color of the module at the given coordinates, which must be in range.
func (qrc *Generator) module(x, y int) bool {
	return qrc.modules[y][x]
}

// Returns true iff the i'th bit of x is set to 1.
func (qrc *Generator) getBit(x int, i uint) bool {
	return ((x >> i) & 1) != 0
}

// Draws function patterns.
//
// Helper method for constructor: drawing function modules.
func (qrc *Generator) drawFunctionPatterns() {
	for i := 0; i < qrc.size; i++ {
		qrc.setFunctionModule(6, i, i%2 == 0)
		qrc.setFunctionModule(i, 6, i%2 == 0)
	}
	qrc.drawFinderPattern(3, 3)
	qrc.drawFinderPattern(qrc.size-4, 3)
	qrc.drawFinderPattern(3, qrc.size-4)
	alignPatPos := []int(qrc.getAlignmentPatternPositions(qrc.version))
	numAlign := len(alignPatPos)

	for i := 0; i < numAlign; i++ {
		for j := 0; j < numAlign; j++ {
			if (i == 0 && j == 0) || (i == 0 && j == numAlign-1) || (i == numAlign-1 && j == 0) {
				continue
			} else {
				qrc.drawAlignmentPattern(alignPatPos[i], alignPatPos[j])
			}
		}
	}
	qrc.drawFormatBits(0)
	qrc.drawVersion()
}

// Draws two copies of the format bits (with its own error correction code)
// based on the given mask and this object's error correction level field.
//
// Helper method for constructor: drawing function modules.
func (qrc *Generator) drawFormatBits(mask int) {
	data := int(qrc.getFormatBits(qrc.errorCorrectionLevel)<<3 | mask)
	rem := int(data)
	for i := 0; i < 10; i++ {
		rem = (rem << 1) ^ ((rem >> 9) * 0x537)
	}
	data = data<<10 | rem
	data ^= 0x5412
	if data>>15 != 0 {
		panic("assertion error")
	}
	for i := 0; i <= 5; i++ {
		qrc.setFunctionModule(8, i, qrc.getBit(int(data), uint(i)))
	}
	qrc.setFunctionModule(8, 7, qrc.getBit(int(data), 6))
	qrc.setFunctionModule(8, 8, qrc.getBit(int(data), 7))
	qrc.setFunctionModule(7, 8, qrc.getBit(int(data), 8))
	for i := 9; i < 15; i++ {
		qrc.setFunctionModule(14-i, 8, qrc.getBit(int(data), uint(i)))
	}
	for i := 0; i <= 7; i++ {
		qrc.setFunctionModule(qrc.size-1-i, 8, qrc.getBit(int(data), uint(i)))
	}
	for i := 8; i < 15; i++ {
		qrc.setFunctionModule(8, qrc.size-15+i, qrc.getBit(int(data), uint(i)))
	}
	qrc.setFunctionModule(8, qrc.size-8, true)
}

// Draws two copies of the version bits (with its own error correction code),
// based on this object's version field (which only has an effect for 7 <= version <= 40).
//
// Helper method for constructor: drawing function modules.
func (qrc *Generator) drawVersion() {
	if qrc.version < 7 {
		return
	}
	rem := int(qrc.version)
	for i := 0; i < 12; i++ {
		rem = (rem << 1) ^ ((rem >> 11) * 0x1F25)
	}
	data := int(qrc.version)<<12 | int(rem)
	if data>>18 != 0 {
		panic("assertion error")
	}
	for i := 0; i < 18; i++ {
		bit := qrc.getBit(int(data), uint(i))
		a, b := int(qrc.size-11+i%3), int(i/3)
		qrc.setFunctionModule(a, b, bit)
		qrc.setFunctionModule(b, a, bit)
	}
}

// Draws a 9*9 finder pattern including the border separator, with the center module at (x, y).
//
// Helper method for constructor: drawing function modules.
func (qrc *Generator) drawFinderPattern(x, y int) {
	for i := -4; i <= 4; i++ {
		for j := -4; j <= 4; j++ {
			dist := int(math.Max(math.Abs(float64(i)), math.Abs(float64(j))))
			xx, yy := int(x+j), int(y+i)
			if 0 <= xx && xx < qrc.size && 0 <= yy && yy < qrc.size {
				qrc.setFunctionModule(xx, yy, dist != 2 && dist != 4)
			}
		}
	}
}

// Draws a 5*5 alignment pattern, with the center module at (x, y).
//
// Helper method for constructor: drawing function modules.
func (qrc *Generator) drawAlignmentPattern(x, y int) {
	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			qrc.setFunctionModule(x+j, y+i, math.Max(math.Abs(float64(i)), math.Abs(float64(j))) != 1)
		}
	}
}

// Sets the color of a module and marks it as a function module.
// Only used by the constructor. Coordinates must be in range.
//
// Helper method for constructor: drawing function modules.
func (qrc *Generator) setFunctionModule(x, y int, isBlack bool) {
	qrc.modules[y][x] = isBlack
	qrc.isFunction[y][x] = true
}

// Returns a new byte string representing the given data with the appropriate error correction
// codewords appended to it, based on this object's version and error correction level.
//
// Helper method for constructor: codewords and masking.
func (qrc *Generator) appendErrorCorrection(data []uint8) []uint8 {
	if len(data) != qrc.getNumDataCodewords(qrc.version, qrc.errorCorrectionLevel) {
		panic("invalid argument")
	}
	numBlocks := int(numErrorCorrectionBlocks[int(qrc.errorCorrectionLevel)][qrc.version])
	blockEccLen := int(eccCodewordsPerBlock[int(qrc.errorCorrectionLevel)][qrc.version])
	rawCodewords := int(qrc.getNumRawDataModules(qrc.version) / 8)
	numShortBlocks := int(numBlocks - rawCodewords%numBlocks)
	shortBlockLen := int(rawCodewords / numBlocks)
	var blocks [][]uint8
	rs := newReedSolomonGenerator(blockEccLen)
	i, start := 0, 0
	for ; i < numBlocks; i++ {
		c := 1
		if i < numShortBlocks {
			c = 0
		}
		end := start + shortBlockLen - blockEccLen + c
		dat := make([]uint8, end-start)
		copy(dat, []uint8(data[start:end]))
		start += len(dat)
		ecc := []uint8(rs.GetRemainder(&dat))
		if i < numShortBlocks {
			dat = append(dat, 0)
		}
		dat = append(dat, ecc...)
		blocks = append(blocks, dat)
	}

	var result []uint8
	for i := 0; i < len(blocks[0]); i++ {
		for j := 0; j < len(blocks); j++ {
			if i != shortBlockLen-blockEccLen || j >= numShortBlocks {
				result = append(result, blocks[j][i])
			}
		}
	}
	if len(result) != rawCodewords {
		panic("assertion error")
	}
	return result
}

// Draws the given sequence of 8-bit codewords (data and error correction) onto the entire
// data area of this QR Code symbol. Function modules need to be marked off before this is called.
//
// Helper method for constructor: codewords and masking.
func (qrc *Generator) drawCodewords(data *[]uint8) {
	if len(*data) != qrc.getNumRawDataModules(qrc.version)/8 {
		panic("invalid argument")
	}
	i := uint(0)
	for right := qrc.size - 1; right >= 1; right -= 2 {
		if right == 6 {
			right = 5
		}
		for vert := 0; vert < qrc.size; vert++ {
			for j := 0; j < 2; j++ {
				x := int(right - j)
				y := vert
				if ((right + 1) & 2) == 0 {
					y = qrc.size - 1 - vert
				}
				if !qrc.isFunction[y][x] && i < uint(len(*data)*8) {
					qrc.modules[y][x] = qrc.getBit(int((*data)[i>>3]), uint(7-(int(i & 7))))
					i++
				}
			}
		}
	}
	if i != uint(len(*data)*8) {
		panic("assertion error")
	}
}

// XORs the data modules in this QR Code with the given mask pattern. Due to XOR's mathematical
// properties, calling applyMask(m) twice with the same value is equivalent to no change at all.
// This means it is possible to apply a mask, undo it, and try another mask. Note that a final
// well-formed QR Code symbol needs exactly one mask applied (not zero, not two, etc.).
//
// Helper method for constructor: codewords and masking.
func (qrc *Generator) applyMask(mask int) {
	if mask < 0 || mask > 7 {
		panic("mask value out of range")
	}
	for y := 0; y < qrc.size; y++ {
		for x := 0; x < qrc.size; x++ {
			var invert bool
			switch mask {
			case 0:
				invert = (x+y)%2 == 0
			case 1:
				invert = y%2 == 0
			case 2:
				invert = x%3 == 0
			case 3:
				invert = (x+y)%3 == 0
			case 4:
				invert = (x/3+y/2)%2 == 0
			case 5:
				invert = x*y%2+x*y%3 == 0
			case 6:
				invert = (x*y%2+x*y%3)%2 == 0
			case 7:
				invert = ((x+y)%2+x*y%3)%2 == 0
			default:
				panic("assertion error")
			}
			qrc.modules[y][x] = xor(qrc.modules[y][x], invert && !qrc.isFunction[y][x])
		}
	}
}

// A messy helper function for the constructors. This QR Code must be in an unmasked state when this
// method is called. The given argument is the requested mask, which is -1 for auto or 0 to 7 for fixed.
// This method applies and returns the actual mask chosen, from 0 to 7.
//
// Helper method for constructor: codewords and masking.
func (qrc *Generator) handleConstructorMasking(mask int) int {
	if mask == -1 {
		minPenalty := math.MaxInt32
		for i := 0; i < 8; i++ {
			qrc.drawFormatBits(i)
			qrc.applyMask(i)
			penalty := qrc.getPenaltyScore()
			if penalty < minPenalty {
				mask = i
				minPenalty = penalty
			}
			qrc.applyMask(i)
		}
	}
	if mask < 0 || mask > 7 {
		panic("assertion error")
	}
	qrc.drawFormatBits(mask)
	qrc.applyMask(mask)
	return mask
}

// Calculates and returns the penalty score based on state of this QR Code's current modules.
// This is used by the automatic mask choice algorithm to find the mask pattern that yields the lowest score.
//
// Helper method for constructor: codewords and masking.
func (qrc *Generator) getPenaltyScore() int {
	result := 0
	for y := 0; y < qrc.size; y++ {
		colorX := false
		runX := -1
		for x := 0; x < qrc.size; x++ {
			if x == 0 || qrc.module(x, y) != colorX {
				colorX = qrc.module(x, y)
				runX = 1
			} else {
				runX++
				if runX == 5 {
					result += penaltyN1
				} else if runX > 5 {
					result++
				}
			}
		}
	}
	for x := 0; x < qrc.size; x++ {
		colorY := false
		runY := -1
		for y := 0; y < qrc.size; y++ {
			if y == 0 || qrc.module(x, y) != colorY {
				colorY = qrc.module(x, y)
				runY = 1
			} else {
				runY++
				if runY == 5 {
					result += penaltyN1
				} else if runY > 5 {
					result++
				}
			}
		}
	}
	for y := 0; y < qrc.size-1; y++ {
		for x := 0; x < qrc.size-1; x++ {
			colour := qrc.module(x, y)
			if colour == qrc.module(x+1, y) && colour == qrc.module(x, y+1) && colour == qrc.module(x+1, y+1) {
				result += penaltyN2
			}
		}
	}
	for y := 0; y < qrc.size; y++ {
		bits := 0
		for x := 0; x < qrc.size; x++ {
			t := 0
			if qrc.module(x, y) {
				t = 1
			}
			bits = ((bits << 1) & 0x7FF) | t
			if x >= 10 && (bits == 0x05D || bits == 0x5D0) {
				result += penaltyN3
			}
		}
	}
	for x := 0; x < qrc.size; x++ {
		bits := 0
		for y := 0; y < qrc.size; y++ {
			t := 0
			if qrc.module(x, y) {
				t = 1
			}
			bits = ((bits << 1) & 0x7FF) | t
			if y >= 10 && (bits == 0x05D || bits == 0x5D0) {
				result += penaltyN3
			}
		}
	}
	black := 0
	for _, row := range qrc.modules {
		for _, colour := range row {
			if colour {
				black++
			}
		}
	}
	total := qrc.size * qrc.size
	for k := 0; black*int(20) < (int(9)-k)*total || black*int(20) > (int(11)+k)*total; k++ {
		result += penaltyN4
	}
	return result
}

// Returns a set of positions of the alignment patterns in ascending order. These positions are
// used on both the x and y axes. Each value in the resulting array is in the range [0, 177).
// This stateless pure function could be implemented as table of 40 variable-length lists of unsigned bytes.
//
// Helper function.
func (qrc *Generator) getAlignmentPatternPositions(ver int) []int {
	if ver < minVersion || ver > maxVersion {
		panic("Version number out of range")
	} else if ver == 1 {
		return []int{}
	} else {
		numAlign := int(ver/7 + 2)
		var step int
		if ver != 32 {
			step = (ver*4 + numAlign*2 + 1) / (2*numAlign - 2) * 2
		} else {
			step = 26
		}
		var result []int
		pos := int(ver*4 + 10)
		for i := 0; i < numAlign-1; i++ {
			result = append([]int{pos}, result...)
			pos -= step
		}
		result = append([]int{6}, result...)
		return result
	}
}

// Returns the number of data bits that can be stored in a QR Code of the given version number, after
// all function modules are excluded. This includes remainder bits, so it might not be a multiple of 8.
// The result is in the range [208, 29648]. This could be implemented as a 40-entry lookup table.
//
// Helper function.
func (qrc *Generator) getNumRawDataModules(ver int) int {
	if ver < minVersion || ver > maxVersion {
		panic("version number out of range")
	}
	result := int((16*ver+128)*ver + 64)
	if ver >= 2 {
		numAlign := int(ver/7 + 2)
		result -= (25*numAlign-10)*numAlign - 55
		if ver >= 7 {
			result -= 18 * 2
		}
	}
	return result
}

// Returns the number of 8-bit data (i.e. not error correction) codewords contained in any
// QR Code of the given version number and error correction level, with remainder bits discarded.
// This stateless pure function could be implemented as a (40*4)-cell lookup table.
//
// Helper function.
func (qrc *Generator) getNumDataCodewords(ver int, ecl eccType) int {
	if ver < minVersion || ver > maxVersion {
		panic("version number out of range")
	}
	return qrc.getNumRawDataModules(ver)/8 - eccCodewordsPerBlock[int(ecl)][ver]*numErrorCorrectionBlocks[int(ecl)][ver]
}
