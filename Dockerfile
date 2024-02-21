# Purpose: This Dockerfile creates an environment for running ZetaChain
# It contains:
# - zetacored: the ZetaChain node binary
# - zetaclientd: the ZetaChain client binary for observers
# - zetae2e: the ZetaChain end-to-end tests CLI

FROM golang:1.20-alpine3.18

ENV GOPATH /go
ENV GOOS=linux
ENV CGO_ENABLED=1

RUN apk --no-cache add git make build-base jq openssh libusb-dev linux-headers bash curl tmux
RUN ssh-keygen -b 2048 -t rsa -f /root/.ssh/localtest.pem -q -N ""

WORKDIR /go/delivery/zeta-node
COPY go.mod .
COPY go.sum .
#RUN --mount=type=cache,target=/root/.cache/go-build \
#    go mod download
RUN go mod download
COPY . .

#RUN --mount=type=cache,target=/root/.cache/go-build \
#    make install
#RUN --mount=type=cache,target=/root/.cache/go-build \
#    make install-zetae2e
RUN make install
RUN make install-zetae2e
#
#FROM golang:1.20-alpine

#RUN apk --no-cache add openssh jq tmux vim curl bash
RUN ssh-keygen -A
WORKDIR /root

RUN cp /root/.ssh/localtest.pem.pub /root/.ssh/authorized_keys

RUN cp /go/bin/zetaclientd /usr/local/bin
RUN cp /go/bin/zetacored /usr/local/bin
RUN cp /go/bin/zetae2e /usr/local/bin

COPY contrib/localnet/scripts /root
COPY contrib/localnet/preparams /root/preparams
COPY contrib/localnet/ssh_config /root/.ssh/config
COPY contrib/localnet/zetacored /root/zetacored
COPY contrib/localnet/tss /root/tss

RUN chmod 755 /root/*.sh
RUN chmod 700 /root/.ssh
RUN chmod 600 /root/.ssh/*

WORKDIR /usr/local/bin
ENV SHELL /bin/sh
EXPOSE 22

ENTRYPOINT ["/usr/sbin/sshd", "-D"]
