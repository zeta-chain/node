#!/bin/bash

zetacored start --pruning=nothing   --log_level info --minimum-gas-prices=0.0001azeta --json-rpc.api eth,txpool,personal,net,debug,web3,miner --api.enable --home ~/.zetacored
