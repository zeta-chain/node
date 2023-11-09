# CHANGELOG

## Version: v10.1.0
### What's Changed:
* fix: updated release ci
* fix: added upgrade name, and allow download. allows to test release can. by @gzukel in https://github.com/zeta-chain/node/pull/1195
* chore: switch back to `cosmos/cosmos-sdk` by @lumtis in https://github.com/zeta-chain/node/pull/1193
* test: Stateful upgrade by @kevinssgh in https://github.com/zeta-chain/node/pull/1135
* fix: address `cosmos-gosec` lint issues by @lumtis in https://github.com/zeta-chain/node/pull/1153
* chore: changed maxNestedMsgs by @CharlieMc0 in https://github.com/zeta-chain/node/pull/1222
* ci: cross-compile release binaries and simplify PR testings by @CharlieMc0 in https://github.com/zeta-chain/node/pull/1218
* feat: external stress test by @kevinssgh in https://github.com/zeta-chain/node/pull/1137
* feat: allow setting liquidity cap for ZRC20 by @lumtis in https://github.com/zeta-chain/node/pull/1205
* refactor: call `onCrossChainCall` when depositing to a contract by @lumtis in https://github.com/zeta-chain/node/pull/1226
* fix(`rpc`): adding namespaces back by @lumtis in https://github.com/zeta-chain/node/pull/1228
* refactor(`cmd`): change default mempool version in config by @lumtis in https://github.com/zeta-chain/node/pull/1238
* fix(`MsgWhitelistERC20`): set unique index for generate cctx by @lumtis in https://github.com/zeta-chain/node/pull/1245
* fix(`observer`): remove error return in `IsAuthorized` by @lumtis in https://github.com/zeta-chain/node/pull/1250
* fix(`GetForeignCoinFromAsset`): Ethereum comparaison checksum/non-checksum format by @lumtis in https://github.com/zeta-chain/node/pull/1261
* feat(`fungible`): add ability to update gas limit by @lumtis in https://github.com/zeta-chain/node/pull/1260
* fix: Blame index update by @kevinssgh in https://github.com/zeta-chain/node/pull/1264
* fix: feed sataoshi/B to zetacore and check actual outTx size by @ws4charlie in https://github.com/zeta-chain/node/pull/1243
* fix: cherry pick all hotfix from v10.0.x (zero-amount, precision, etc.) by @ws4charlie in https://github.com/zeta-chain/node/pull/1235
* fix: register emissions grpc server by @kingpinXD in https://github.com/zeta-chain/node/pull/1257
* feat: Bitcoin block header and merkle proof by @ws4charlie in https://github.com/zeta-chain/node/pull/1263
* fix: read gas limit from smart contract by @lumtis in https://github.com/zeta-chain/node/pull/1277
* fix(`fungible`): add CLI command to query system contract by @lumtis in https://github.com/zeta-chain/node/pull/1252
* fix(`cmd`): add notice when using `--ledger` with Ethereum HD path by @lumtis in https://github.com/zeta-chain/node/pull/1285
* fix: gosec issues by @lumtis in https://github.com/zeta-chain/node/pull/1290
* refactor: remove duplicate funtion name IsEthereum by @lukema95 in https://github.com/zeta-chain/node/pull/1279
* fix: query outtx tracker by chain using prefixed store by @ws4charlie in https://github.com/zeta-chain/node/pull/1283
* refactor: skip gas stability pool funding when gasLimit is equal gasUsed by @lukema95 in https://github.com/zeta-chain/node/pull/1289
* fix: minor fixes to stateful upgrade by @kevinssgh in https://github.com/zeta-chain/node/pull/1280
* ci: add mainnet builds to goreleaser by @CharlieMc0 in https://github.com/zeta-chain/node/pull/1302
* feat(`fungible`): add query command to get all gas stability pool balances by @lumtis in https://github.com/zeta-chain/node/pull/1247
* feat: tss funds migration by @kingpinXD in https://github.com/zeta-chain/node/pull/1143
* fix(`gas-payment`): remove check `gasObtained == outTxGasFee` by @lumtis in https://github.com/zeta-chain/node/pull/1304
* chore: sync from mockmain  by @brewmaster012 in https://github.com/zeta-chain/node/pull/1265
* chore: increment handler version by @kingpinXD in https://github.com/zeta-chain/node/pull/1307
* fix: begin blocker for mock mainnet by @kingpinXD in https://github.com/zeta-chain/node/pull/1308



