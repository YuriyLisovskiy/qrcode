package main

import "github.com/YuriyLisovskiy/qrcode/src/qr"

func main() {
//	text := "HTTPS://GITHUB.COM/YURIYLISOVSKIY"
	text := "https://github.com/YuriyLisovskiy/best-friend-telegram-bot"
//	text := "10"

	qrcode := qr.Generator{}
	qrcode = qrcode.EncodeText(text)
	qrcode.Draw(4)
}
