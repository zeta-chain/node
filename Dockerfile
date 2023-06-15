FROM golang:1.19-alpine

ENV GOPATH /go
ENV GOOS=linux
ENV CGO_ENABLED=1

RUN apk --no-cache add git make build-base jq openssh libusb-dev linux-headers bash curl tmux
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
#
#FROM golang:1.19-alpine

#RUN apk --no-cache add openssh jq tmux vim curl bash
RUN ssh-keygen -A
WORKDIR /root

RUN cp /root/.ssh/localtest.pem.pub /root/.ssh/authorized_keys

RUN cp /go/bin/zetaclientd /usr/local/bin
RUN cp /go/bin/zetacored /usr/local/bin
RUN cp /go/bin/smoketest /usr/local/bin

COPY contrib/localnet/scripts /root
COPY contrib/localnet/preparams /root/preparams
COPY contrib/localnet/ssh_config /root/.ssh/config
COPY contrib/localnet/zetacored /root/zetacored
COPY contrib/localnet/tss /root/tss

RUN chmod 755 /root/*.sh

WORKDIR /usr/local/bin
ENV SHELL /bin/sh
EXPOSE 22

ENTRYPOINT ["/usr/sbin/sshd", "-D"]
