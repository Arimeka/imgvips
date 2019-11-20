FROM ubuntu:focal

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
		g++ \
		gcc \
		git \
		libc6-dev \
		make \
		pkg-config \
		wget \
		libvips-dev=8.8.3-3 && \
	  rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.11.13

RUN wget -O go.tgz https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go.tgz && rm go.tgz && \
    export PATH="/usr/local/go/bin:$PATH" && \
    go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN go get golang.org/x/lint/golint && \
    go get -v github.com/go-lintpack/lintpack/... && \
    go get -v github.com/go-critic/go-critic/... && \
    go get -v github.com/fzipp/gocyclo && \
    go get -v github.com/gordonklaus/ineffassign && \
    lintpack build -o $GOPATH/bin/gocritic -linter.name='gocritic' github.com/go-critic/go-critic/checkers

WORKDIR /test

COPY . .

CMD exec make test
