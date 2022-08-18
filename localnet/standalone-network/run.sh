#!/usr/bin/env bash

killall zetacored
zetacored start --trace \
--home ~/.zetacored \
--p2p.laddr 0.0.0.0:27655  \
--grpc.address 0.0.0.0:9096 \
--grpc-web.address 0.0.0.0:9093 \
--address tcp://0.0.0.0:27659 \
--rpc.laddr tcp://127.0.0.1:26657 \
--pruning custom \
--pruning-keep-recent 1 \
--pruning-keep-every 10  \
--pruning-interval 10 \
--state-sync.snapshot-interval 10 \
--state-sync.snapshot-keep-recent 1
