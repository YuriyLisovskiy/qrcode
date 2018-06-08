package utils

type ReedSolomonGenerator struct {
	Coefficients []uint8
}

func NewReedSolomonGenerator(degree int) ReedSolomonGenerator {
	newRSG := ReedSolomonGenerator{}
	if degree < 1 || degree > 255 {
		panic("degree out of range")
	}
	newRSG.Coefficients = make([]uint8, degree)
	newRSG.Coefficients[degree - 1] = 1
	root := uint8(1)
	for i := 0; i < degree; i++ {
		for j := 0; j < len(newRSG.Coefficients); j++ {
			newRSG.Coefficients[j] = newRSG.Multiply(newRSG.Coefficients[j], root)
			if j + 1 < len(newRSG.Coefficients) {
				newRSG.Coefficients[j] ^= newRSG.Coefficients[j + 1]
			}
		}
		root = newRSG.Multiply(root, 0x02)
	}
	return newRSG
}

func (rsg ReedSolomonGenerator) GetReminder(data *[]uint8) []uint8 {
	result := make([]uint8, len(rsg.Coefficients))
	for _, b := range *data {
		factor := b ^ result[0]
		result = result[1:]
		result = append(result, 0)
		for j := 0; j < len(result); j++ {
			result[j] ^= rsg.Multiply(rsg.Coefficients[j], factor)
		}
	}
	return result
}

func (rsg ReedSolomonGenerator) Multiply(x, y uint8) uint8 {
	z := uint8(0)
	for i := 7; i >= 0; i-- {
		z = (z << 1) ^ ((z >> 7) * 0x11D)
		z ^= ((y >> uint8(i)) & 1) * x
	}
	if z >> 8 != 0 {
		panic("assertion error")
	}
	return z
}
