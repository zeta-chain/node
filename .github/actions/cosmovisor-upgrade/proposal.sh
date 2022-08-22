contents="$(jq '.app_state.gov.voting_params.voting_period = "10s"' "$DAEMON_HOME"/config/genesis.json)" && \
echo "${contents}" > "$DAEMON_HOME"/config/genesis.json

# Add state data here if required

# cosmovisor start --home ~/.zetacore/ --p2p.laddr 0.0.0.0:27655  --grpc.address 0.0.0.0:9096 --grpc-web.address 0.0.0.0:9093 --address tcp://0.0.0.0:27659 --rpc.laddr tcp://127.0.0.1:26657 >> zetanode.log 2>&1  &

sleep 7
zetacored tx gov submit-proposal software-upgrade "$UpgradeName" --from zeta --deposit 100000000stake --upgrade-height 10 --title "$UpgradeName" --description "$UpgradeName" --keyring-backend test --chain-id localnet --yes
sleep 7
zetacored tx gov vote 1 yes --from zeta --keyring-backend test --chain-id localnet --yes
clear
sleep 7
zetacored query gov proposal 1

tail -f zetanode.log

