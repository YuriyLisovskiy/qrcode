package vars

import "github.com/YuriyLisovskiy/qrcode/src/utils"

var (
	NUMERIC      = utils.NewMode(0x1, 10, 12, 14)
	ALPHANUMERIC = utils.NewMode(0x2, 9, 11, 13)
	BYTE         = utils.NewMode(0x4, 8, 16, 16)
	KANJI        = utils.NewMode(0x8, 8, 10, 12)
	ECI          = utils.NewMode(0x7, 0, 0, 0)
)
