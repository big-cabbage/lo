
FROM golang:1.23.1

WORKDIR /go/src/github.com/big-cabbage/lo

COPY Makefile go.* ./

RUN make tools
