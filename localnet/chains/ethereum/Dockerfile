FROM ethereum/client-go:v1.10.1

ARG ACCOUNT_PASSWORD

COPY ./genesis.json /tmp

RUN geth init /tmp/genesis.json \
    && rm -f ~/.ethereum/geth/nodekey \
    && echo ${ACCOUNT_PASSWORD} > /tmp/password \
    && geth account new --password /tmp/password 

ENTRYPOINT ["geth"]