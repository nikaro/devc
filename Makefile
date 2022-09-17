APP=devc
PREFIX?=/usr/local
_INSTDIR=${DESTDIR}${PREFIX}
BINDIR?=${_INSTDIR}/bin
SHAREDIR?=${_INSTDIR}/share
MANDIR?=${_INSTDIR}/share/man

GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

.PHONY: all
all: build

.PHONY: setup
## setup: Setup go modules
setup:
	go get -u all
	go mod tidy
	go mod vendor

.PHONY: build
## build: Build for the current target
build:
	@echo "Building..."
	env CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -mod vendor -o build/${APP}-${GOOS}-${GOARCH} .

.PHONY: man
## man: Build manpage
man:
	@echo "Building manpages..."
	build/${APP}-${GOOS}-${GOARCH} man

.PHONY: completion
## man: Build completions
completion:
	@echo "Building completions..."
	build/${APP}-${GOOS}-${GOARCH} completion bash > completions/${APP}.bash
	build/${APP}-${GOOS}-${GOARCH} completion fish > completions/${APP}.fish
	build/${APP}-${GOOS}-${GOARCH} completion zsh > completions/${APP}.zsh

.PHONY: install
## install: Install the application
install:
	@echo "Installing..."
	install -d ${BINDIR}
	install -m 755 build/${APP}-${GOOS}-${GOARCH} ${BINDIR}/${APP}
	install -d ${MANDIR}/man1
	install -m 644 $(wildcard man/${APP}*.1) ${MANDIR}/man1/
	install -d ${SHAREDIR}/bash-completion/completions
	install -m644 completions/${APP}.bash ${SHAREDIR}/bash-completion/completions/${APP}
	install -d ${SHAREDIR}/fish/vendor_completions.d
	install -m644 completions/${APP}.fish ${SHAREDIR}/fish/vendor_completions.d/${APP}.fish
	install -d ${SHAREDIR}/zsh/site-functions
	install -m644 completions/${APP}.zsh ${SHAREDIR}/zsh/site-functions/_${APP}

.PHONY: uninstall
## uninstall: Uninstall the application
uninstall:
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
## format: Runs goimports on the project
format:
	@echo "Formatting..."
	fd -t file -e go -E vendor/ | xargs goimports -l -w

.PHONY: lint
## lint: Run linters
lint:
	@echo "Linting..."
	golangci-lint run

.PHONY: test
## test: Runs go test
test:
	@echo "Testing..."
	go test ./...

.PHONY: run
## run: Runs go run
run:
	go run -race ${APP}.go

.PHONY: clean
## clean: Cleans the binary
clean:
	@echo "Cleaning..."
	@rm -rf build/
	@rm -rf dist/

.PHONY: help
## help: Prints this help message
help:
	@echo -e "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
