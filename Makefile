BINARY=cfgm

.DEFAULT_GOAL := help

# Git current hash
GIT_HASH=`git rev-parse --short HEAD`

# Git current tag
GIT_TAG=`git tag -l --contains HEAD`

# Build date time
BUILD=`date +%FT%T%z`

# Setup the -ldflags options
LDFLAGS=-ldflags "-X main.GitHash=${GIT_HASH}${GIT_TAG} -X main.Build=${BUILD}"

dev-start:
	docker-compose -f docker-compose.dev.yml up -d

dev-stop:
	docker-compose -f docker-compose.dev.yml down

dep:
	dep ensure -v -vendor-only
	
build:
	@echo "Build - ${GIT_HASH}${GIT_TAG}"
	CGO_ENABLED=0 go build ${LDFLAGS} -v -a -installsuffix cgo -o ${BINARY} ./cmd/cfgm

docker:
	docker build --build-arg SSH_DEP_PRIVATE_KEY="$$pk" -t $$c .

check: vet test ## Runs all tests

test: ## Run the unit tests
	go test -race -cover -v $(shell go list ./... | grep -v /vendor/)

vet: ## Run the vet tool
	go vet $(shell go list ./... | grep -v /vendor/)

clean: ## Clean up build artifacts
	go clean

help: ## Display this help message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.SILENT: build test lint vet clean docker-build docker-push help
