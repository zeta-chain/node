FROM golang:latest

ARG ACCOUNT_PASSWORD

RUN apt-get update -y && apt-get upgrade -y \
    && apt install build-essential git -y \
    && mkdir -p /bor

COPY ./bor .
RUN make bor
COPY genesis.json /datadir/
RUN cp ./build/bin/bor /usr/local/bin/bor

RUN #bor init /datadir/genesis.json
RUN rm -f ~/.ethereum/bor/nodekey
RUN echo ${ACCOUNT_PASSWORD} > /tmp/password
RUN echo ${ACCOUNT_PASSWORD} | bor account new
RUN #bor account new --password /tmp/password

ENV SHELL /bin/bash
EXPOSE 8545 8546 8547 30303 30303/udp

ENTRYPOINT ["bor"]