GO = go
GOTEST ?= go test
export TEST_COUNT ?= 1
export TEST_ARGS ?= 1
bin:
	mkdir -p bin/

.PHONY: build
build: bin
	$(GO) build  -ldflags "${LDFLAGS}" -o bin/monitor ./cmd/monitor/*.go

clear:
	rm -rf bin

.PHONY: test
test:
	$(GO) test -v ./... -count $(TEST_COUNT) -race $(TEST_ARGS)