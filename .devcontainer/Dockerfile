FROM nikaro/alpine-dev:latest

RUN \
    sudo apk add go && \
    go get -v golang.org/x/tools/gopls && \
    go get -v github.com/uudashr/gopkgs/v2/cmd/gopkgs && \
    go get -v github.com/ramya-rao-a/go-outline && \
    go get -v github.com/go-delve/delve/cmd/dlv && \
    go get -v honnef.co/go/tools/cmd/staticcheck && \
    :
