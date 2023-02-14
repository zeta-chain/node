################
# BUILD ZETAVISOR
################
FROM golang:1.18 AS cosmovisor

ENV GOBIN=/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV GOOS=linux

COPY ./* ./
RUN go get github.com/zeta-chain/cosmos-sdk/cosmovisor/cmd/cosmovisor@zetavisor-v0.1.3 && \
    go install github.com/zeta-chain/cosmos-sdk/cosmovisor/cmd/cosmovisor@zetavisor-v0.1.3

################
# BUILD MAIN IMAGE
################
FROM ubuntu:latest

EXPOSE 1317
EXPOSE 26656
EXPOSE 26657

COPY --from=cosmovisor /go/bin/cosmovisor /usr/bin/cosmovisor
COPY ./configs/* /app_version/
COPY ./scripts/network-genesis.sh /scripts/network-genesis.sh
COPY ./upgrades/ /app_version/

RUN apt update && \
    apt install curl wget vim jq glibc-source -y && \
    chmod -R a+x /scripts && \
    chmod -R a+x /app_version/ && \
    chmod -R 766 /app_version

CMD bash /scripts/network-genesis.sh