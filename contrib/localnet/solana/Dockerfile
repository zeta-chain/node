FROM ghcr.io/zeta-chain/solana-docker:1.18.15

WORKDIR /data
COPY ./start-solana.sh /usr/bin/start-solana.sh
RUN chmod +x /usr/bin/start-solana.sh
COPY ./gateway.so .
COPY ./gateway-keypair.json .

ENTRYPOINT [ "bash" ]
CMD [ "/usr/bin/start-solana.sh" ]