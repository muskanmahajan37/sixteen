# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_DIR=output/
PACKAGE_NAME=sixteen
BINARY_LINUX=$(BINARY_DIR)/$(PACKAGE_NAME)_linux
BINARY_MACOS=$(BINARY_DIR)/$(PACKAGE_NAME)_macos
BINARY_WINDOWS=$(BINARY_DIR)/$(PACKAGE_NAME)_windows

all: test build
build:
	$(GOBUILD) -o $(BINARY_DIR) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)
run:
	$(GOBUILD) -o $(BINARY_DIR) -v ./...
	./$(BINARY_DIR)

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v
build-macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_MACOS) -v
