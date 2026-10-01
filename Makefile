.PHONY: build test install clean

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -s -w \
	-X 'github.com/burakberk/lubelogger-cli/internal/version.Version=$(VERSION)' \
	-X 'github.com/burakberk/lubelogger-cli/internal/version.Commit=$(COMMIT)' \
	-X 'github.com/burakberk/lubelogger-cli/internal/version.Date=$(DATE)'

build:
	@echo "Building lubelog $(VERSION)..."
	@mkdir -p cli/bin
	@cd cli && go build -ldflags="$(LDFLAGS)" -o bin/lubelog ./cmd/lubelog
	@echo "Built: cli/bin/lubelog"

test:
	@cd cli && go test -v ./...

install: build
	@cp cli/bin/lubelog /usr/local/bin/lubelog
	@echo "Installed lubelog to /usr/local/bin/lubelog"

clean:
	@rm -rf cli/bin dist
