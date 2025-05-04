FROM ghcr.io/zeta-chain/solana-docker:2.0.24

# Set working directory for the container
WORKDIR /data

# Setup SSH server
# Expose port 22 for SSH access
EXPOSE 22
# Install and configure SSH server
RUN apt update && apt install -y openssh-server

# Copy and make executable the main Solana startup script
COPY ./start-solana.sh /usr/bin/start-solana.sh
RUN chmod +x /usr/bin/start-solana.sh

# Copy Solana program files and their keypairs
COPY ./gateway.so .
COPY ./gateway-keypair.json .
COPY ./connected.so .
COPY ./connected-keypair.json .
COPY ./connected_spl.so .
COPY ./connected_spl-keypair.json .
COPY ./gateway_upgrade.so .

RUN mkdir /run/sshd

ENTRYPOINT [ "bash" ]
CMD [ "/usr/bin/start-solana.sh" ]