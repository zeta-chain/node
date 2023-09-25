import os
import requests
import json

os.environ['NODE'] = "http://127.0.0.1:26657"
CURRENT_HEIGHT = requests.get(f"{os.environ['NODE']}/status").json()["result"]["sync_info"]["latest_block_height"]
UPGRADE_HEIGHT = int(CURRENT_HEIGHT) + (
            int(os.environ['PROPOSAL_TIME_SECONDS']) / int(os.environ['BLOCK_TIME_SECONDS'])) + 20
github_file = open(os.environ["GITHUB_ENV"], "a+")
github_file.write(f"UPGRADE_HEIGHT={UPGRADE_HEIGHT}")
github_file.close()

proposal_json = {
    "messages": [
        {
            "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
            "authority": os.environ["GOV_ADDRESS"],
            "plan": {
                "name": os.environ['VERSION'],
                "time": "0001-01-01T00:00:00Z",
                "height": str(UPGRADE_HEIGHT).split('.')[0],
                "info": os.environ["UPGRADE_INFO"],
                "upgraded_client_state": None
            }
        }
    ],
    "metadata": os.environ["METADATA"],
    "deposit": os.environ["DEPOSIT"]
}

proposal_json = json.dumps(proposal_json)
write_gov_json = open("gov.json", "w")
write_gov_json.write(proposal_json)
write_gov_json.close()

# GOV_PROPOSAL = f"""zetacored tx gov submit-proposal gov.json \
# --from {os.environ['MONIKER']} \
# --chain-id "{os.environ['CHAINID']}" \
# --keyring-backend test \
# --node "{os.environ['NODE']}" \
# --gas=auto \
# --gas-adjustment=2 \
# --gas-prices={os.environ['GAS_PRICES']} \
# -y
# """

GOV_PROPOSAL = f"""zetacored tx gov submit-legacy-proposal software-upgrade "{os.environ['VERSION']}" \
    --from "{os.environ['MONIKER']}" \
    --deposit {os.environ["DEPOSIT"]} \
    --upgrade-height "{str(UPGRADE_HEIGHT).split('.')[0]}" \
    --upgrade-info '{os.environ["UPGRADE_INFO"]}' \
    --title "{os.environ['VERSION']}" \
    --description "Zeta Release {os.environ['VERSION']}" \
    --chain-id "{os.environ['CHAINID']}" \
    --node "{os.environ['NODE']}" \
    --keyring-backend test \
    --gas=auto \
    --gas-adjustment=2 \
    --gas-prices={os.environ['GAS_PRICES']} \
    -y \
    --no-validate"""

print(GOV_PROPOSAL)
