# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
BINARY_NAME=benz
MAIN_PATH=main.go

all:test build

build:
	$(GOBUILD) -o dist/$(BINARY_NAME) -v $(MAIN_PATH)

test:
	$(GOTEST) -v ./...

testcover:
	$(GOTEST) -v ./... -coverprofile=benz.cov ./...

cover:
	$(GOTEST) cover -func benz.cov

