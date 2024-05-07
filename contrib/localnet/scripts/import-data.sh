#!/bin/bash
mkdir ~/genesis_export/
echo "Download Latest State Export"
LATEST_EXPORT_URL=$(curl https://snapshots.zetachain.com/latest-state-export | jq -r .mainnet)
echo "LATEST EXPORT URL: ${LATEST_EXPORT_URL}"
wget -q ${LATEST_EXPORT_URL} -O ~/genesis_export/exported-genesis.json