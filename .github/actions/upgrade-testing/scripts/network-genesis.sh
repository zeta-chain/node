export PATH="/usr/local/go/bin:${PATH}"

echo 'export PATH=/usr/local/go/bin:'${PATH} >> /root/.bashrc
cat /root/.bashrc

log_it () {
  echo "********************************"
  echo "$1"
  echo "********************************"
  echo ""
  echo ""
}

chmod -R 777 /app_version
cp /app_version/${STARTING_VERSION}/bin/${DAEMON_NAME} /usr/bin/${DAEMON_NAME}
echo "${ZETA_MNEMONIC}" | ${DAEMON_NAME} keys add ${MONIKER} --keyring-backend test --recover
${DAEMON_NAME} init "${MONIKER}" --chain-id "${CHAIN_ID}"
cp /app_version/app.toml ${DAEMON_HOME}config/app.toml
cp /app_version/config.toml ${DAEMON_HOME}config/config.toml
mkdir -p ${DAEMON_HOME}/cosmovisor/genesis/bin
mkdir -p ${DAEMON_HOME}/cosmovisor/upgrades

genesis_account=$(${DAEMON_NAME} keys show ${MONIKER} -a --keyring-backend test)
log_it "${genesis_account}"

validator_account=$(${DAEMON_NAME} keys show ${MONIKER} -a --bech val --keyring-backend=test)
log_it "${validator_account}"

log_it "Add Genesis Account"
${DAEMON_NAME} add-genesis-account ${genesis_account} 500000000000000000000000000000000${DENOM}

log_it "GenerateTX"
${DAEMON_NAME} gentx ${MONIKER} 1000000000000000000000000${DENOM} --chain-id "${CHAIN_ID}" --ip "127.0.0.1" --keyring-backend test --moniker ${MONIKER}

log_it "Collect GenTX"
${DAEMON_NAME} collect-gentxs

log_it "Modify Genesis File"
echo $(jq --arg a "${VOTING_PERIOD}" '.app_state.gov.voting_params.voting_period = ($a)' ${DAEMON_HOME}/config/genesis.json) > ${DAEMON_HOME}/config/genesis.json
echo $(jq --arg a "${DENOM}" '.app_state.crisis.constant_fee.denom = ($a)' ${DAEMON_HOME}/config/genesis.json) > ${DAEMON_HOME}/config/genesis.json
echo $(jq --arg a "${DENOM}" '.app_state.mint.params.mint_denom = ($a)' ${DAEMON_HOME}/config/genesis.json) > ${DAEMON_HOME}/config/genesis.json
echo $(jq --arg a "${DENOM}" '.app_state.gov.deposit_params.min_deposit[0].denom = ($a)' ${DAEMON_HOME}/config/genesis.json) > ${DAEMON_HOME}/config/genesis.json
echo $(jq --arg a "${DENOM}" '.app_state.staking.params.bond_denom = ($a)' ${DAEMON_HOME}/config/genesis.json) > ${DAEMON_HOME}/config/genesis.json
echo $(jq --arg a "${DENOM}" '.app_state.evm.params.evm_denom = ($a)' ${DAEMON_HOME}/config/genesis.json) > ${DAEMON_HOME}/config/genesis.json

log_it "*********DEBUG GENSIS FILE***********"
cat ${DAEMON_HOME}/config/genesis.json | jq .app_state.gov.voting_params.voting_period
cat ${DAEMON_HOME}/config/genesis.json | jq .app_state.crisis.constant_fee.denom
cat ${DAEMON_HOME}/config/genesis.json | jq .app_state.mint.params.mint_denom
cat ${DAEMON_HOME}/config/genesis.json | jq .app_state.gov.deposit_params.min_deposit[0].denom
cat ${DAEMON_HOME}/config/genesis.json | jq .app_state.staking.params.bond_denom
cat ${DAEMON_HOME}/config/genesis.json | jq .app_state.evm.params.evm_denom
log_it "**************************************"


log_it "Copy Binaries to Cosmovisor Upgrades Folder"
cp -r /app_version/* ${DAEMON_HOME}/cosmovisor/upgrades/

log_it "***************"
log_it "Cosmos Upgrades"
ls -lah ${DAEMON_HOME}/cosmovisor/upgrades/

log_it "Copy Starting Binary to Cosmovisor Genesis Bin Folder"
cp /usr/bin/${DAEMON_NAME} ${DAEMON_HOME}/cosmovisor/genesis/bin

chmod -R 777 ${DAEMON_HOME}/cosmovisor
chmod -R a+x ${DAEMON_HOME}/cosmovisor/

log_it "Validate Genesis File"
${DAEMON_NAME} validate-genesis --home ${DAEMON_HOME}/

nohup cosmovisor start --rpc.laddr tcp://0.0.0.0:26657 --minimum-gas-prices ${GAS_PRICES} "--grpc.enable=true" > cosmovisor.log 2>&1 &
tail -n 1000 -f cosmovisor.log
