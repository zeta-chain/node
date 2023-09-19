import json
import os

genesis = open(os.environ["NEW_GENESIS"], "r").read()
genesis_json_object = json.loads(genesis)

# genesis_json_object["staking"]["params"]["bond_denom"] = os.environ["DENOM"]
# genesis_json_object["crisis"]["constant_fee"]["denom"] = os.environ["DENOM"]
# genesis_json_object["gov"]["deposit_params"]["min_deposit"][0]["denom"] = os.environ["DENOM"]
# genesis_json_object["mint"]["params"]["mint_denom"] = os.environ["DENOM"]
# genesis_json_object["evm"]["params"]["evm_denom"] = os.environ["DENOM"]
# genesis_json_object["block"]["max_gas"] = os.environ["MAX_GAS"]
# genesis_json_object["gov"]["voting_params"]["voting_period"] = f'{os.environ["PROPOSAL_TIME_SECONDS"]}s'

"""
#Set config to use azeta
check cat $HOME/.zetacored/config/genesis.json | jq '.app_state["staking"]["params"]["bond_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["crisis"]["constant_fee"]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["mint"]["params"]["mint_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.app_state["evm"]["params"]["evm_denom"]="azeta"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
cat $HOME/.zetacored/config/genesis.json | jq '.consensus_params["block"]["max_gas"]="10000000"' > $HOME/.zetacored/config/tmp_genesis.json && mv $HOME/.zetacored/config/tmp_genesis.json $HOME/.zetacored/config/genesis.json
contents="$(jq '.app_state.gov.voting_params.voting_period = "10s"' $DAEMON_HOME/config/genesis.json)" && \
echo "${contents}" > $DAEMON_HOME/config/genesis.json
sed -i '/\[api\]/,+3 s/enable = false/enable = true/' ~/.zetacored/config/app.toml
"""

exported_genesis = open(os.environ["OLD_GENESIS"], "r").read()
exported_genesis_json_object = json.loads(exported_genesis)

crosschain = exported_genesis_json_object["app_state"]["crosschain"]
observer = exported_genesis_json_object["app_state"]["observer"]
emissions = exported_genesis_json_object["app_state"]["emissions"]
fungible = exported_genesis_json_object["app_state"]["fungible"]
evm = exported_genesis_json_object["app_state"]["evm"]
auth_accounts = exported_genesis_json_object["app_state"]["auth"]["accounts"]

genesis_json_object["app_state"]["auth"]["accounts"] = genesis_json_object["app_state"]["auth"]["accounts"] + auth_accounts
genesis_json_object["app_state"]["crosschain"] = crosschain
genesis_json_object["app_state"]["observer"] = observer
genesis_json_object["app_state"]["emissions"] = emissions
genesis_json_object["app_state"]["fungible"] = fungible

evm_accounts = []
for index, account in enumerate(evm["accounts"]):
    if account["address"] == "0x0000000000000000000000000000000000000001":
        print("pop account", account["address"])
    elif account["address"] == "0x0000000000000000000000000000000000000006":
        print("pop account", account["address"])
    elif account["address"] == "0x0000000000000000000000000000000000000002":
        print("pop account", account["address"])
    elif account["address"] == "0x0000000000000000000000000000000000000002":
        print("pop account", account["address"])
    elif account["address"] == "0x0000000000000000000000000000000000000008":
        print("pop account", account["address"])
    else:
        evm_accounts.append(account)

evm["accounts"] = evm_accounts
genesis_json_object["app_state"]["evm"] = evm

genesis = open("genesis-edited.json", "w")
genesis_string = json.dumps(genesis_json_object, indent=2)
dumped_genesis_object = genesis_string.replace("0x0000000000000000000000000000000000000001","0x387A12B28fe02DcAa467c6a1070D19B82F718Bb5")
genesis.write(genesis_string)
genesis.close()
