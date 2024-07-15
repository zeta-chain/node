#!/bin/bash

CHAINID="athens_101-1"
UPGRADE_AUTHORITY_ACCOUNT="zeta10d07y265gmmuvt4z0w9aw880jnsr700jvxasvr"

if [[ -z $ZETACORED_URL ]]; then
    ZETACORED_URL='http://upgrade-host:8000/zetacored'
fi
if [[ -z $ZETACLIENTD_URL ]]; then
    ZETACLIENTD_URL='http://upgrade-host:8000/zetaclientd'
fi

# Wait for authorized_keys file to exist (populated by zetacore0)
while [ ! -f ~/.ssh/authorized_keys ]; do
    echo "Waiting for authorized_keys file to exist..."
    sleep 1
done

while ! curl -s -o /dev/null zetacore0:26657/status ; do
    echo "Waiting for zetacore0 rpc"
    sleep 1
done

# wait for minimum height
CURRENT_HEIGHT=0
while [[ $CURRENT_HEIGHT -lt 1 ]]
do
    CURRENT_HEIGHT=$(curl -s zetacore0:26657/status | jq -r '.result.sync_info.latest_block_height')
    echo "current height is ${CURRENT_HEIGHT}, waiting for 1"
    sleep 1
done

# copy zetacore0 config and keys if not running on zetacore0
if [[ $(hostname) != "zetacore0" ]]; then
  scp -r 'zetacore0:~/.zetacored/config' 'zetacore0:~/.zetacored/os_info' 'zetacore0:~/.zetacored/config' 'zetacore0:~/.zetacored/keyring-file' 'zetacore0:~/.zetacored/keyring-test' ~/.zetacored/
  sed -i 's|tcp://localhost:26657|tcp://zetacore0:26657|g' ~/.zetacored/config/client.toml
fi

# get new zetacored version
curl -L -o /tmp/zetacored.new "${ZETACORED_URL}"
chmod +x /tmp/zetacored.new
UPGRADE_NAME=$(/tmp/zetacored.new version)

# if explicit upgrade height not provided, use dumb estimator
if [[ -z $UPGRADE_HEIGHT ]]; then
    UPGRADE_HEIGHT=$(( $(curl -s zetacore0:26657/status | jq '.result.sync_info.latest_block_height' | tr -d '"') + 60))
    echo "Upgrade height was not provided. Estimating ${UPGRADE_HEIGHT}."
fi

cat > upgrade.json <<EOF
{
  "messages": [
    {
      "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
      "plan": {
        "height": "${UPGRADE_HEIGHT}",
        "info": "",
        "name": "${UPGRADE_NAME}",
        "time": "0001-01-01T00:00:00Z",
        "upgraded_client_state": null
      },
      "authority": "${UPGRADE_AUTHORITY_ACCOUNT}"
    }
  ],
  "metadata": "",
  "deposit": "1000000000000000000000azeta",
  "title": "${UPGRADE_NAME}",
  "summary": "${UPGRADE_NAME}"
}
EOF

# convert uname arch to goarch style
UNAME_ARCH=$(uname -m)
case "$UNAME_ARCH" in
    x86_64)    GOARCH=amd64;;
    i686)      GOARCH=386;;
    armv7l)    GOARCH=arm;;
    aarch64)   GOARCH=arm64;;
    *)         GOARCH=unknown;;
esac

cat > upgrade_plan_info.json <<EOF
{
    "binaries": {
        "linux/${GOARCH}": "${ZETACORED_URL}",
        "zetaclientd-linux/${GOARCH}": "${ZETACLIENTD_URL}"
    }
}
EOF

cat upgrade.json | jq --arg info "$(cat upgrade_plan_info.json)" '.messages[0].plan.info = $info' | tee upgrade_full.json

echo "Submitting upgrade proposal"

zetacored tx gov submit-proposal upgrade_full.json --from operator --keyring-backend test --chain-id $CHAINID --yes --fees 2000000000000000azeta -o json | tee proposal.json
PROPOSAL_TX_HASH=$(jq -r .txhash proposal.json)
PROPOSAL_ID=""
# WARN: this seems to be unstable
while [[ -z $PROPOSAL_ID ]]; do
    echo "waiting to get proposal_id"
    sleep 1
    # v0.47 version
    # proposal_id=$(zetacored query tx $PROPOSAL_TX_HASH -o json | jq -r '.events[] | select(.type == "submit_proposal") | .attributes[] | select(.key == "proposal_id") | .value')
    
    # v0.46 version
    PROPOSAL_ID=$(zetacored query tx $PROPOSAL_TX_HASH -o json | jq -r '.logs[0].events[] | select(.type == "proposal_deposit") | .attributes[] | select(.key == "proposal_id") | .value')
done
echo "proposal id is ${PROPOSAL_ID}"

zetacored tx gov vote "${PROPOSAL_ID}" yes --from operator --keyring-backend test --chain-id $CHAINID --yes --fees=2000000000000000azeta