#!/bin/bash
# Install ZetaCore PreReqs
apt-get -y install make gcc jq

wget https://golang.org/dl/go1.17.3.linux-arm64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.3.linux-arm64.tar.gz
rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.17.3.linux-arm64.tar.gz

cd /root
rm -rf ./zeta-node
git clone --branch athens-aws git@github.com:zeta-chain/zeta-node.git
cd zeta-node

pwd 
ls
whoami

export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:/root/go/bin
export GOPATH=/root/go
export HOME=/root

make install