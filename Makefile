# Makefile for Go project

# Variables
GOCMD = go
GOLINTCMD = golangci-lint
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOLINT = $(GOLINTCMD) run 
GOTEST = $(GOCMD) test
GOMOD = $(GOCMD) mod
GOFMT = $(GOCMD) fmt
BINARY_NAME = omicron

# Targets
all: clean deps update build 

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v -race ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	$(GOMOD) download
	$(GOMOD) vendor

lint:
	$(GOLINT) ./...

fmt:
	$(GOFMT) ./...

security-scan:
	gosec -severity medium -confidence medium -quiet ./... || true

.PHONY: all build test clean run deps lint fmt security-scan
