//  Copyright (c) 2018 Yuriy Lisovskiy
//  Distributed under the Apache License Version 2.0,
//  see the accompanying file LICENSE or https://opensource.org/licenses/Apache-2.0

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
