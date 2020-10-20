FROM nikaro/debian-dev:latest

# go path
ENV GOPATH ${HOME}/go
# add go to path
ENV PATH ${GOPATH}/bin:/opt/go/bin:${PATH}

ENV GO_VERSION 1.15.3
RUN \
    # golang
    sudo curl -LOSs https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    sudo tar -C /opt -xaf go${GO_VERSION}.linux-amd64.tar.gz && \
    sudo rm -f go${GO_VERSION}.linux-amd64.tar.gz && \
    :

RUN \
    # go packages
    go get \
        github.com/mdempsky/gocode \
        github.com/uudashr/gopkgs/cmd/gopkgs \
        github.com/ramya-rao-a/go-outline \
        github.com/acroca/go-symbols \
        golang.org/x/tools/cmd/guru \
        golang.org/x/tools/cmd/gorename \
        github.com/cweill/gotests/gotests \
        github.com/fatih/gomodifytags \
        github.com/josharian/impl \
        github.com/davidrjenni/reftools/cmd/fillstruct \
        github.com/haya14busa/goplay/cmd/goplay \
        github.com/godoctor/godoctor \
        github.com/go-delve/delve/cmd/dlv \
        github.com/stamblerre/gocode \
        github.com/rogpeppe/godef \
        github.com/sqs/goreturns \
        golang.org/x/lint/golint \
        golang.org/x/tools/cmd/goimports \
        golang.org/x/tools/gopls \
    && \
    go build -o ${GOPATH}/bin/gocode-gomod github.com/stamblerre/gocode && \
    :
