# QR Code Generator
[![License](https://img.shields.io/badge/BSD-3--Clause-orange.svg)](LICENSE)
[![Language](https://img.shields.io/badge/Go-1.10-blue.svg)](https://golang.org/)
[![Build Status](https://travis-ci.org/YuriyLisovskiy/qrcode.svg?branch=master)](https://travis-ci.org/YuriyLisovskiy/qrcode)
### Installation
```
$ go get github.com/YuriyLisovskiy/qrcode/src/qr
```
### Usage
Import package with
```go
import "github.com/YuriyLisovskiy/qrcode/src/qr"
```
```go
qrGen := qr.Generator{}
qrGen = qrGen.EncodeText("https://github.com/YuriyLisovskiy")
qrGen.DrawImage("qr.png", 4, 15)
```
#### Result:
<p align="center">
  <a href="https://github.com/YuriyLisovskiy/qrcode/blob/master/sample/qr.png"><img src="sample/qr.png"></a>
</p>
### Run
* tests: `make test`
* demo: `make demo`
### Author
* [Yuriy Lisovskiy](https://github.com/YuriyLisovskiy)
### License
This project is licensed under BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.