# QR Code Generator
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Minimum Go Version](https://img.shields.io/badge/Go-v1.3-blue.svg)](https://golang.org/)
[![Go Documentation](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/YuriyLisovskiy/qrcode/qr)
[![Build Status](https://travis-ci.org/YuriyLisovskiy/qrcode.svg?branch=master)](https://travis-ci.org/YuriyLisovskiy/qrcode)
[![Coverage Status](https://coveralls.io/repos/github/YuriyLisovskiy/qrcode/badge.svg?branch=master)](https://coveralls.io/github/YuriyLisovskiy/qrcode?branch=master)
### Installation
```
$ go get -u github.com/YuriyLisovskiy/qrcode/qr
```
### Usage
Import package with
```go
import "github.com/YuriyLisovskiy/qrcode/qr"
```
```go
// Text to encode.
const TEXT =
`
Nibh. Pulvinar enim porttitor tellus litora
nec vestibulum sit montes class, euismod odio
litora venenatis suscipit mi cras arcu a
dictum risus vestibulum parturient pellentesque
sociosqu vel rutrum a eros.
`

// Create an instance of Generator.
qrGen := qr.Generator{}

// Encode the text.
qrGen = qrGen.EncodeText(TEXT)

// Create an image of generated qr code.
qrGen.DrawImage("path/to/qr.png", 4, 500)
```
#### Result image:
<p align="center">
  <a href="https://github.com/YuriyLisovskiy/qrcode/blob/master/sample/qr.png"><img src="sample/qr.png"></a>
</p>

### Run
* tests: `make test`
* demo: `make demo`
### Author
* [Yuriy Lisovskiy](https://github.com/YuriyLisovskiy)
### License
This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
