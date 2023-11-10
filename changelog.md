# CHANGELOG

## Unreleasd

### Features

### Fixes

### Refactoring

### Chores

### Tests

### CI

## Version: v10.1.2

### Features
* external stress test by @kevinssgh in https://github.com/zeta-chain/node/pull/1137
* allow setting liquidity cap for ZRC20 by @lumtis in https://github.com/zeta-chain/node/pull/1205
* add ability to update gas limit by @lumtis in https://github.com/zeta-chain/node/pull/1260
* Bitcoin block header and merkle proof by @ws4charlie in https://github.com/zeta-chain/node/pull/1263
* add query command to get all gas stability pool balances by @lumtis in https://github.com/zeta-chain/node/pull/1247
* tss funds migration by @kingpinXD in https://github.com/zeta-chain/node/pull/1143

### Fixes

* added upgrade name, and allow download. allows to test release can. by @gzukel in https://github.com/zeta-chain/node/pull/1195
* address `cosmos-gosec` lint issues by @lumtis in https://github.com/zeta-chain/node/pull/1153
* adding namespaces back by @lumtis in https://github.com/zeta-chain/node/pull/1228
* set unique index for generate cctx by @lumtis in https://github.com/zeta-chain/node/pull/1245
* remove error return in `IsAuthorized` by @lumtis in https://github.com/zeta-chain/node/pull/1250
* Ethereum comparaison checksum/non-checksum format by @lumtis in https://github.com/zeta-chain/node/pull/1261
* Blame index update by @kevinssgh in https://github.com/zeta-chain/node/pull/1264
* feed sataoshi/B to zetacore and check actual outTx size by @ws4charlie in https://github.com/zeta-chain/node/pull/1243
* cherry pick all hotfix from v10.0.x (zero-amount, precision, etc.) by @ws4charlie in https://github.com/zeta-chain/node/pull/1235
* register emissions grpc server by @kingpinXD in https://github.com/zeta-chain/node/pull/1257
* read gas limit from smart contract by @lumtis in https://github.com/zeta-chain/node/pull/1277
* add CLI command to query system contract by @lumtis in https://github.com/zeta-chain/node/pull/1252
* add notice when using `--ledger` with Ethereum HD path by @lumtis in https://github.com/zeta-chain/node/pull/1285
* gosec issues by @lumtis in https://github.com/zeta-chain/node/pull/1290
* query outtx tracker by chain using prefixed store by @ws4charlie in https://github.com/zeta-chain/node/pull/1283
* minor fixes to stateful upgrade by @kevinssgh in https://github.com/zeta-chain/node/pull/1280
* remove check `gasObtained == outTxGasFee` by @lumtis in https://github.com/zeta-chain/node/pull/1304
* begin blocker for mock mainnet by @kingpinXD in https://github.com/zeta-chain/node/pull/1308

### Refactoring

* call `onCrossChainCall` when depositing to a contract by @lumtis in https://github.com/zeta-chain/node/pull/1226
* change default mempool version in config by @lumtis in https://github.com/zeta-chain/node/pull/1238
* remove duplicate funtion name IsEthereum by @lukema95 in https://github.com/zeta-chain/node/pull/1279
* skip gas stability pool funding when gasLimit is equal gasUsed by @lukema95 in https://github.com/zeta-chain/node/pull/1289

### Chores

* switch back to `cosmos/cosmos-sdk` by @lumtis in https://github.com/zeta-chain/node/pull/1193
* changed maxNestedMsgs by @CharlieMc0 in https://github.com/zeta-chain/node/pull/1222
* sync from mockmain  by @brewmaster012 in https://github.com/zeta-chain/node/pull/1265
* increment handler version by @kingpinXD in https://github.com/zeta-chain/node/pull/1307

### Tests

* Stateful upgrade by @kevinssgh in https://github.com/zeta-chain/node/pull/1135

### CI

* cross-compile release binaries and simplify PR testings by @CharlieMc0 in https://github.com/zeta-chain/node/pull/1218
* add mainnet builds to goreleaser by @CharlieMc0 in https://github.com/zeta-chain/node/pull/1302







