package qrcode

import (
	"github.com/YuriyLisovskiy/qrcode/src/utils"
	"github.com/YuriyLisovskiy/qrcode/src/vars"
	"math"
)

type QrCode struct {
	Version              int
	Size                 int
	ErrorCorrectionLevel utils.Ecc
	Mask                 int
	Modules              [][]bool
	IsFunction           [][]bool
}

func NewQrCode(ver int, ecl utils.Ecc, dataCodewords []uint8, mask int) QrCode {
	if ver < vars.MinVersion || ver > vars.MaxVersion || mask < -1 || mask > 7 {
		panic("value out of range")
	}
	newQrCode := QrCode{Version: ver, ErrorCorrectionLevel: ecl}
	if vars.MinVersion <= ver && ver <= vars.MaxVersion {
		newQrCode.Size = ver * 4 + 17
	} else {
		newQrCode.Size = -1
	}
	newQrCode.Modules = make([][]bool, newQrCode.Size)
	newQrCode.IsFunction = make([][]bool, newQrCode.Size)
	newQrCode.DrawFunctionPatterns()
	allCodewords := newQrCode.AppendErrorCorrection(&dataCodewords)
	newQrCode.DrawCodewords(&allCodewords)
	newQrCode.Mask = newQrCode.HandleConstructorMasking(mask)
	return newQrCode
}

func (qrc *QrCode) GetFormatBits(ecl utils.Ecc) int {
	switch ecl {
	case utils.EccLOW:
		return 1
	case utils.EccMEDIUM:
		return 0
	case utils.EccQUARTILE:
		return 3
	case utils.EccHIGH:
		return 2
	default:
		panic("assertion error")
	}
}

func (qrc *QrCode) EncodeText(text string, ecl utils.Ecc) QrCode {
	segs := MakeSegments(text)
	return qrc.EncodeSegments(&segs, ecl, 1, 40, -1, true)
}

func (qrc *QrCode) EncodeBinary(data *[]uint8, ecl utils.Ecc) QrCode {
	segs := []QrSegment{MakeBytes(data)}
	return qrc.EncodeSegments(&segs, ecl, 1, 40, -1, true)
}

func (qrc *QrCode) EncodeSegments(segs *[]QrSegment, ecl utils.Ecc, minVersion, maxVersion, mask int, boostEcl bool) QrCode {
	if !(vars.MinVersion <= minVersion && minVersion <= maxVersion && maxVersion <= vars.MaxVersion) || mask < -1 || mask > 7 {
		panic("invalid value")
	}
	var version, dataUsedBits int
	for version = minVersion; ; version++ {
		dataCapacityBits := qrc.GetNumDataCodewords(version, ecl) * 8
		dataUsedBits = GetTotalBits(segs, version)
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
	for _, newEcl := range []utils.Ecc{utils.EccMEDIUM, utils.EccQUARTILE, utils.EccHIGH} {
		if boostEcl && dataUsedBits <= qrc.GetNumDataCodewords(version, newEcl) * 8 {
			ecl = newEcl
		}
	}
	dataCapacityBits := qrc.GetNumDataCodewords(version, ecl) * 8
	var bitBuf utils.BitBuffer
	for _, seg := range *segs {
		bitBuf = bitBuf.AppendBits(uint32(seg.GetMode().GetModeBits()), 4)
		bitBuf = bitBuf.AppendBits(uint32(seg.GetNumChars()), seg.GetMode().NumCharCountBits(version))
		bitBuf = append(bitBuf, *seg.GetData()...)
	}
	bitBuf = bitBuf.AppendBits(0, int(math.Min(float64(4), float64(dataCapacityBits - len(bitBuf)))))
	bitBuf = bitBuf.AppendBits(0, (8 - len(bitBuf) % 8) % 8)
	for padByte := 0xEC; len(bitBuf) < dataCapacityBits; padByte ^= 0xEC ^ 0x11 {
		bitBuf = bitBuf.AppendBits(uint32(padByte), 8)
	}
	if len(bitBuf) % 8 != 0 {
		panic("assertion error")
	}
	return NewQrCode(version, ecl, bitBuf.GetBytes(), mask)
}

func (qrc *QrCode) GetVersion() int {
	return qrc.Version
}
func (qrc *QrCode) GetSize() int {
	return qrc.Size
}
func (qrc *QrCode) GetErrorCorrectionLevel() utils.Ecc {
	return qrc.ErrorCorrectionLevel
}
func (qrc *QrCode) GetMask() int {
	return qrc.Mask
}
func (qrc *QrCode) GetModule(x, y int) bool {
	return 0 <= x && x < qrc.Size && 0 <= y && y < qrc.Size && qrc.Module(x, y)
}
func (qrc *QrCode) Module(x, y int) bool {
	return qrc.Modules[x][y]
}
func (qrc *QrCode) GetBit(x int64, i int) bool {
	return ((x >> uint(i)) & 1) != 0
}

func (qrc *QrCode) DrawFunctionPatterns() {
	for i := 0; i < qrc.Size; i++ {
		qrc.SetFunctionModule(6, i, i % 2 == 0)
		qrc.SetFunctionModule(i, 6, i % 2 == 0)
	}
	qrc.DrawFinderPattern(3, 3)
	qrc.DrawFinderPattern(qrc.Size - 4, 3)
	qrc.DrawFinderPattern(3, qrc.Size - 4)
	alignPatPos := qrc.GetAlignmentPatternPositions(qrc.Version)
	numAlign := len(alignPatPos)

	for i := 0; i < numAlign; i++ {
		for j := 0; j < numAlign; j++ {
			if (i == 0 && j == 0) || (i == 0 && j == numAlign - 1) || (i == numAlign - 1 && j == 0) {
				continue
			} else {
				qrc.DrawAlignmentPattern(alignPatPos[i], alignPatPos[j])
			}
		}
	}
	qrc.DrawFormatBits(0)
	qrc.DrawVersion()
}

func (qrc *QrCode) DrawFormatBits(mask int) {
	data := qrc.GetFormatBits(qrc.ErrorCorrectionLevel) << 3 | mask
	rem := data
	for i := 0; i < 10; i++ {
		rem = (rem << 1) ^ ((rem >> 9) * 0x537)
	}
	data = data << 10 | rem
	data ^= 0x5412
	if data >> 15 != 0 {
		panic("assertion error")
	}
	for i := 0; i <= 5; i++ {
		qrc.SetFunctionModule(8, i, qrc.GetBit(int64(data), i))
	}
	qrc.SetFunctionModule(8, 7, qrc.GetBit(int64(data), 6))
	qrc.SetFunctionModule(8, 8, qrc.GetBit(int64(data), 7))
	qrc.SetFunctionModule(7, 8, qrc.GetBit(int64(data), 8))
	for i := 9; i < 15; i++ {
		qrc.SetFunctionModule(14 - i, 8, qrc.GetBit(int64(data), i))
	}
	for i := 0; i <= 7; i++ {
		qrc.SetFunctionModule(qrc.Size - 1 - i, 8, qrc.GetBit(int64(data), i))
	}
	for i := 8; i < 15; i++ {
		qrc.SetFunctionModule(8, qrc.Size - 15 + i, qrc.GetBit(int64(data), i))
	}
	qrc.SetFunctionModule(8, qrc.Size - 8, true)
}

func (qrc *QrCode) DrawVersion() {
	if qrc.Version < 7 {
		return
	}
	rem := qrc.Version
	for i := 0; i < 12; i++ {
		rem = (rem << 1) ^ ((rem >> 11) * 0x1F25)
	}
	data := int64(qrc.Version) << 12 | int64(rem)
	if data >> 18 != 0 {
		panic("assertion error")
	}
	for i := 0; i < 18; i++ {
		bit := qrc.GetBit(data, i)
		a, b := qrc.Size - 11 + i % 3, i / 3
		qrc.SetFunctionModule(a, b, bit)
		qrc.SetFunctionModule(b, a, bit)
	}
}

func (qrc *QrCode) DrawFinderPattern(x, y int) {
	for i := -4; i <= 4; i++ {
		for j := -4; j <= 4; j++ {
			dist := math.Max(math.Abs(float64(i)), math.Abs(float64(j)))
			xx, yy := x + j, y + i
			if 0 <= xx && xx < qrc.Size && 0 <= yy && yy < qrc.Size {
				qrc.SetFunctionModule(xx, yy, dist != 2 && dist != 4)
			}
		}
	}
}

func (qrc *QrCode) DrawAlignmentPattern(x, y int) {
	for i := -2; i <= 2; i++ {
		for j := -2; j <= 2; j++ {
			qrc.SetFunctionModule(x + j, y + i, math.Max(math.Abs(float64(i)), math.Abs(float64(j))) != 1)
		}
	}
}

func (qrc *QrCode) SetFunctionModule(x, y int, isBlack bool) {
	qrc.Modules[y][x] = isBlack
	qrc.IsFunction[y][x] = true
}

func (qrc *QrCode) AppendErrorCorrection(data *[]uint8) []uint8 {
	if len(*data) != qrc.GetNumDataCodewords(qrc.Version, qrc.ErrorCorrectionLevel) {
		panic("invalid argument")
	}
	numBlocks := vars.NumErrorCorrectionBlocks[qrc.ErrorCorrectionLevel][qrc.Version]
	blockEccLen := vars.EccCodewordsPerBlock[qrc.ErrorCorrectionLevel][qrc.Version]
	rawCodewords := qrc.GetNumRawDataModules(qrc.Version) / 8
	numShortBlocks := numBlocks - rawCodewords % numBlocks
	shortBlockLen := rawCodewords / numBlocks
	var blocks [][]uint8
	rs := utils.NewReedSolomonGenerator(blockEccLen)
	i, k := 0, 0
	for ; i < numBlocks; i++ {
		c := 1
		if i < numShortBlocks {
			c = 0
		}
		dat := (*data)[k:(k + shortBlockLen - blockEccLen + c)]
		k += len(dat)
		ecc := rs.GetRemainder(&dat)
		if i < numShortBlocks {
			dat = append(dat, 0)
		}
		dat = append(dat, ecc...)
		blocks = append(blocks, dat)
	}
	var result []uint8
	for i := 0; i < len(blocks[0]); i++ {
		for j := 0; j < len(blocks); j++ {
			if i != shortBlockLen - blockEccLen || j >= numShortBlocks {
				result = append(result, blocks[j][i])
			}
		}
	}
	if len(result) != rawCodewords {
		panic("assertion error")
	}
	return result
}

func (qrc *QrCode) DrawCodewords(data *[]uint8) {
	if len(*data) != qrc.GetNumRawDataModules(qrc.Version) / 8 {
		panic("invalid argument")
	}
	i := 0
	for right := qrc.Size - 1; right >= 1; right -= 2 {
		if right == 6 {
			right = 5
		}
		for vert := 0; vert < qrc.Size; vert++ {
			for j := 0; j < 2; j++ {
				x := right - j
				y := vert
				if ((right + 1) & 2) == 0 {
					y = qrc.Size - 1 - vert
				}
				if qrc.IsFunction[y][x] && i < len(*data) * 8 {
					qrc.Modules[y][x] = qrc.GetBit(int64((*data)[i >> 3]), 7 - (i & 7))
					i++
				}
			}
		}
	}
	if i != len(*data) * 8 {
		panic("assertion error")
	}
}

func (qrc *QrCode) ApplyMask(mask int) {
	if mask < 0 || mask > 7 {
		panic("mask value out of range")
	}
	for y := 0; y < qrc.Size; y++ {
		for x := 0; x < qrc.Size; x++ {
			var invert bool
			switch mask {
			case 0:
				invert = (x + y) % 2 == 0
			case 1:
				invert = y % 2 == 0
			case 2:
				invert = x % 3 == 0
			case 3:
				invert = (x + y) % 3 == 0
			case 4:
				invert = (x / 3 + y / 2) % 2 == 0
			case 5:
				invert = x * y % 2 + x * y % 3 == 0
			case 6:
				invert = (x * y % 2 + x * y % 3) % 2 == 0
			case 7:
				invert = ((x + y) % 2 + x * y % 3) % 2 == 0
			default:
				panic("assertion error")
			}
			qrc.Modules[y][x] = qrc.Modules[y][x] != (invert && !qrc.IsFunction[y][x])
		}
	}
}
