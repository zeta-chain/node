export VERSION=$1

echo "********************RESTART VARS********************" > /root/.zetacored/restart.log
echo "VERSION: ${VERSION}" >> /root/.zetacored/restart.log
echo "GAS_PRICES: ${GAS_PRICES}" >> /root/.zetacored/restart.log
echo "DAEMON_HOME: ${DAEMON_HOME}" >> /root/.zetacored/restart.log
echo "********************RESTART VARS********************" >> /root/.zetacored/restart.log

source /root/.bashrc

cd ${DAEMON_HOME}

echo "CHECK CURRENT BINARY" >> /root/.zetacored/restart.log
ls -lah cosmovisor/genesis/bin/zetacored >> /root/.zetacored/restart.log

echo "COPY BINARY TO CURRENT ONE" >> /root/.zetacored/restart.log
cp cosmovisor/upgrades/${VERSION}/bin/zetacored cosmovisor/genesis/bin/zetacored

echo "CHECK CURRENT BINARY" >> /root/.zetacored/restart.log
ls -lah cosmovisor/genesis/bin/zetacored >> /root/.zetacored/restart.log

echo "KILL ALL COSMOVISOR" >> /root/.zetacored/restart.log
killall cosmovisor

echo "RESTART COSMOVISOR" >> /root/.zetacored/restart.log
nohup cosmovisor start --rpc.laddr tcp://0.0.0.0:26657 --minimum-gas-prices ${GAS_PRICES} "--grpc.enable=true" > cosmovisor.log 2>&1 &

echo "SLEEP FOR 15 SECONDS" >> /root/.zetacored/restart.log
sleep 15

echo "CHECK VERSION" >> /root/.zetacored/restart.log
cosmovisor version >> /root/.zetacored/restart.log