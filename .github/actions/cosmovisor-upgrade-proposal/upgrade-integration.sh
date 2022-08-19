#!/usr/bin/env bash


clibuilder()
{
   echo ""
   echo "Usage: $0 -u UpgradeName -c CurrentBinary -n NewBinary"
   echo -e "\t-u Name of the upgrade [Must match a handler defined in setup-handlers.go in NewBinary]"
   echo -e "\t-c Branch name for old binary (Upgrade From)"
   echo -e "\t-n Branch name for new binary (Upgrade To)"
   exit 1 # Exit script after printing help
}

while getopts "u:c:n:" opt
do
   case "$opt" in
      u ) UpgradeName="$OPTARG" ;;
      c ) CurrentBinary="$OPTARG" ;;
      n ) NewBinary="$OPTARG" ;;
      ? ) clibuilder ;; # Print cliBuilder in case parameter is non-existent
   esac
done

if [ -z "$UpgradeName" ] || [ -z "$CurrentBinary" ] || [ -z "$NewBinary" ]
then
   echo "Some or all of the parameters are empty";
   clibuilder
fi


export DAEMON_HOME=$HOME/.zetacore
export DAEMON_NAME=zetacored
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
export DAEMON_RESTART_AFTER_UPGRADE=true
export CLIENT_DAEMON_NAME=zetaclientd
export CLIENT_DAEMON_ARGS="-enable-chains,GOERLI,-val zeta"
#export DAEMON_DATA_BACKUP_DIR=$DAEMON_HOME

make clean
rm -rf ~/.zetacore
rm -rf zetacore.log

rm -rf $GOPATH/bin/zetacored
rm -rf $GOPATH/bin/old/zetacored
rm -rf $GOPATH/bin/new/zetacored

# Setup old binary and start chain
mkdir -p  $GOPATH/bin/old
mkdir -p  $GOPATH/bin/new

git checkout $CurrentBinary
make install-zetacore
cp $GOPATH/bin/zetacored $GOPATH/bin/old/
zetacored init test --chain-id=localnet -o

echo "Generating deterministic account - zeta"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | zetacored keys add zeta --recover --keyring-backend=test

echo "Generating deterministic account - mario"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | zetacored keys add mario --recover --keyring-backend=test


zetacored add-genesis-account $(zetacored keys show zeta -a --keyring-backend=test) 500000000000000000000000000000000stake --keyring-backend=test
zetacored add-genesis-account $(zetacored keys show mario -a --keyring-backend=test)  500000000000000000000000000000000stake --keyring-backend=test

zetacored gentx zeta 1000000000000000000000000stake --chain-id=localnet --keyring-backend=test

echo "Collecting genesis txs..."
zetacored collect-gentxs

echo "Validating genesis file..."
zetacored validate-genesis


mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
mkdir -p $DAEMON_HOME/cosmovisor/upgrades/$UpgradeName/bin


# Setup new binary
git checkout $NewBinary
rm -rf $GOPATH/bin/zetacored
make install
cp $GOPATH/bin/zetacored $GOPATH/bin/new/


# Setup cosmovisor
# Genesis
cp $GOPATH/bin/old/zetacored $DAEMON_HOME/cosmovisor/genesis/bin
cp $GOPATH/bin/zetaclientd $DAEMON_HOME/cosmovisor/genesis/bin

#Upgrades
cp $GOPATH/bin/new/zetacored $DAEMON_HOME/cosmovisor/upgrades/$UpgradeName/bin/

#Permissions
chmod +x $DAEMON_HOME/cosmovisor/genesis/bin/zetacored
chmod +x $DAEMON_HOME/cosmovisor/genesis/bin/zetaclientd
chmod +x $DAEMON_HOME/cosmovisor/upgrades/$UpgradeName/bin/zetacored

contents="$(jq '.app_state.gov.voting_params.voting_period = "10s"' $DAEMON_HOME/config/genesis.json)" && \
echo "${contents}" > $DAEMON_HOME/config/genesis.json

# Add state data here if required

cosmovisor start --home ~/.zetacore/ --p2p.laddr 0.0.0.0:27655  --grpc.address 0.0.0.0:9096 --grpc-web.address 0.0.0.0:9093 --address tcp://0.0.0.0:27659 --rpc.laddr tcp://127.0.0.1:26657 >> zetanode.log 2>&1  &

sleep 7
zetacored tx gov submit-proposal software-upgrade $UpgradeName --from zeta --deposit 100000000stake --upgrade-height 10 --title $UpgradeName --description $UpgradeName --keyring-backend test --chain-id localnet --yes
sleep 7
zetacored tx gov vote 1 yes --from zeta --keyring-backend test --chain-id localnet --yes
clear
sleep 7
zetacored query gov proposal 1

#tail -f zetanode.log

