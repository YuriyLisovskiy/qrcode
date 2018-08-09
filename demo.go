package main

import "github.com/YuriyLisovskiy/qrcode/qr"

const TEXT =
`
Nibh. Pulvinar enim porttitor tellus litora
nec vestibulum sit montes class, euismod odio
litora venenatis suscipit mi cras arcu a
dictum risus vestibulum parturient pellentesque
sociosqu vel rutrum a eros.
`

func main() {
	qrGenerator := qr.Generator{}
	qrGenerator = qrGenerator.EncodeText(TEXT)
	qrGenerator.DrawImage("sample/qr.png", 4, 500)
	qrGenerator.Draw(4)
}
