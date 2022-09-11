#!/bin/bash

zetacored start --pruning=nothing   --minimum-gas-prices=0.0001azeta --json-rpc.api eth,txpool,personal,net,debug,web3,miner --api.enable --home ~/.zetacore
