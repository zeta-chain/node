# Build Stage
FROM golang:1.20-alpine3.18 AS builder

ENV GOPATH /go
ENV GOOS=linux
ENV CGO_ENABLED=1

# Install build dependencies
RUN apk --no-cache add git make build-base jq openssh libusb-dev linux-headers bash curl tmux python3 py3-pip

# Set the working directory
WORKDIR /go/delivery/zeta-node

# Copy module files and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the source code and build the application
COPY . .

RUN make install

# Run Stage
FROM alpine:3.18

# Copy Start Script Helpers
COPY contrib/docker-scripts/* /scripts/

# Install runtime dependencies
RUN apk --no-cache add git jq bash curl tmux python3 libusb-dev linux-headers make build-base wget py3-pip qemu-img qemu-system-x86_64 && \
    pip install requests && \
    chmod a+x -R /scripts && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.31-r0/glibc-2.31-r0.apk && \
    apk add --force-overwrite --allow-untrusted glibc-2.31-r0.apk

# Copy the binaries from the build stage
COPY --from=builder /go/bin/zetaclientd /usr/local/bin/zetaclientd
COPY --from=builder /go/bin/zetacored /usr/local/bin/zetacored

# Set the working directory
WORKDIR /usr/local/bin

# Set the default shell
ENV SHELL /bin/bash

EXPOSE 26656
EXPOSE 1317
EXPOSE 8545
EXPOSE 8546
EXPOSE 9090
EXPOSE 26657
EXPOSE 9091
