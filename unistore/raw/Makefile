APP_NAME = \
	unistore \
	unistore-cli

BIN_DIR = bin
TOOLS_DIR = .tools
SCHEMAS_DIR = schemas

VERSION ?= dev
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

TOOLS_SOURCES = \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 \
	google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1

TOOLS_BIN = $(addprefix $(TOOLS_DIR)/,$(patsubst %/,%,$(dir $(subst @,/,$(TOOLS_SOURCES)))))

all: $(APP_NAME)

$(APP_NAME): dependencies
	go build \
	    -ldflags "-w -s -X unistore/internal/meta.Version=$(VERSION)" \
	    -o $(BIN_DIR)/$@/$@-$(GOOS)-$(GOARCH) \
	    -trimpath \
	    cmd/$@/*.go

$(TOOLS_BIN):
	GOBIN="$(shell pwd)/$(dir $@)" go install $(subst $(TOOLS_DIR)/,,$@)

tools: $(TOOLS_BIN)

schemas: tools
	go run github.com/uber/prototool/cmd/prototool generate

dependencies:
	go mod download -x
	go mod tidy -v

lint: dependencies
	! gofmt -s -d . | grep "^"
	! go run golang.org/x/tools/cmd/goimports -d cmd/ internal/ pkg/ | grep "^"
	go run golang.org/x/lint/golint --set_exit_status ./...
	go vet ./...
	go run github.com/uber/prototool/cmd/prototool lint
	go run github.com/uber/prototool/cmd/prototool format -d

clean:
	rm -rfv $(BIN_DIR) $(TOOLS_DIR)
	rm -fv $(SCHEMAS_DIR)/**/*.pb.go

.PHONY: schemas dependencies lint clean
