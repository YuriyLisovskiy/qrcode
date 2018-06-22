package main

import "github.com/YuriyLisovskiy/qrcode/src/qr"

func main() {
	text := "https://github.com/YuriyLisovskiy"
	qrGenerator := qr.Generator{}
	qrGenerator = qrGenerator.EncodeText(text)
	qrGenerator.DrawImage("qr.png", 4)
}
