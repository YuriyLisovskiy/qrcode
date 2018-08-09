#  Copyright (c) 2018 Yuriy Lisovskiy
#  Distributed under the Apache License Version 2.0,
#  see the accompanying file LICENSE or https://opensource.org/licenses/Apache-2.0

PACKAGES = ./qr
COVER_OUT = coverage.out

all: test demo

test:
	@echo Running tests...
	@go test -v -timeout 1h -covermode=count -coverprofile=$(COVER_OUT) ${PACKAGES}
	@echo Done

clean:
	@rm -rf $(COVER_OUT)

demo:
	@go run demo.go
