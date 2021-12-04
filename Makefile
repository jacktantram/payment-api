
BIN="/usr/local/bin"
VERSION="1.0.0-rc6"
BINARY_NAME="buf"
LINTER_VERSION=v1.42.1

.PHONY:  build clean run

# Installs buf for protobuf generation
get-buf:
	curl -sSL \
        "https://github.com/bufbuild/buf/releases/download/v${VERSION}/${BINARY_NAME}-$(shell uname -s)-$(shell uname -m)" \
        -o "${BIN}/${BINARY_NAME}" && \
      chmod +x "${BIN}/${BINARY_NAME}"

# Generates protobuf using buf
proto-generate: get-buf
	rm -rf build
	${BIN}/${BINARY_NAME} generate

# Runs the go testing stage
test:
	go test ./...

# Runs the go testing stage checking for race conditions.
test-race:
	go test ./... -race


get-linter:
	command -v golangci-lint || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh  | sh -s --  -b ${GOPATH}/bin ${LINTER_VERSION}
# Runs the Go linter
lint: get-linter
	golangci-lint run


# Builds the projects docker containers
build:
	docker-compose build

# Runs the projects docker containers
run:
	docker-compose up -d

# Runs the go tests with integration tests
test-integration:
	go test ./... -tags=integration -count=1
