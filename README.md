# QR Code Generator
[![License](https://img.shields.io/badge/BSD-3--Clause-orange.svg)](LICENSE)
[![Language](https://img.shields.io/badge/Go-1.10-blue.svg)](https://golang.org/)
[![Build Status](https://travis-ci.org/YuriyLisovskiy/qrcode.svg?branch=master)](https://travis-ci.org/YuriyLisovskiy/qrcode)
[![GoDoc](https://img.shields.io/badge/go-docs-blue.svg)](https://godoc.org/github.com/YuriyLisovskiy/qrcode/qr)
### Installation
```
$ go get github.com/YuriyLisovskiy/qrcode/qr
```
### Usage
Import package with
```go
import "github.com/YuriyLisovskiy/qrcode/qr"
```
```go
qrGen := qr.Generator{}
qrGen = qrGen.EncodeText("https://github.com/YuriyLisovskiy")
qrGen.DrawImage("qr.png", 4, 500)
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
