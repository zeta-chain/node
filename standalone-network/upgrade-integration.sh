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

KEYRING=test
CHAINID="localnet_101-1"
export DAEMON_HOME=$HOME/.zetacored
export DAEMON_NAME=zetacored
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
export DAEMON_RESTART_AFTER_UPGRADE=true
export CLIENT_DAEMON_NAME=zetaclientd
export CLIENT_DAEMON_ARGS="-enable-chains,GOERLI,-val,zeta"
export DAEMON_DATA_BACKUP_DIR=$DAEMON_HOME
export CLIENT_SKIP_UPGRADE=true
export CLIENT_START_PROCESS=false
export UNSAFE_SKIP_BACKUP=true

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
git pull
make install-zetacore
cp $GOPATH/bin/zetacored $GOPATH/bin/old/
zetacored init test --chain-id=localnet_101-1 -o
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
echo "Generating deterministic account - zeta"
echo "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" | zetacored keys add zeta --algo secp256k1 --recover --keyring-backend=test
echo "Generating deterministic account - mario"
echo "hand inmate canvas head lunar naive increase recycle dog ecology inhale december wide bubble hockey dice worth gravity ketchup feed balance parent secret orchard" | zetacored keys add mario --algo secp256k1 --recover --keyring-backend=test



zetacored add-observer-list standalone-network/observers.json --keygen-block=0 --tss-pubkey="tsspubkey"
zetacored gentx zeta 1000000000000000000000azeta --chain-id=$CHAINID --keyring-backend=$KEYRING

echo "Collecting genesis txs..."
zetacored collect-gentxs

echo "Validating genesis file..."
zetacored validate-genesis


mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin
mkdir -p $DAEMON_HOME/cosmovisor/upgrades/$UpgradeName/bin


# Setup new binary
git checkout $NewBinary
git pull
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

cosmovisor start --home ~/.zetacored/ --p2p.laddr 0.0.0.0:27655  --grpc.address 0.0.0.0:9096 --grpc-web.address 0.0.0.0:9093 --address tcp://0.0.0.0:27659 --rpc.laddr tcp://127.0.0.1:26657 >> zetanode.log 2>&1  &
sleep 8
zetacored tx gov submit-legacy-proposal software-upgrade $UpgradeName --from zeta --deposit 100000000azeta --upgrade-height 6 --title $UpgradeName --description $UpgradeName --keyring-backend test --chain-id localnet_101-1 --yes --no-validate --fees=200azeta --broadcast-mode block
sleep 8
zetacored tx gov vote 1 yes --from zeta --keyring-backend test --chain-id localnet_101-1 --yes --fees=200azeta --broadcast-mode block
clear
sleep 7
zetacored query gov proposal 1

tail -f zetanode.log

