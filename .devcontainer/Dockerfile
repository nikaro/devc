FROM nikaro/alpine-dev:latest

ENV PATH "${HOME}/go/bin:${PATH}"

RUN \
    sudo apk add go && \
    go install -v golang.org/x/tools/gopls@latest && \
    go install -v github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \
    go install -v github.com/ramya-rao-a/go-outline@latest && \
    go install -v github.com/go-delve/delve/cmd/dlv@latest && \
    go install -v honnef.co/go/tools/cmd/staticcheck@latest && \
    go install -v github.com/goreleaser/goreleaser@latest && \
    :
