PREFIX?=/usr/local
_INSTDIR=${DESTDIR}${PREFIX}
BINDIR?=${_INSTDIR}/bin
MANDIR?=${_INSTDIR}/share/man
APP=devc

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

.PHONY: all
all: build

.PHONY: build
## build: Build for the current target
build:
	@echo "Building..."
	@env CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -mod vendor -o build/${APP}-${GOOS}-${GOARCH} main.go

.PHONY: build-all
## build-all: Build for all targets
build-all:
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(MAKE) build
	@env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(MAKE) build
	@env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(MAKE) build
	@env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(MAKE) build
	@env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(MAKE) build

.PHONY: check
## check: Check that the build is working
check:
	@./${APP}

.PHONY: install
## install: Install the application
install:
	@echo "Installing..."
	@install build/${APP}-${GOOS}-${GOARCH} ${BINDIR}/${APP}

.PHONY: uninstall
## uninstall: Uninstall the application
uninstall:
	@echo "Uninstalling..."
	@rm -rf ${BINDIR}/${APP}

.PHONY: run
## run: Runs go run main.go
run:
	go run -race main.go

.PHONY: clean
## clean: Cleans the binary
clean:
	@echo "Cleaning..."
	@rm -rf build/
	@rm -rf dist/

.PHONY: setup
## setup: Setup go modules
setup:
	@-go mod init
	@go get -u all
	@go mod tidy
	@go mod vendor

.PHONY: lint
## lint: Runs golint linter on the project
lint:
	@golint .

.PHONY: format
## format: Runs goimports on the project
format:
	@goimports -l -w .

.PHONY: test
## test: Runs go test
test:
	@go test ./...

.PHONY: help
## help: Prints this help message
help:
	@echo -e "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
