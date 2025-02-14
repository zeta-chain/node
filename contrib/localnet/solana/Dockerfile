FROM ghcr.io/zeta-chain/solana-docker:2.0.24

WORKDIR /data
COPY ./start-solana.sh /usr/bin/start-solana.sh
RUN chmod +x /usr/bin/start-solana.sh
COPY ./gateway.so .
COPY ./gateway-keypair.json .
COPY ./connected.so .
COPY ./connected-keypair.json .


ENTRYPOINT [ "bash" ]
CMD [ "/usr/bin/start-solana.sh" ]