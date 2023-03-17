FROM golang:1.19-alpine AS builder

RUN apk --no-cache add git make build-base jq

ENV GOPATH /go
WORKDIR /go/delivery/zeta-node
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN make install
RUN make install-smoketest

FROM golang:1.19-alpine

RUN apk --no-cache add openssh jq tmux vim curl bash
RUN ssh-keygen -A
RUN mkdir /root/.ssh

COPY --from=builder /go/bin/zetaclientd /usr/local/bin
COPY --from=builder /go/bin/zetacored /usr/local/bin
COPY --from=builder /go/bin/smoketest /usr/local/bin
COPY contrib/localnet/meta.pem.pub /root/.ssh/authorized_keys
COPY contrib/localnet/meta.pem /root/.ssh/meta.pem
COPY contrib/localnet/scripts /root
COPY contrib/localnet/preparams /root/preparams
COPY contrib/localnet/ssh_config /root/.ssh/config
COPY contrib/localnet/zetacored /root/zetacored
COPY contrib/localnet/tss /root/tss

RUN chmod 755 /root/reset-testnet.sh
RUN chmod 755 /root/start-zetacored.sh
RUN chmod 755 /root/start-zetaclientd.sh
RUN chmod 755 /root/seed.sh
RUN chmod 755 /root/keygen.sh

WORKDIR /usr/local/bin
ENV SHELL /bin/sh
EXPOSE 22

ENTRYPOINT ["/usr/sbin/sshd", "-D"]
