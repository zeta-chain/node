#!/usr/bin/env bash

kill -9 $(lsof -ti:26657)
export DAEMON_HOME=$HOME/.zetacore
export DAEMON_NAME=zetacored
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
export DAEMON_RESTART_AFTER_UPGRADE=true
export CLIENT_DAEMON_NAME=zetaclientd
export CLIENT_DAEMON_ARGS="-enable-chains,GOERLI,-val,zeta"
export DAEMON_DATA_BACKUP_DIR=$DAEMON_HOME
export CLIENT_SKIP_UPGRADE=true
export UNSAFE_SKIP_BACKUP=true

rm -rf ~/.zetacore
rm -rf zetacore.log
rm -rf zetanode.log
rm -rf zetacore-debug.log
rm -rf GOERLI_debug.log
rm -rf ZetaClient.log

make install
# Genesis
mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
#mkdir -p $DAEMON_HOME/cosmovisor/upgrades/0.2.1/bin
cp $GOPATH/bin/zetacored $DAEMON_HOME/cosmovisor/genesis/bin
cp $GOPATH/bin/zetaclientd $DAEMON_HOME/cosmovisor/genesis/bin


chmod +x $DAEMON_HOME/cosmovisor/genesis/bin/zetacored
chmod +x $DAEMON_HOME/cosmovisor/genesis/bin/zetaclientd

zetacored init test --chain-id=localnet -o
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | zetacored keys add zeta --recover --keyring-backend=test
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | zetacored keys add mario --recover --keyring-backend=test
zetacored add-genesis-account $(zetacored keys show zeta -a --keyring-backend=test) 500000000000000000000000000000000stake --keyring-backend=test
zetacored add-genesis-account $(zetacored keys show mario -a --keyring-backend=test)  500000000000000000000000000000000stake --keyring-backend=test
zetacored gentx zeta 1000000000000000000000000stake --chain-id=localnet --keyring-backend=test
zetacored collect-gentxs
zetacored validate-genesis


contents="$(jq '.app_state.gov.voting_params.voting_period = "10s"' $DAEMON_HOME/config/genesis.json)" && \
echo "${contents}" > $DAEMON_HOME/config/genesis.json


cosmovisor start --home ~/.zetacore/ --p2p.laddr 0.0.0.0:27655  --grpc.address 0.0.0.0:9096 --grpc-web.address 0.0.0.0:9093 --address tcp://0.0.0.0:27659 --rpc.laddr tcp://127.0.0.1:26657 >> zetanode.log 2>&1  &

sleep 7
printf "Raising the governance proposal:\n"
zetacored tx gov submit-proposal software-upgrade 0.2.1 \
  --from zeta \
  --deposit 10000000000000000000stake \
  --upgrade-height 6 \
  --upgrade-info '{"binaries":{"zetaclientd-darwin/arm64":"https://filebin.net/4awhitgraq8eenpd/zetaclientd","zetacored-darwin/arm64":"https://filebin.net/4awhitgraq8eenpd/zetacored"}}' \
  --description "test-upgrade" \
  --title "test-upgrade" \
  --from zeta \
  --keyring-backend test \
  --chain-id localnet \
  --yes
sleep 7
zetacored tx gov vote 1 yes --from zeta --keyring-backend test --chain-id localnet --yes
clear
sleep 10
zetacored query gov proposal 1
tail -f zetanode.log