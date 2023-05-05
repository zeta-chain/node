FROM golang:1.19-alpine AS builder

ENV GOPATH /go
ENV GOOS=linux
ENV CGO_ENABLED=1

RUN apk --no-cache add git make build-base jq openssh libusb-dev linux-headers
RUN ssh-keygen -b 2048 -t rsa -f /root/.ssh/localtest.pem -q -N ""

WORKDIR /go/delivery/zeta-node
COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    make install
RUN --mount=type=cache,target=/root/.cache/go-build \
    make install-smoketest

FROM golang:1.19-alpine

RUN apk --no-cache add openssh jq tmux vim curl bash 
RUN ssh-keygen -A

COPY --from=builder /root/.ssh/localtest.pem.pub /root/.ssh/authorized_keys
COPY --from=builder /root/.ssh/localtest.pem.pub /root/.ssh/localtest.pem.pub
COPY --from=builder /root/.ssh/localtest.pem /root/.ssh/localtest.pem
COPY --from=builder /go/bin/zetaclientd /usr/local/bin
COPY --from=builder /go/bin/zetacored /usr/local/bin
COPY --from=builder /go/bin/smoketest /usr/local/bin

COPY contrib/localnet/scripts /root
COPY contrib/localnet/preparams /root/preparams
COPY contrib/localnet/ssh_config /root/.ssh/config
COPY contrib/localnet/zetacored /root/zetacored
COPY contrib/localnet/tss /root/tss

RUN chmod 755 /root/reset-testnet.sh
RUN chmod 755 /root/start-zetacored.sh
RUN chmod 755 /root/start-zetaclientd.sh
RUN chmod 755 /root/start-zetaclientd-genesis.sh
RUN chmod 755 /root/genesis.sh
RUN chmod 755 /root/seed.sh
RUN chmod 755 /root/keygen.sh
RUN chmod 755 /root/os-info.sh
RUN chmod 755 /root/start-zetaclientd-p2p-diag.sh

WORKDIR /usr/local/bin
ENV SHELL /bin/sh
EXPOSE 22

ENTRYPOINT ["/usr/sbin/sshd", "-D"]
