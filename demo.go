package main

import "github.com/YuriyLisovskiy/qrcode/qr"

func main() {
	text := "https://github.com/YuriyLisovskiy"
	qrGenerator := qr.Generator{}
	qrGenerator = qrGenerator.EncodeText(text)
	qrGenerator.DrawImage("sample/qr.png", 4, 500)
	qrGenerator.Draw(4)
}
