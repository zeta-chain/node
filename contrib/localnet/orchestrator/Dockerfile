FROM zetanode:latest as zeta
FROM ethereum/client-go:v1.10.26 as geth
FROM golang:1.19-alpine as orchestrator

RUN apk --no-cache add jq curl bash make git build-base

COPY --from=zeta /usr/local/bin/zetacored /usr/local/bin/
COPY --from=zeta /usr/local/bin/zetaclientd /usr/local/bin/
COPY --from=geth /usr/local/bin/geth /usr/local/bin/

COPY --from=zeta /root/.ssh/localtest.pem.pub /root/.ssh/authorized_keys
COPY --from=zeta /root/.ssh/localtest.pem /root/.ssh/localtest.pem

COPY contrib/localnet/orchestrator/start.sh /work/
COPY contrib/localnet/orchestrator/start-upgrade.sh /work/
RUN chmod +x /work/*.sh

ENV GOPATH /go
ENV GOOS=linux
ENV CGO_ENABLED=1
WORKDIR /go/delivery/zeta-node
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make install-smoketest
RUN mv /go/bin/smoketest /usr/local/bin/smoketest
RUN chmod +x /usr/local/bin/smoketest

WORKDIR /work
