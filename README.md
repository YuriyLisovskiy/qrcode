# QR Code Generator
[![License](https://img.shields.io/badge/BSD-3--Clause-orange.svg)](LICENSE)
[![Language](https://img.shields.io/badge/Go-1.10-blue.svg)](https://golang.org/)
[![Build Status](https://travis-ci.org/YuriyLisovskiy/qrcode.svg?branch=master)](https://travis-ci.org/YuriyLisovskiy/qrcode)
### Usage
```go
package main

import "github.com/YuriyLisovskiy/qrcode/src/qr"

func main() {
	qrGen := qr.Generator{}
	qrGen = qrGen.EncodeText("https://github.com/YuriyLisovskiy")
	qrGen.Draw(4)
}
```
#### Result:
![sample image](sample/qr.png =250x)
### Run
* tests: `make test`
* demo: `make demo`
### Author
* [Yuriy Lisovskiy](https://github.com/YuriyLisovskiy)
### License
This project is licensed under BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.