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
	staticcheck ./...

fmt:
	go fmt ./...

security-scan:
	gosec ./...

.PHONY: all build test clean run deps
