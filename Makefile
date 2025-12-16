VERSION=$(shell awk -F'"' '/"version":/ {print $$4}' version.json)
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date -u -Iseconds)
GOFLAGS=-ldflags="-X github.com/volmedo/padron/pkg/build.version=$(VERSION) -X github.com/volmedo/padron/pkg/build.Commit=$(COMMIT) -X github.com/volmedo/padron/pkg/build.Date=$(DATE) -X github.com/volmedo/padron/pkg/build.BuiltBy=make"
TAGS?=

.PHONY: build padron

all: build

build: padron

# padron depends on Go sources - use shell to check if rebuild needed
padron: FORCE
	@if [ ! -f padron ] || \
	   [ -n "$$(find cmd pkg internal -name '*.go' -type f -newer padron 2>/dev/null)" ]; then \
		echo "Building padron..."; \
		go build $(GOFLAGS) $(TAGS) -o ./padron github.com/volmedo/padron/cmd; \
	fi

FORCE:

test:
	go test -v ./...
