all: test demo

test:
	go test ./qr/...

demo:
	go run demo.go
