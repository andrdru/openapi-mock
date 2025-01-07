-include .env
export

GOLANGCI_LINT_IMAGE="golangci/golangci-lint:v1.63.4"
GOLANGCI_LINT_CMD=docker run --rm \
                   	--volume $(shell pwd):/app \
                   	--volume $(shell go env GOPATH)/pkg/mod:/go/pkg/mod \
                   	--workdir /app \
                   	--env HOME="/tmp" \
                   	"${GOLANGCI_LINT_IMAGE}"

LOCAL_IMAGE="openapi_mock_builder"

LOCAL_IMAGE_CMD=docker run --rm \
                   	--volume $(shell pwd):/app \
                   	--volume $(shell go env GOPATH)/pkg/mod:/go/pkg/mod \
                   	--workdir /app \
                   	--env HOME="/tmp" \
                   	--network="host" \
                   	"${LOCAL_IMAGE}"


# lint run golangci-lint
.PHONY: lint
lint:
	@ ${GOLANGCI_LINT_CMD} golangci-lint run ./... -v

# build local image
#
# make build
# make build FLAGS="--no-cache"
.PHONY: build
build:
	@ docker build --rm ${FLAGS} docker/ -t "${LOCAL_IMAGE}"

# gen run all codegen
.PHONY: gen
gen: build gen-go

# gen-go generate all //go:generate
.PHONY: gen-go
gen-go:
	@ ${LOCAL_IMAGE_CMD} go generate -v ./...
	@ git add .
