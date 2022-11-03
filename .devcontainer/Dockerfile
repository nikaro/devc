FROM alpine:edge

RUN \
    apk add --no-cache \
        build-base \
        git \
        go \
        neovim \
    && \
    :

RUN \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0 && \
    go install golang.org/x/tools/cmd/goimports@v0.1.12 && \
    go install golang.org/x/tools/gopls@v0.9.5 && \
    :
