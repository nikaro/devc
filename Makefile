APP=devc
PREFIX?=/usr/local
_INSTDIR=${DESTDIR}${PREFIX}
BINDIR?=${_INSTDIR}/bin
SHAREDIR?=${_INSTDIR}/share
MANDIR?=${_INSTDIR}/share/man

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)
VERSION = $(shell git describe --always --dirty)

.PHONY: all
all: build

.PHONY: setup
setup: ## Setup go modules
	go get -u all
	go mod tidy
	go mod vendor

.PHONY: build
build: ## Build for the current target
	@echo "Building..."
	env CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -mod vendor -ldflags="-s -w -X 'main.version=${VERSION}'" -o build/${APP}-${GOOS}-${GOARCH} .

.PHONY: man
man: ## Build manpage
	@echo "Building manpages..."
	build/${APP}-${GOOS}-${GOARCH} man

.PHONY: completion
completion: ## Build completions
	@echo "Building completions..."
	build/${APP}-${GOOS}-${GOARCH} completion bash > completions/${APP}
	build/${APP}-${GOOS}-${GOARCH} completion fish > completions/${APP}.fish
	build/${APP}-${GOOS}-${GOARCH} completion zsh > completions/_${APP}

.PHONY: install
install: ## Install the application
	@echo "Installing..."
	install -d ${BINDIR}
	install -m 755 build/${APP}-${GOOS}-${GOARCH} ${BINDIR}/${APP}
	install -d ${MANDIR}/man1
	install -m 644 $(wildcard man/${APP}*.1) ${MANDIR}/man1/
	install -d ${SHAREDIR}/bash-completion/completions
	install -m 644 completions/${APP} ${SHAREDIR}/bash-completion/completions/${APP}
	install -d ${SHAREDIR}/fish/vendor_completions.d
	install -m 644 completions/${APP}.fish ${SHAREDIR}/fish/vendor_completions.d/${APP}.fish
	install -d ${SHAREDIR}/zsh/site-functions
	install -m 644 completions/_${APP} ${SHAREDIR}/zsh/site-functions/_${APP}

.PHONY: uninstall
uninstall: ## Uninstall the application
	@echo "Uninstalling..."
	rm -f ${BINDIR}/${APP}
	rm -f $(wildcard ${MANDIR}/man1/${APP}*.1)
	rm -f ${SHAREDIR}/bash-completion/completions/${APP}
	rm -f ${SHAREDIR}/fish/vendor_completions.d/${APP}.fish
	rm -f ${SHAREDIR}/zsh/site-functions/_${APP}
	-rmdir ${BINDIR}
	-rmdir ${SHAREDIR}
	-rmdir ${MANDIR}/man1
	-rmdir ${MANDIR}
	-rmdir ${SHAREDIR}/bash-completion/completions
	-rmdir ${SHAREDIR}/bash-completion
	-rmdir ${SHAREDIR}/fish/vendor_completions.d
	-rmdir ${SHAREDIR}/fish
	-rmdir ${SHAREDIR}/zsh/site-functions
	-rmdir ${SHAREDIR}/zsh

.PHONY: format
format: ## Runs goimports on the project
	@echo "Formatting..."
	find . -type f -name '*.go' -not -path './vendor/*' | xargs goimports -l -w

.PHONY: lint
lint: ## Run linters
	@echo "Linting..."
	golangci-lint run

.PHONY: test
test: ## Runs go test
	@echo "Testing..."
	go test ./...

.PHONY: run
run: ## Runs go run
	go run -race ${APP}.go

.PHONY: clean
clean: ## Cleans the binary
	@echo "Cleaning..."
	@rm -rf build/
	@rm -rf dist/

.PHONY: help
help: ## Print this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
