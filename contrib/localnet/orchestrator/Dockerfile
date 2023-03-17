FROM zetanode:latest as zeta
FROM ethereum/client-go:v1.10.26 as geth
FROM golang:1.19-alpine as orchestrator

RUN apk --no-cache add jq curl bash make git build-base

COPY --from=zeta /usr/local/bin/zetacored /usr/local/bin/
COPY --from=zeta /usr/local/bin/zetaclientd /usr/local/bin/
COPY --from=geth /usr/local/bin/geth /usr/local/bin/


COPY contrib/localnet/orchestrator/start.sh /work/
RUN chmod +x /work/start.sh

ENV GOPATH /go
WORKDIR /go/delivery/zeta-node
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make install-smoketest
RUN mv /go/bin/smoketest /usr/local/bin/smoketest
RUN chmod +x /usr/local/bin/smoketest

WORKDIR /work
