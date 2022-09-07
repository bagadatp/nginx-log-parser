GOPRIVATE := "*.autoiterative.com"

VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X main.version=${VERSION} -extldflags '-static'"
TARGET_DIR := target
CMD := log-parser
BINARY := ${TARGET_DIR}/${CMD}
CMD_FILES := $(shell find cmd -name "*.go")
PKG_FILES := $(shell find pkg -name "*.go")
TEST_INPUT_FILE := integration/testdata/sample.log

.PHONY: all
all: ${BINARY}

${TARGET_DIR}:
	@mkdir -p ${TARGET_DIR}

${BINARY}: ${TARGET_DIR} ${CMD_FILES} ${PKG_FILES}
	GOPRIVATE=${GOPRIVATE} go build ${LDFLAGS} -o ${BINARY} ./cmd/${CMD}

.PHONY: install
install:
	@go install ${LDFLAGS} ./cmd/${CMD}

.PHONY: run
run: install
	@log-parser --in=${TEST_INPUT_FILE} --out=- | jq -S

.PHONY: test
test:
	@CGO_ENABLED=1 go test -race ./...

.PHONY: integration
integration: ${BINARY}
	@./integration/test.sh
 
.PHONY: clean
clean:
	@rm -f ${BINARY}
	@rm -f ${BINARY_TEMP}
	@rm -rf ${TARGET_DIR}
	@rm -f integration/.log-parser*

