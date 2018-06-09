## QR Code Generator
[![License](https://img.shields.io/badge/BSD-3--Clause-orange.svg)](LICENSE)
[![Language](https://img.shields.io/badge/Go-1.10-blue.svg)](https://golang.org/)
[![Build Status](https://travis-ci.org/YuriyLisovskiy/qrcode.svg?branch=master)](https://travis-ci.org/YuriyLisovskiy/qrcode)
## Download
 ```
 $ git clone https://github.com/YuriyLisovskiy/qrcode.git
 ```
 ## Build
 > Minimum version of Go language required: `go1.10`, see [golang installation](https://golang.org/doc/install)
 
 To build for all supported operating systems:
 ```
 $ make build
 ``` 
 To run tests and build for all supported operating systems:
 ```
 $ make
 ```
 Available operating systems build:
 * **Linux**: `$ make linux` - to build all supported linux architectures, other:
 	* `$ make linux-386`
 	* `$ make linux-amd64`
 	* `$ make linux-arm`
 	* `$ make linux-arm64`
 * **Windows**: `$ make windows` - to build all supported windows architectures, other:
 	* `$ make windows-386`
 	* `$ make windows-amd64`
 * **Darwin**: `$ make darwin` - to build all supported darwin architectures, other:
 	* `$ make darwin-386`
 	* `$ make darwin-amd64`
 * **FreeBSD**: `$ make freebsd` - to build all supported freebsd architectures, other:
 	* `$ make freebsd-386`
 	* `$ make freebsd-amd64`
 	* `$ make freebsd-arm`
 * **Solaris**(only amd64):
 	* `$ make solaris`
 ## Testing
 To run tests:
 ```
 $ make test
 ```
## Author
 * [Yuriy Lisovskiy](https://github.com/YuriyLisovskiy)
 ## License
  This project is licensed under BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.