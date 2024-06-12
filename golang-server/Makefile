# Makefile for Go project

# Variables
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
BINARY_NAME = omicron

# Targets
all: clean build

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
	
static-check:
	staticcheck ./... || true

fmt:
	go fmt ./... 
	
vendor:
    go mod vendor

security-scan:
	gosec-severity medium -confidence medium -quiet ./... || true

.PHONY: all build test clean run deps
