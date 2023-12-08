# CHANGELOG

## Unreleased

### Features
* [1395](https://github.com/zeta-chain/node/pull/1395) - Add state variable to track aborted zeta amount
* [1387](https://github.com/zeta-chain/node/pull/1387) - Add HSM capability for zetaclient hot key
* enable zetaclients to use dynamic gas price on zetachain - enables >0 min_gas_price in feemarket module

### Fixes
* 6582f6b42f4f0eb0358e6fdefe5278295c791166 - masked zetaclient config at startup
* 1b8fdd3394c35afa47e830a5858b3c5c133fad8e - Added check for redeployment of gas and asset token contracts
* [1372](https://github.com/zeta-chain/node/pull/1372) - Include Event Index as part for inbound tx digest
* [1367](https://github.com/zeta-chain/node/pull/1367) - fix minRelayTxFee issue and check misuse of bitcoin mainnet/testnet addresses
* [1358](https://github.com/zeta-chain/node/pull/1358) - add a new thread to zetaclient which checks zeta supply in all connected chains in every block
* prevent deposits for paused zrc20
* [1406](https://github.com/zeta-chain/node/pull/1406) - improve log prints and speed up evm outtx inclusion
* fix Athens-3 issue - include bitcoin outtx regardless of the cctx status
* add check to verify new tss has been produced when triggering tss funds migration

* fix Athens-3 log print issue - avoid posting uncessary outtx confirmation


### Refactoring
* [1211](https://github.com/zeta-chain/node/issues/1211) - use `grpc` and `msg` for query and message files
* refactor cctx scheduler - decouple evm cctx scheduler from btc cctx scheduler
* [1391](https://github.com/zeta-chain/node/pull/1391) - consolidate node builds
* update `MsgUpdateContractBytecode` to use code hash instead of contract address
* move tss state date from crosschain to observer
### Chores

### Tests
- Add unit tests for adding votes to a ballot 
### CI

## Version: v10.1.2
### Unreleased:
* add a new thread to zetaclient which checks zeta supply in all connected chains in every block
* add a new tx to update an observer, this can be either be run a tombstoned observer/validator or via admin_policy_group_2.

### Features
* [1137](https://github.com/zeta-chain/node/pull/1137) - external stress testing
* [1205](https://github.com/zeta-chain/node/pull/1205) - allow setting liquidity cap for ZRC20
* [1260](https://github.com/zeta-chain/node/pull/1260) - add ability to update gas limit
* [1263](https://github.com/zeta-chain/node/pull/1263) - Bitcoin block header and merkle proof
* [1247](https://github.com/zeta-chain/node/pull/1247) - add query command to get all gas stability pool balances
* [1143](https://github.com/zeta-chain/node/pull/1143) - tss funds migration capability
* [1358](https://github.com/zeta-chain/node/pull/1358) - zetaclient thread for zeta supply checks
* [1384](https://github.com/zeta-chain/node/pull/1384) - tx to update an observer
### Fixes

* [1195](https://github.com/zeta-chain/node/pull/1195) - added upgrade name, and allow download. allows to test release
* [1153](https://github.com/zeta-chain/node/pull/1153) - address `cosmos-gosec` lint issues
* [1128](https://github.com/zeta-chain/node/pull/1228) - adding namespaces back in rpc
* [1245](https://github.com/zeta-chain/node/pull/1245) - set unique index for generate cctx
* [1250](https://github.com/zeta-chain/node/pull/1250) - remove error return in `IsAuthorized`
* [1261](https://github.com/zeta-chain/node/pull/1261) - Ethereum comparaison checksum/non-checksum format
* [1264](https://github.com/zeta-chain/node/pull/1264) - Blame index update
* [1243](https://github.com/zeta-chain/node/pull/1243) - feed sataoshi/B to zetacore and check actual outTx size
* [1235](https://github.com/zeta-chain/node/pull/1235) - cherry pick all hotfix from v10.0.x (zero-amount, precision, etc.)
* [1257](https://github.com/zeta-chain/node/pull/1257) - register emissions grpc server
* [1277](https://github.com/zeta-chain/node/pull/1277) - read gas limit from smart contract
* [1252](https://github.com/zeta-chain/node/pull/1252) - add CLI command to query system contract
* [1285](https://github.com/zeta-chain/node/pull/1285) - add notice when using `--ledger` with Ethereum HD path
* [1283](https://github.com/zeta-chain/node/pull/1283) - query outtx tracker by chain using prefixed store
* [1280](https://github.com/zeta-chain/node/pull/1280) - minor fixes to stateful upgrade
* [1304](https://github.com/zeta-chain/node/pull/1304) - remove check `gasObtained == outTxGasFee`
* [1308](https://github.com/zeta-chain/node/pull/1308) - begin blocker for mock mainnet

### Refactoring

* [1226](https://github.com/zeta-chain/node/pull/1226) - call `onCrossChainCall` when depositing to a contract
* [1238](https://github.com/zeta-chain/node/pull/1238) - change default mempool version in config 
* [1279](https://github.com/zeta-chain/node/pull/1279) - remove duplicate funtion name IsEthereum
* [1289](https://github.com/zeta-chain/node/pull/1289) - skip gas stability pool funding when gasLimit is equal gasUsed

### Chores

* [1193](https://github.com/zeta-chain/node/pull/1193) - switch back to `cosmos/cosmos-sdk`
* [1222](https://github.com/zeta-chain/node/pull/1222) - changed maxNestedMsgs
* [1265](https://github.com/zeta-chain/node/pull/1265) - sync from mockmain
* [1307](https://github.com/zeta-chain/node/pull/1307) - increment handler version

### Tests

* [1135](https://github.com/zeta-chain/node/pull/1135) - Stateful upgrade for smoke tests

### CI

* [1218](https://github.com/zeta-chain/node/pull/1218) - cross-compile release binaries and simplify PR testings
* [1302](https://github.com/zeta-chain/node/pull/1302) - add mainnet builds to goreleaser

