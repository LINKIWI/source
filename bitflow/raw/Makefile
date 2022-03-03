APP_NAME = bitflow
BIN_DIR = bin

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

all: $(APP_NAME)

$(APP_NAME):
	go build \
	    -ldflags "-w -s -X bitflow/internal/meta.Version=$(VERSION)" \
	    -o $(BIN_DIR)/$(APP_NAME)-$(GOOS)-$(GOARCH) \
	    -trimpath \
	    cmd/bitflow/main.go

lint:
	! gofmt -s -d . | grep "^"
	! go run golang.org/x/tools/cmd/goimports -d . | grep "^"
	go run golang.org/x/lint/golint --set_exit_status ./...
	go vet ./...

clean:
	rm -rfv $(BIN_DIR)

.PHONY: lint clean
