# CHANGELOG

## Unreleased

* `zetaclientd start` : 2 inputs required from stdin

### Refactor

* [1630](https://github.com/zeta-chain/node/pull/1630) added password prompts for hotkey and tss keyshare in zetaclient
  Starting zetaclient now requires two passwords to be input; one for the hotkey and another for the tss key-share.
* [1731](https://github.com/zeta-chain/node/pull/1731) added doc for hotkey and tss key-share password prompts.

### Fixes

* [1690](https://github.com/zeta-chain/node/issues/1690) - double watched gas prices and fix btc scheduler
* [1687](https://github.com/zeta-chain/node/pull/1687) - only use EVM supported chains for gas stability pool
* [1692](https://github.com/zeta-chain/node/pull/1692) - fix get params query for emissions module

### Tests

* [1584](https://github.com/zeta-chain/node/pull/1584) - allow to run E2E tests on any networks

## Version: v12.2.4

### Fixes

* [1638](https://github.com/zeta-chain/node/issues/1638) - additional check to make sure external chain height always increases
* [1672](https://github.com/zeta-chain/node/pull/1672) - paying 50% more than base gas price to buffer EIP1559 gas price increase
* [1642](https://github.com/zeta-chain/node/pull/1642) - Change WhitelistERC20 authorization from group1 to group2
* [1610](https://github.com/zeta-chain/node/issues/1610) - add pending outtx hash to tracker after monitoring for 10 minutes
* [1656](https://github.com/zeta-chain/node/issues/1656) - schedule bitcoin keysign with intervals to avoid keysign failures
* [1662](https://github.com/zeta-chain/node/issues/1662) - skip Goerli BlobTxType transactions introduced in Dencun upgrade
* [1663](https://github.com/zeta-chain/node/issues/1663) - skip Mumbai empty block if ethclient sanity check fails
* [1661](https://github.com/zeta-chain/node/issues/1661) - use estimated SegWit tx size for Bitcoin gas fee calculation
* [1667](https://github.com/zeta-chain/node/issues/1667) - estimate SegWit tx size in uinit of vByte

## Version: v12.1.0

### Tests

* [1577](https://github.com/zeta-chain/node/pull/1577) - add chain header tests in E2E tests and fix admin tests

### Features
* [1658](https://github.com/zeta-chain/node/pull/1658) - modify emission distribution to use fixed block rewards

### Fixes
* [1535](https://github.com/zeta-chain/node/issues/1535) - Avoid voting on wrong ballots due to false blockNumber in EVM tx receipt
* [1588](https://github.com/zeta-chain/node/pull/1588) - fix chain params comparison logic
* [1650](https://github.com/zeta-chain/node/pull/1605) - exempt (discounted) *system txs* from min gas price check and gas fee deduction
* [1632](https://github.com/zeta-chain/node/pull/1632) - set keygen to `KeygenStatus_KeyGenSuccess` if its in `KeygenStatus_PendingKeygen`.
* [1576](https://github.com/zeta-chain/node/pull/1576) - Fix zetaclient crash due to out of bound integer conversion and log prints.
* [1575](https://github.com/zeta-chain/node/issues/1575) - Skip unsupported chain parameters by IsSupported flag

### CI

* [1580](https://github.com/zeta-chain/node/pull/1580) - Fix release pipelines cleanup step.

### Chores

* [1585](https://github.com/zeta-chain/node/pull/1585) - Updated release instructions
* [1615](https://github.com/zeta-chain/node/pull/1615) - Add upgrade handler for version v12.1.0

### Features

* [1591](https://github.com/zeta-chain/node/pull/1591) - support lower gas limit for voting on inbound and outbound transactions
* [1592](https://github.com/zeta-chain/node/issues/1592) - check inbound tracker tx hash against Tss address and some refactor on inTx observation

### Refactoring

* [1628](https://github.com/zeta-chain/node/pull/1628) optimize return and simplify code

### Refactoring
* [1619](https://github.com/zeta-chain/node/pull/1619) - Add evm fee calculation to tss migration of evm chains

## Version: v12.0.0

### Breaking Changes

TSS and chain validation related queries have been moved from `crosschain` module to `observer` module:
* `PendingNonces` :Changed from `/zeta-chain/crosschain/pendingNonces/{chain_id}/{address}` to `/zeta-chain/observer/pendingNonces/{chain_id}/{address}` . It returns all the pending nonces for a chain id and address. This returns the current pending nonces for the chain.
* `ChainNonces` : Changed from `/zeta-chain/crosschain/chainNonces/{chain_id}` to`/zeta-chain/observer/chainNonces/{chain_id}` . It returns all the chain nonces for a chain id. This returns the current nonce of the TSS address for the chain.
* `ChainNoncesAll` :Changed from `/zeta-chain/crosschain/chainNonces` to `/zeta-chain/observer/chainNonces` . It returns all the chain nonces for all chains. This returns the current nonce of the TSS address for all chains.

All chains now have the same observer set:
* `ObserversByChain`: `/zeta-chain/observer/observers_by_chain/{observation_chain}` has been removed and replaced with `/zeta-chain/observer/observer_set`. All chains have the same observer set.
* `AllObserverMappers`: `/zeta-chain/observer/all_observer_mappers` has been removed. `/zeta-chain/observer/observer_set` should be used to get observers.

Observer params and core params have been merged into chain params:
* `Params`: `/zeta-chain/observer/params` no longer returns observer params. Observer params data have been moved to chain params described below.
* `GetCoreParams`: Renamed into `GetChainParams`. `/zeta-chain/observer/get_core_params` moved to `/zeta-chain/observer/get_chain_params`.
* `GetCoreParamsByChain`: Renamed into `GetChainParamsForChain`. `/zeta-chain/observer/get_core_params_by_chain` moved to `/zeta-chain/observer/get_chain_params_by_chain`.

Getting the correct TSS address for Bitcoin now requires proviidng the Bitcoin chain id:
* `GetTssAddress` : Changed from `/zeta-chain/observer/get_tss_address/` to `/zeta-chain/observer/getTssAddress/{bitcoin_chain_id}` . Optional bitcoin chain id can now be passed as a parameter to fetch the correct tss for required BTC chain. This parameter only affects the BTC tss address in the response.

### Features
* [1498](https://github.com/zeta-chain/node/pull/1498) - Add monitoring(grafana, prometheus, ethbalance) for localnet testing
* [1395](https://github.com/zeta-chain/node/pull/1395) - Add state variable to track aborted zeta amount
* [1410](https://github.com/zeta-chain/node/pull/1410) - `snapshots` commands
* enable zetaclients to use dynamic gas price on zetachain - enables >0 min_gas_price in feemarket module
* add static chain data for Sepolia testnet
* added metrics to track the burn rate of the hotkey in the telemetry server as well as prometheus

### Fixes

* [1554](https://github.com/zeta-chain/node/pull/1554) - Screen out unconfirmed UTXOs that are not created by TSS itself
* [1560](https://github.com/zeta-chain/node/issues/1560) - Zetaclient post evm-chain outtx hashes only when receipt is available
* [1516](https://github.com/zeta-chain/node/issues/1516) - Unprivileged outtx tracker removal
* [1537](https://github.com/zeta-chain/node/issues/1537) - Sanity check events of ZetaSent/ZetaReceived/ZetaRevertedWithdrawn/Deposited
* [1530](https://github.com/zeta-chain/node/pull/1530) - Outbound tx confirmation/inclusion enhancement
* [1496](https://github.com/zeta-chain/node/issues/1496) - post block header for enabled EVM chains only
* [1518](https://github.com/zeta-chain/node/pull/1518) - Avoid duplicate keysign if an outTx is already pending
* fix Code4rena issue - zetaclients potentially miss inTx when PostSend (or other RPC) fails
* fix go-staticcheck warnings for zetaclient
* fix Athens-3 issue - incorrect pending-tx inclusion and incorrect confirmation count
* masked zetaclient config at startup
* set limit for queried pending cctxs
* add check to verify new tss has been produced when triggering tss funds migration
* fix Athens-3 log print issue - avoid posting uncessary outtx confirmation
* fix docker build issues with version: golang:1.20-alpine3.18
* [1525](https://github.com/zeta-chain/node/pull/1525) - relax EVM chain block header length check 1024->4096
* [1522](https://github.com/zeta-chain/node/pull/1522/files) - block `distribution` module account from receiving zeta
* [1528](https://github.com/zeta-chain/node/pull/1528) - fix panic caused on decoding malformed BTC addresses
* [1536](https://github.com/zeta-chain/node/pull/1536) - add index to check previously finalized inbounds
* [1556](https://github.com/zeta-chain/node/pull/1556) - add emptiness check for topic array in event parsing
* [1546](https://github.com/zeta-chain/node/pull/1546) - fix reset of pending nonces on genesis import
* [1555](https://github.com/zeta-chain/node/pull/1555) - Reduce websocket message limit to 10MB
* [1567](https://github.com/zeta-chain/node/pull/1567) - add bitcoin chain id to fetch the tss address rpc endpoint
* [1501](https://github.com/zeta-chain/node/pull/1501) - fix stress test - use new refactored config file and smoketest runner
* [1589](https://github.com/zeta-chain/node/pull/1589) - add bitcoin chain id to `get tss address` and `get tss address historical` cli query

### Refactoring

* [1552](https://github.com/zeta-chain/node/pull/1552) - requires group2 to enable header verification
* [1211](https://github.com/zeta-chain/node/issues/1211) - use `grpc` and `msg` for query and message files
* refactor cctx scheduler - decouple evm cctx scheduler from btc cctx scheduler
* move tss state from crosschain to observer
* move pending nonces, chain nonces and nonce to cctx to observer
* move tss related cli from crosschain to observer
* reorganize smoke tests structure
* Add pagination to queries which iterate over large data sets InTxTrackerAll ,PendingNoncesAll ,AllBlameRecord ,TssHistory
* GetTssAddress now returns only the current tss address for ETH and BTC
* Add a new query GetTssAddressesByFinalizedBlockHeight to get any other tss addresses for a finalized block height
* Move observer params into core params
* Remove chain id from the index for observer mapper and rename it to observer set.
* Add logger to smoke tests
* [1521](https://github.com/zeta-chain/node/pull/1521) - replace go-tss lib version with one that reverts back to thorchain tss-lib
* [1558](https://github.com/zeta-chain/node/pull/1558) - change log level for gas stability pool iteration error
* Update --ledger flag hint

### Chores
* [1446](https://github.com/zeta-chain/node/pull/1446) - renamed file `zetaclientd/aux.go` to `zetaclientd/utils.go` to avoid complaints from go package resolver. 
* [1499](https://github.com/zeta-chain/node/pull/1499) - Add scripts to localnet to help test gov proposals
* [1442](https://github.com/zeta-chain/node/pull/1442) - remove build types in `.goreleaser.yaml`
* [1504](https://github.com/zeta-chain/node/pull/1504) - remove `-race` in the `make install` commmand
*  [1564](https://github.com/zeta-chain/node/pull/1564) - bump ti-actions/changed-files

### Tests

* [1538](https://github.com/zeta-chain/node/pull/1538) - improve stateful e2e testing

### CI
* Removed private runners and unused GitHub Action

## Version: v11.0.0

### Features

* [1387](https://github.com/zeta-chain/node/pull/1387) - Add HSM capability for zetaclient hot key
* add a new thread to zetaclient which checks zeta supply in all connected chains in every block
* add a new tx to update an observer, this can be either be run a tombstoned observer/validator or via admin_policy_group_2.

### Fixes

* Added check for redeployment of gas and asset token contracts
* [1372](https://github.com/zeta-chain/node/pull/1372) - Include Event Index as part for inbound tx digest
* [1367](https://github.com/zeta-chain/node/pull/1367) - fix minRelayTxFee issue and check misuse of bitcoin mainnet/testnet addresses
* [1358](https://github.com/zeta-chain/node/pull/1358) - add a new thread to zetaclient which checks zeta supply in all connected chains in every block
* prevent deposits for paused zrc20
* [1406](https://github.com/zeta-chain/node/pull/1406) - improve log prints and speed up evm outtx inclusion
* fix Athens-3 issue - include bitcoin outtx regardless of the cctx status

### Refactoring

* [1391](https://github.com/zeta-chain/node/pull/1391) - consolidate node builds
* update `MsgUpdateContractBytecode` to use code hash instead of contract address

### Chores

### Tests
- Add unit tests for adding votes to a ballot 

### CI

## Version: v10.1.2

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

