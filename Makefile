BINARY = qrcode

ARCH32_386 = 386
ARCH32_arm = arm
ARCH64_amd64 = amd64
ARCH64_arm64 = arm64

FLAGS = main.go

all: clean test linux windows darwin freebsd solaris

all-os: clean linux windows darwin freebsd solaris

linux: linux-386 linux-arm linux-amd64 linux-arm64

darwin: darwin-386 darwin-amd64

windows: windows-386 windows-amd64

freebsd: freebsd-386 freebsd-arm freebsd-amd64


linux-386:
	GOOS=linux GOARCH=${ARCH32_386} go build -o bin/linux/${ARCH32_386}/${BINARY} ${FLAGS}

linux-arm:
	GOOS=linux GOARCH=${ARCH32_arm} go build -o bin/linux/${ARCH32_arm}/${BINARY} ${FLAGS}

linux-amd64:
	GOOS=linux GOARCH=${ARCH64_amd64} go build -o bin/linux/${ARCH64_amd64}/${BINARY} ${FLAGS}

linux-arm64:
	GOOS=linux GOARCH=${ARCH64_arm64} go build -o bin/linux/${ARCH64_arm64}/${BINARY} ${FLAGS}


darwin-386:
	GOOS=darwin GOARCH=${ARCH32_386} go build -o bin/darwin/${ARCH32_386}/${BINARY} ${FLAGS}

darwin-amd64:
	GOOS=darwin GOARCH=${ARCH64_amd64} go build -o bin/darwin/${ARCH64_amd64}/${BINARY} ${FLAGS}


windows-386:
	GOOS=windows GOARCH=${ARCH32_386} go build -o bin/windows/${ARCH32_386}/${BINARY}.exe ${FLAGS}

windows-amd64:
	GOOS=windows GOARCH=${ARCH64_amd64} go build -o bin/windows/${ARCH64_amd64}/${BINARY}.exe ${FLAGS}


freebsd-386:
	GOOS=freebsd GOARCH=${ARCH32_386} go build -o bin/freebsd/${ARCH32_386}/${BINARY}.out ${FLAGS}

freebsd-arm:
	GOOS=freebsd GOARCH=${ARCH32_arm} go build -o bin/freebsd/${ARCH32_arm}/${BINARY}.out ${FLAGS}

freebsd-amd64:
	GOOS=freebsd GOARCH=${ARCH64_amd64} go build -o bin/freebsd/${ARCH64_amd64}/${BINARY}.out ${FLAGS}


solaris:
	GOOS=solaris GOARCH=${ARCH64_amd64} go build -o bin/solaris/${ARCH64_amd64}/${BINARY}.bin ${FLAGS}

test:
	go test ./test/...

clean:
	-rm -rf bin/

.PHONY: linux-386 linux-arm linux-amd64 linux-arm64 darwin-386 darwin-amd64 windows-386 windows-amd64 freebsd-386 freebsd-arm freebsd-amd64 solaris-amd64 test clean
