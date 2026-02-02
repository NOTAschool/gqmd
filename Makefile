GO := GOROOT=/opt/homebrew/Cellar/go/1.25.5/libexec GOSUMDB=sum.golang.org /opt/homebrew/bin/go
BINARY := gqmd

.PHONY: all build test clean tidy

all: build

tidy:
	$(GO) mod tidy

build: tidy
	$(GO) build -o $(BINARY) ./cmd/gqmd

test:
	$(GO) test -v ./...

clean:
	rm -f $(BINARY)
