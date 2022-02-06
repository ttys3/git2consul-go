BINARY = git2consul
COMMIT := $(shell git rev-parse HEAD)
BRANCH := $(shell git symbolic-ref --short -q HEAD || echo HEAD)
DATE := $(shell date -u +%Y%m%d-%H:%M:%S)
TAG := $(shell git describe --tags `git rev-list --tags --max-count=1`)
VERSION_PKG = github.com/KohlsTechnology/git2consul-go/pkg/version
LDFLAGS := "-X ${VERSION_PKG}.Branch=${BRANCH} -X ${VERSION_PKG}.BuildDate=${DATE} \
	-X ${VERSION_PKG}.GitSHA1=${COMMIT} -X ${VERSION_PKG}.Version=${TAG}"

.PHONY: all
all: build

.PHONY: clean
clean:
	rm -rf $(BINARY) dist/

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o $(BINARY) -ldflags $(LDFLAGS)

dump:
	./$(BINARY) -dump > ./config.sample.yaml

run:
	./$(BINARY) -loglvl debug -config ./config.sample.yaml

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: image
image:
	docker build . -t quay.io/kohlstechnology/git2consul:latest

.PHONY: test
test: lint-all test-unit

.PHONY: test-unit
test-unit:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

# Make sure go.mod and go.sum are not modified
.PHONY: test-dirty
test-dirty: vendor build
	go mod tidy
	git diff --exit-code
	# TODO: also check that there are no untracked files, e.g. extra .go

# Make sure goreleaser is working
.PHONY: test-release
test-release:
	BRANCH=$(BRANCH) COMMIT=$(COMMIT) DATE=$(DATE) VERSION_PKG=$(VERSION_PKG) goreleaser release --snapshot --skip-publish --rm-dist

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run

.PHONY: lint-all
lint-all: golangci-lint

# Requires GITHUB_TOKEN environment variable to be set
.PHONY: release
release:
	BRANCH=$(BRANCH) COMMIT=$(COMMIT) DATE=$(DATE) VERSION_PKG=$(VERSION_PKG) goreleaser  release --rm-dist
