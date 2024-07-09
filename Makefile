# Makefile for Go project

# Variables
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
BINARY_NAME = omicron

# Targets
all: clean deps update build 

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	$(GOGET) -v ./...

update:
	$(GOMOD) tidy
	$(GOGET) -u ./...
	$(GOMOD) vendor
	$(GOMOD) tidy

static-check:
	staticcheck ./... || true

fmt:
	go fmt ./...

security-scan:
	gosec -severity medium -confidence medium -quiet ./... || true

.PHONY: all build test clean run deps update static-check fmt security-scan
