GO := GOROOT=/opt/homebrew/Cellar/go/1.25.5/libexec GOSUMDB=sum.golang.org /opt/homebrew/bin/go
BINARY := gqmd
VERSION := 0.26.2.2

.PHONY: all build test clean tidy build-macos

all: build

tidy:
	$(GO) mod tidy

build: tidy
	$(GO) build -o $(BINARY) ./cmd/gqmd

test:
	$(GO) test -v ./...

clean:
	rm -f $(BINARY) $(BINARY)-darwin-*

build-macos: tidy
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags="-s -w" -o $(BINARY)-darwin-arm64 ./cmd/gqmd
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="-s -w" -o $(BINARY)-darwin-amd64 ./cmd/gqmd
