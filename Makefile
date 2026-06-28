BINARY_NAME := pkgui
BUILD_DIR := build

.PHONY: build install run clean test lint

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

install: build
	@GOBIN=$${GOBIN:-$$(go env GOPATH)/bin}; \
	echo "installing to $$GOBIN/$(BINARY_NAME)"; \
	cp $(BUILD_DIR)/$(BINARY_NAME) "$$GOBIN/$(BINARY_NAME)"

run:
	go run .

clean:
	@rm -rf $(BUILD_DIR)

test:
	go test ./...

lint:
	golangci-lint run ./...
