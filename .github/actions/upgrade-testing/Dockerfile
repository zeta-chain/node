################
# BUILD ZETAVISOR
################
FROM golang:1.20 AS cosmovisor

ENV GOBIN=/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN apt update && \
    apt install unzip wget -y && \
    wget https://github.com/zeta-chain/cosmos-sdk/archive/refs/heads/main.zip && \
    unzip main.zip && \
    cd cosmos-sdk-main/tools/cosmovisor && \
    make cosmovisor && \
    go get github.com/zeta-chain/cosmos-sdk/cosmovisor/cmd/cosmovisor@zetavisor-v0.1.3 && \
    go install github.com/zeta-chain/cosmos-sdk/cosmovisor/cmd/cosmovisor@zetavisor-v0.1.3

################
# BUILD MAIN IMAGE
################
FROM golang:1.20

EXPOSE 1317
EXPOSE 26656
EXPOSE 26657

COPY --from=cosmovisor /go/bin/cosmovisor /usr/bin/cosmovisor
COPY ./configs/* /app_version/
COPY ./scripts/network-genesis.sh /scripts/network-genesis.sh
COPY ./upgrades/ /app_version/

RUN apt update && \
    apt-get update && \
    apt-get dist-upgrade -y && \
    apt upgrade libc6 -y && \
    apt install curl wget vim jq bash glibc-source -y && \
    chmod -R a+x /scripts && \
    chmod -R a+x /app_version/ && \
    chmod -R 766 /app_version

CMD bash /scripts/network-genesis.sh