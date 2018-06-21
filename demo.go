package main

import (
	"github.com/YuriyLisovskiy/qrcode/src/utils"
	"github.com/YuriyLisovskiy/qrcode/src/qrcode"
)

func printQr(qr *qrcode.QrCode) {
	border := 4
	for y := -border; y < qr.GetSize() + border; y++ {
		for x := -border; x < qr.GetSize() + border; x++ {
			if qr.GetModule(x, y) {
				print("  ")
			} else {
				print("\u2588\u2588")
			}
		}
		println()
	}
	println()
}

func main() {
//	text := "HTTPS://GITHUB.COM/YURIYLISOVSKIY"
	text := "https://github.com/YuriyLisovskiy"
//	text := "10"
	qr := qrcode.QrCode{}
	qr = qr.EncodeText(text, utils.EccLOW)
	printQr(&qr)
}
