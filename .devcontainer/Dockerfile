FROM nikaro/debian-dev:latest

# go path
ENV GOPATH ${HOME}/go
# add pyvenv to path
ENV PATH ${GOPATH}/bin:/usr/local/go/bin:${PATH}

ARG GO_VERSION=1.13.5
ARG GO_SHA256SUM=512103d7ad296467814a6e3f635631bd35574cab3369a97a323c9a585ccaa569
RUN \
    # golang
    sudo curl -LOSs https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz && \
	echo -n "${GO_SHA256SUM}	go${GO_VERSION}.linux-amd64.tar.gz" | sha256sum --check --status && \
    sudo tar -C /usr/local/ -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    sudo rm -f go${GO_VERSION}* && \
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
    :
