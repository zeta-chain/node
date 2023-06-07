export VERSION=$1
echo "********************RESTART VARS********************"
echo "VERSION: ${VERSION}"
echo "GAS_PRICES: ${GAS_PRICES}"
echo "DAEMON_HOME: ${DAEMON_HOME}"
echo "********************RESTART VARS********************"

source /root/.bashrc
cd ${DAEMON_HOME}

echo "COPY BINARY TO CURRENT ONE"
cp cosmovisor/upgrades/${VERSION}/bin/zetacored cosmovisor/genesis/bin/zetacored

echo "CHECK CURRENT BINARY"
ls -lah cosmovisor/genesis/bin/zetacored

echo "KILL ALL COSMOVISOR"
killall cosmovisor

echo "RESTART COSMOVISOR"
nohup cosmovisor start --rpc.laddr tcp://0.0.0.0:26657 --minimum-gas-prices ${GAS_PRICES} "--grpc.enable=true" > cosmovisor.log 2>&1 &

echo "SLEEP FOR 15 SECONDS"
sleep 15

echo "CHECK VERSION"
cosmovisor version