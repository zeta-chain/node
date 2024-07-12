# CHANGELOG

## Unreleased

### Breaking Changes

* [2460](https://github.com/zeta-chain/node/pull/2460) - Upgrade to go 1.22. This required us to temporarily remove the QUIC backend from [go-libp2p](https://github.com/libp2p/go-libp2p). If you are a zetaclient operator and have configured quic peers, you need to switch to tcp peers.

### Features

* [2032](https://github.com/zeta-chain/node/pull/2032) - improve some general structure of the ZetaClient codebase
* [2100](https://github.com/zeta-chain/node/pull/2100) - cosmos v0.47 upgrade
* [2145](https://github.com/zeta-chain/node/pull/2145) - add `ibc` and `ibc-transfer` modules
* [2135](https://github.com/zeta-chain/node/pull/2135) - add develop build version logic
* [2152](https://github.com/zeta-chain/node/pull/2152) - custom priority nonce mempool
* [2113](https://github.com/zeta-chain/node/pull/2113) - add zetaclientd-supervisor process
* [2154](https://github.com/zeta-chain/node/pull/2154) - add `ibccrosschain` module
* [2282](https://github.com/zeta-chain/node/pull/2282) - modify rpc methods to support synthetic txs
* [2258](https://github.com/zeta-chain/node/pull/2258) - add Optimism and Base in static chain information
* [2287](https://github.com/zeta-chain/node/pull/2287) - implement `MsgUpdateChainInfo` message
* [2279](https://github.com/zeta-chain/node/pull/2279) - add a CCTXGateway field to chain static data
* [2275](https://github.com/zeta-chain/node/pull/2275) - add ChainInfo singleton state variable in authority
* [2291](https://github.com/zeta-chain/node/pull/2291) - initialize cctx gateway interface
* [2289](https://github.com/zeta-chain/node/pull/2289) - add an authorization list to keep track of all authorizations on the chain
* [2305](https://github.com/zeta-chain/node/pull/2305) - add new messages `MsgAddAuthorization` and `MsgRemoveAuthorization` that can be used to update the authorization list
* [2313](https://github.com/zeta-chain/node/pull/2313) - add `CheckAuthorization` function to replace the `IsAuthorized` function. The new function uses the authorization list to verify the signer's authorization
* [2312](https://github.com/zeta-chain/node/pull/2312) - add queries `ShowAuthorization` and `ListAuthorizations`
* [2319](https://github.com/zeta-chain/node/pull/2319) - use `CheckAuthorization` function in all messages
* [2325](https://github.com/zeta-chain/node/pull/2325) - revert telemetry server changes
* [2339](https://github.com/zeta-chain/node/pull/2339) - add binaries related question to syncing issue form
* [2366](https://github.com/zeta-chain/node/pull/2366) - add migration script for adding authorizations table
* [2372](https://github.com/zeta-chain/node/pull/2372) - add queries for tss fund migration info
* [2416g](https://github.com/zeta-chain/node/pull/2416) - add Solana chain information

### Refactor

* [2094](https://github.com/zeta-chain/node/pull/2094) - upgrade go-tss to use cosmos v0.47
* [2110](https://github.com/zeta-chain/node/pull/2110) - move non-query rate limiter logic to zetaclient side and code refactor
* [2032](https://github.com/zeta-chain/node/pull/2032) - improve some general structure of the ZetaClient codebase
* [2097](https://github.com/zeta-chain/node/pull/2097) - refactor lightclient verification flags to account for individual chains
* [2071](https://github.com/zeta-chain/node/pull/2071) - Modify chains struct to add all chain related information
* [2118](https://github.com/zeta-chain/node/pull/2118) - consolidate inbound and outbound naming
* [2124](https://github.com/zeta-chain/node/pull/2124) - removed unused variables and method
* [2150](https://github.com/zeta-chain/node/pull/2150) - created `chains` `zetacore` `orchestrator` packages in zetaclient and reorganized source files accordingly
* [2210](https://github.com/zeta-chain/node/pull/2210) - removed uncessary panics in the zetaclientd process
* [2205](https://github.com/zeta-chain/node/pull/2205) - remove deprecated variables pre-v17
* [2226](https://github.com/zeta-chain/node/pull/2226) - improve Go formatting with imports standardization and max line length to 120
* [2262](https://github.com/zeta-chain/node/pull/2262) - refactor MsgUpdateZRC20 into MsgPauseZrc20 and MsgUnPauseZRC20
* [2290](https://github.com/zeta-chain/node/pull/2290) - rename `MsgAddBlameVote` message to `MsgVoteBlame`
* [2269](https://github.com/zeta-chain/node/pull/2269) - refactor MsgUpdateCrosschainFlags into MsgEnableCCTX, MsgDisableCCTX and MsgUpdateGasPriceIncreaseFlags
* [2306](https://github.com/zeta-chain/node/pull/2306) - refactor zetaclient outbound transaction signing logic
* [2296](https://github.com/zeta-chain/node/pull/2296) - move `testdata` package to `testutil` to organize test-related utilities
* [2317](https://github.com/zeta-chain/node/pull/2317) - add ValidateOutbound method for cctx orchestrator
* [2340](https://github.com/zeta-chain/node/pull/2340) - add ValidateInbound method for cctx orchestrator
* [2344](https://github.com/zeta-chain/node/pull/2344) - group common data of EVM/Bitcoin signer and observer using base structs
* [2357](https://github.com/zeta-chain/node/pull/2357) - integrate base Signer structure into EVM/Bitcoin Signer
* [2359](https://github.com/zeta-chain/node/pull/2359) - integrate base Observer structure into EVM/Bitcoin Observer
* [2375](https://github.com/zeta-chain/node/pull/2375) - improve & speedup code formatting
* [2380](https://github.com/zeta-chain/node/pull/2380) - use `ChainInfo` in `authority` to allow dynamically support new chains
* [2395](https://github.com/zeta-chain/node/pull/2395) - converge AppContext with ZetaCoreContext in zetaclient
* [2428](https://github.com/zeta-chain/node/pull/2428) - propagate context across codebase & refactor zetacore client

### Tests

* [2047](https://github.com/zeta-chain/node/pull/2047) - fix liquidity cap advanced test
* [2181](https://github.com/zeta-chain/node/pull/2181) - add more assertion and test cases in ZEVM message passing E2E tests
* [2184](https://github.com/zeta-chain/node/pull/2184) - add tx priority checks to e2e tests
* [2199](https://github.com/zeta-chain/node/pull/2199) - custom priority mempool unit tests
* [2240](https://github.com/zeta-chain/node/pull/2240) - removed hard-coded Bitcoin regnet chainID in E2E withdraw tests
* [2266](https://github.com/zeta-chain/node/pull/2266) - try fixing E2E test `crosschain_swap` failure `btc transaction not signed`
* [2294](https://github.com/zeta-chain/node/pull/2294) - add and fix existing ethermint rpc unit test
* [2329](https://github.com/zeta-chain/node/pull/2329) - fix TODOs in rpc unit tests
* [2342](https://github.com/zeta-chain/node/pull/2342) - extend rpc unit tests with testing extension to include synthetic ethereum txs
* [2299](https://github.com/zeta-chain/node/pull/2299) - add `zetae2e` command to deploy test contracts
* [2364](https://github.com/zeta-chain/node/pull/2364) - add stateful upgrade test
* [2360](https://github.com/zeta-chain/node/pull/2360) - add stateful e2e tests.
* [2349](https://github.com/zeta-chain/node/pull/2349) - add TestBitcoinDepositRefund and WithdrawBitcoinMultipleTimes E2E tests
* [2368](https://github.com/zeta-chain/node/pull/2368) - eliminate panic usage across testing suite
* [2369](https://github.com/zeta-chain/node/pull/2369) - fix random cross-chain swap failure caused by using tiny UTXO
* [2549](https://github.com/zeta-chain/node/pull/2459) - add separate accounts for each policy in e2e tests
* [2415](https://github.com/zeta-chain/node/pull/2415) - add e2e test for upgrade and test admin functionalities

### Fixes

* [1484](https://github.com/zeta-chain/node/issues/1484) - replaced hard-coded `MaxLookaheadNonce` with a default lookback factor
* [2125](https://github.com/zeta-chain/node/pull/2125) - fix develop upgrade test
* [2222](https://github.com/zeta-chain/node/pull/2222) - removed `maxHeightDiff` to let observer scan from Bitcoin height where it left off
* [2233](https://github.com/zeta-chain/node/pull/2233) - fix `IsSupported` flag not properly updated in zetaclient's context
* [2243](https://github.com/zeta-chain/node/pull/2243) - fix incorrect bitcoin outbound height in the CCTX outbound parameter
* [2256](https://github.com/zeta-chain/node/pull/2256) - fix rate limiter falsely included reverted non-withdraw cctxs
* [2327](https://github.com/zeta-chain/node/pull/2327) - partially cherry picked the fix to Bitcoin outbound dust amount
* [2362](https://github.com/zeta-chain/node/pull/2362) - set 1000 satoshis as minimum BTC amount that can be withdrawn from zEVM
* [2382](https://github.com/zeta-chain/node/pull/2382) - add tx input and gas in rpc methods for synthetic eth txs
* [2396](https://github.com/zeta-chain/node/issues/2386) - special handle bitcoin testnet gas price estimator
* [2434](https://github.com/zeta-chain/node/pull/2434) - the default database when running `zetacored init` is now pebbledb

### CI
* [2388](https://github.com/zeta-chain/node/pull/2388) - added GitHub attestations of binaries produced in the release workflow. 
* [2285](https://github.com/zeta-chain/node/pull/2285) - added nightly EVM performance testing pipeline, modified localnet testing docker image to utilitze debian:bookworm, removed build-jet runners where applicable, removed deprecated/removed upgrade path testing pipeline
* [2268](https://github.com/zeta-chain/node/pull/2268) - updated the publish-release pipeline to utilize the Github Actions Ubuntu 20.04 Runners
* [2070](https://github.com/zeta-chain/node/pull/2070) - Added commands to build binaries from the working branch as a live full node rpc to test non-governance changes
* [2119](https://github.com/zeta-chain/node/pull/2119) - Updated the release pipeline to only run on hotfix/ and release/ branches. Added option to only run pre-checks and not cut release as well. Switched approval steps to use environments
* [2189](https://github.com/zeta-chain/node/pull/2189) - Updated the docker tag when a release trigger runs to be the github event for the release name which should be the version. Removed mac specific build as the arm build should handle that
* [2191](https://github.com/zeta-chain/node/pull/2191) - Fixed conditional logic for the docker build step for non release builds to not overwrite the github tag
* [2192](https://github.com/zeta-chain/node/pull/2192) - Added release status checker and updater pipeline that will update release statuses when they go live on network
* [2335](https://github.com/zeta-chain/node/pull/2335) - ci: updated the artillery report to publish to artillery cloud
* [2377](https://github.com/zeta-chain/node/pull/2377) - ci: adjusted sast-linters.yml to not scan itself, nor alert on removal of nosec.
* [2400](https://github.com/zeta-chain/node/pull/2400) - ci: adjusted the performance test to pass or fail pipeline based on test results, alert slack, and launch network with state. Fixed connection issues as well.
* [2425](https://github.com/zeta-chain/node/pull/2425) - Added verification to performance testing pipeline to ensure p99 aren't above 2000ms and p50 aren't above 40ms, Tweaked the config to 400 users requests per second. 425 is the current max before it starts failing.

### Documentation

* [2321](https://github.com/zeta-chain/node/pull/2321) - improve documentation for ZetaClient functions and packages

## v17.0.0

### Fixes

* [2249](https://github.com/zeta-chain/node/pull/2249) - fix inbound and outbound validation for BSC chain
* [2265](https://github.com/zeta-chain/node/pull/2265) - fix rate limiter query for revert cctxs

## v16.0.0

### Breaking Changes

* Admin policies have been moved from `observer` to a new module `authority`
  * Updating admin policies now requires to send a governance proposal executing the `UpdatePolicies` message in the `authority` module
  * The `Policies` query of the `authority` module must be used to get the current admin policies
  * `PolicyType_group1` has been renamed into `PolicyType_groupEmergency` and `PolicyType_group2` has been renamed into `PolicyType_groupAdmin`

* A new module called `lightclient` has been created for the blocker header and proof functionality to add inbound and outbound trackers in a permissionless manner (currently deactivated on live networks)
  * The list of block headers are now stored in the `lightclient` module instead of the `observer` module
    * The message to vote on new block headers is still in the `observer` module but has been renamed to `MsgVoteBlockHeader` instead of `MsgAddBlockHeader`
    * The `GetAllBlockHeaders` query has been moved to the `lightclient` module and renamed to `BlockHeaderAll`
    * The `GetBlockHeaderByHash` query has been moved to the `lightclient` module and renamed to `BlockHeader`
    * The `GetBlockHeaderStateByChain` query has been moved to the `lightclient` module and renamed to `ChainState`
    * The `Prove` query has been moved to the `lightclient` module
    * The `BlockHeaderVerificationFlags` has been deprecated in `CrosschainFlags`, `VerificationFlags` should be used instead

* `MsgGasPriceVoter` message in the `crosschain` module has been renamed to `MsgVoteGasPrice`
  * The structure of the message remains the same

* `MsgCreateTSSVoter` message in the `crosschain` module has been moved to the `observer` module and renamed to `MsgVoteTSS`
  * The structure of the message remains the same

### Refactor

* [1511](https://github.com/zeta-chain/node/pull/1511) - move ballot voting logic from `crosschain` to `observer`
* [1783](https://github.com/zeta-chain/node/pull/1783) - refactor zetaclient metrics naming and structure
* [1774](https://github.com/zeta-chain/node/pull/1774) - split params and config in zetaclient
* [1831](https://github.com/zeta-chain/node/pull/1831) - removing unnecessary pointers in context structure
* [1864](https://github.com/zeta-chain/node/pull/1864) - prevent panic in param management
* [1848](https://github.com/zeta-chain/node/issues/1848) - create a method to observe deposits to tss address in one evm block
* [1885](https://github.com/zeta-chain/node/pull/1885) - change important metrics on port 8123 to be prometheus compatible
* [1863](https://github.com/zeta-chain/node/pull/1863) - remove duplicate ValidateChainParams function
* [1914](https://github.com/zeta-chain/node/pull/1914) - move crosschain flags to core context in zetaclient
* [1948](https://github.com/zeta-chain/node/pull/1948) - remove deprecated GetTSSAddress query in crosschain module
* [1936](https://github.com/zeta-chain/node/pull/1936) - refactor common package into subpackages and rename to pkg
* [1966](https://github.com/zeta-chain/node/pull/1966) - move TSS vote message from crosschain to observer
* [1853](https://github.com/zeta-chain/node/pull/1853) - refactor vote inbound tx and vote outbound tx
* [1815](https://github.com/zeta-chain/node/pull/1815) - add authority module for authorized actions
* [1976](https://github.com/zeta-chain/node/pull/1976) - add lightclient module for header and proof functionality
* [2001](https://github.com/zeta-chain/node/pull/2001) - replace broadcast mode block with sync and remove fungible params
* [1989](https://github.com/zeta-chain/node/pull/1989) - simplify `IsSendOutTxProcessed` method and add unit tests
* [2013](https://github.com/zeta-chain/node/pull/2013) - rename `GasPriceVoter` message to `VoteGasPrice`
* [2059](https://github.com/zeta-chain/node/pull/2059) - Remove unused params from all functions in zetanode
* [2071](https://github.com/zeta-chain/node/pull/2071) - Modify chains struct to add all chain related information
* [2076](https://github.com/zeta-chain/node/pull/2076) - automatically deposit native zeta to an address if it doesn't exist on ZEVM
* [2169](https://github.com/zeta-chain/node/pull/2169) - Limit zEVM revert transactions to coin type ZETA

### Features

* [1789](https://github.com/zeta-chain/node/issues/1789) - block cross-chain transactions that involve restricted addresses
* [1755](https://github.com/zeta-chain/node/issues/1755) - use evm JSON RPC for inbound tx (including blob tx) observation
* [1884](https://github.com/zeta-chain/node/pull/1884) - added zetatool cmd, added subcommand to filter deposits
* [1942](https://github.com/zeta-chain/node/pull/1982) - support Bitcoin P2TR, P2WSH, P2SH, P2PKH addresses
* [1935](https://github.com/zeta-chain/node/pull/1935) - add an operational authority group
* [1954](https://github.com/zeta-chain/node/pull/1954) - add metric for concurrent keysigns
* [1979](https://github.com/zeta-chain/node/pull/1979) - add script to import genesis data into an existing genesis file
* [2006](https://github.com/zeta-chain/node/pull/2006) - add Amoy testnet static chain information
* [2045](https://github.com/zeta-chain/node/pull/2046) - add grpc query with outbound rate limit for zetaclient to use
* [2046](https://github.com/zeta-chain/node/pull/2046) - add state variable in crosschain for rate limiter flags
* [2034](https://github.com/zeta-chain/node/pull/2034) - add support for zEVM message passing
* [1825](https://github.com/zeta-chain/node/pull/1825) - add a message to withdraw emission rewards

### Tests

* [1767](https://github.com/zeta-chain/node/pull/1767) - add unit tests for emissions module begin blocker
* [1816](https://github.com/zeta-chain/node/pull/1816) - add args to e2e tests
* [1791](https://github.com/zeta-chain/node/pull/1791) - add e2e tests for feature of restricted address
* [1787](https://github.com/zeta-chain/node/pull/1787) - add unit tests for cross-chain evm hooks and e2e test failed withdraw to BTC legacy address
* [1840](https://github.com/zeta-chain/node/pull/1840) - fix code coverage test failures ignored in CI
* [1870](https://github.com/zeta-chain/node/pull/1870) - enable emissions pool in local e2e testing
* [1868](https://github.com/zeta-chain/node/pull/1868) - run e2e btc tests locally
* [1851](https://github.com/zeta-chain/node/pull/1851) - rename usdt to erc20 in e2e tests
* [1872](https://github.com/zeta-chain/node/pull/1872) - remove usage of RPC in unit test
* [1805](https://github.com/zeta-chain/node/pull/1805) - add admin and performance test and fix upgrade test
* [1879](https://github.com/zeta-chain/node/pull/1879) - full coverage for messages in types packages
* [1899](https://github.com/zeta-chain/node/pull/1899) - add empty test files so packages are included in coverage
* [1900](https://github.com/zeta-chain/node/pull/1900) - add testing for external chain migration
* [1903](https://github.com/zeta-chain/node/pull/1903) - common package tests
* [1961](https://github.com/zeta-chain/node/pull/1961) - improve observer module coverage
* [1967](https://github.com/zeta-chain/node/pull/1967) - improve crosschain module coverage
* [1955](https://github.com/zeta-chain/node/pull/1955) - improve emissions module coverage
* [1941](https://github.com/zeta-chain/node/pull/1941) - add unit tests for zetacore package
* [1985](https://github.com/zeta-chain/node/pull/1985) - improve fungible module coverage
* [1992](https://github.com/zeta-chain/node/pull/1992) - remove setupKeeper from crosschain module
* [2008](https://github.com/zeta-chain/node/pull/2008) - add test for connector bytecode update
* [2047](https://github.com/zeta-chain/node/pull/2047) - fix liquidity cap advanced test
* [2076](https://github.com/zeta-chain/node/pull/2076) - automatically deposit native zeta to an address if it doesn't exist on ZEVM

### Fixes

* [1861](https://github.com/zeta-chain/node/pull/1861) - fix `ObserverSlashAmount` invalid read
* [1880](https://github.com/zeta-chain/node/issues/1880) - lower the gas price multiplier for EVM chains
* [1883](https://github.com/zeta-chain/node/issues/1883) - zetaclient should check 'IsSupported' flag to pause/unpause a specific chain
* * [2076](https://github.com/zeta-chain/node/pull/2076) - automatically deposit native zeta to an address if it doesn't exist on ZEVM
* [1633](https://github.com/zeta-chain/node/issues/1633) - zetaclient should be able to pick up new connector and erc20Custody addresses
* [1944](https://github.com/zeta-chain/node/pull/1944) - fix evm signer unit tests
* [1888](https://github.com/zeta-chain/node/issues/1888) - zetaclient should stop inbound/outbound txs according to cross-chain flags
* [1970](https://github.com/zeta-chain/node/issues/1970) - remove the timeout in the evm outtx tracker processing thread

### Chores

* [1814](https://github.com/zeta-chain/node/pull/1814) - fix code coverage ignore for protobuf generated files

### CI

* [1958](https://github.com/zeta-chain/node/pull/1958) - Fix e2e advanced test debug checkbox
* [1945](https://github.com/zeta-chain/node/pull/1945) - update advanced testing pipeline to not execute tests that weren't selected so they show skipped instead of skipping steps
* [1940](https://github.com/zeta-chain/node/pull/1940) - adjust release pipeline to be created as pre-release instead of latest
* [1867](https://github.com/zeta-chain/node/pull/1867) - default restore_type for full node docker-compose to snapshot instead of statesync for reliability
* [1891](https://github.com/zeta-chain/node/pull/1891) - fix typo that was introduced to docker-compose and a typo in start.sh for the docker start script for full nodes
* [1894](https://github.com/zeta-chain/node/pull/1894) - added download binaries and configs to the start sequence so it will download binaries that don't exist
* [1953](https://github.com/zeta-chain/node/pull/1953) - run E2E tests for all PRs

## Version: v15.0.0

### Features
* [1912](https://github.com/zeta-chain/node/pull/1912) - add reset chain nonces msg

## Version: v14.0.1

- [1817](https://github.com/zeta-chain/node/pull/1817) - Add migration script to fix pending and chain nonces on testnet

## Version: v13.0.0

### Breaking Changes

* `zetaclientd start`: now requires 2 inputs from stdin: hotkey password and tss keyshare password
  Starting zetaclient now requires two passwords to be input; one for the hotkey and another for the tss key-share

### Features

* [1698](https://github.com/zeta-chain/node/issues/1698) - bitcoin dynamic depositor fee

### Docs

* [1731](https://github.com/zeta-chain/node/pull/1731) added doc for hotkey and tss key-share password prompts

### Features

* [1728] (https://github.com/zeta-chain/node/pull/1728) - allow aborted transactions to be refunded by minting tokens to zEvm

### Refactor

* [1766](https://github.com/zeta-chain/node/pull/1766) - Refactors the `PostTxProcessing` EVM hook functionality to deal with invalid withdraw events
* [1630](https://github.com/zeta-chain/node/pull/1630) - added password prompts for hotkey and tss keyshare in zetaclient
* [1760](https://github.com/zeta-chain/node/pull/1760) - Make staking keeper private in crosschain module
* [1809](https://github.com/zeta-chain/node/pull/1809) - Refactored tryprocessout function in evm signer

### Fixes

* [1678](https://github.com/zeta-chain/node/issues/1678) - clean cached stale block to fix evm outtx hash mismatch
* [1690](https://github.com/zeta-chain/node/issues/1690) - double watched gas prices and fix btc scheduler
* [1687](https://github.com/zeta-chain/node/pull/1687) - only use EVM supported chains for gas stability pool
* [1692](https://github.com/zeta-chain/node/pull/1692) - fix get params query for emissions module
* [1706](https://github.com/zeta-chain/node/pull/1706) - fix CLI crosschain show-out-tx-tracker
* [1707](https://github.com/zeta-chain/node/issues/1707) - fix bitcoin fee rate estimation
* [1712](https://github.com/zeta-chain/node/issues/1712) - increase EVM outtx inclusion timeout to 20 minutes
* [1733](https://github.com/zeta-chain/node/pull/1733) - remove the unnecessary 2x multiplier in the convertGasToZeta RPC
* [1721](https://github.com/zeta-chain/node/issues/1721) - zetaclient should provide bitcoin_chain_id when querying TSS address
* [1744](https://github.com/zeta-chain/node/pull/1744) - added cmd to encrypt tss keyshare file, allowing empty tss password for backward compatibility

### Tests

* [1584](https://github.com/zeta-chain/node/pull/1584) - allow to run E2E tests on any networks
* [1746](https://github.com/zeta-chain/node/pull/1746) - rename smoke tests to e2e tests
* [1753](https://github.com/zeta-chain/node/pull/1753) - fix gosec errors on usage of rand package
* [1762](https://github.com/zeta-chain/node/pull/1762) - improve coverage for fungibile module
* [1782](https://github.com/zeta-chain/node/pull/1782) - improve coverage for fungibile module system contract

### CI

* Adjusted the release pipeline to be a manually executed pipeline with an approver step. The pipeline now executes all the required tests run before the approval step unless skipped
* Added pipeline to build and push docker images into dockerhub on release for ubuntu and macos
* Adjusted the pipeline for building and pushing docker images for MacOS to install and run docker
* Added docker-compose and make commands for launching full nodes. `make mainnet-zetarpc-node`  `make mainnet-bitcoind-node`
* Made adjustments to the docker-compose for launching mainnet full nodes to include examples of using the docker images build from the docker image build pipeline
* [1736](https://github.com/zeta-chain/node/pull/1736) - chore: add Ethermint endpoints to OpenAPI
* Re-wrote Dockerfile for building Zetacored docker images
* Adjusted the docker-compose files for Zetacored nodes to utilize the new docker image
* Added scripts for the new docker image that facilitate the start up automation
* Adjusted the docker pipeline slightly to pull the version on PR from the app.go file
* [1781](https://github.com/zeta-chain/node/pull/1781) - add codecov coverage report in CI
* fixed the download binary script to use relative pathing from binary_list file

### Features

* [1425](https://github.com/zeta-chain/node/pull/1425) add `whitelist-erc20` command

### Chores

* [1729](https://github.com/zeta-chain/node/pull/1729) - add issue templates
* [1754](https://github.com/zeta-chain/node/pull/1754) - cleanup expected keepers

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
* [1675](https://github.com/zeta-chain/node/issues/1675) - use chain param ConfirmationCount for bitcoin confirmation

## Chores

* [1694](https://github.com/zeta-chain/node/pull/1694) - remove standalone network, use require testing package for the entire node folder

## Version: v12.1.0

### Tests

* [1577](https://github.com/zeta-chain/node/pull/1577) - add chain header tests in E2E tests and fix admin tests

### Features

* [1658](https://github.com/zeta-chain/node/pull/1658) - modify emission distribution to use fixed block rewards

### Fixes

* [1535](https://github.com/zeta-chain/node/issues/1535) - Avoid voting on wrong ballots due to false blockNumber in EVM tx receipt
* [1588](https://github.com/zeta-chain/node/pull/1588) - fix chain params comparison logic
* [1650](https://github.com/zeta-chain/node/pull/1605) - exempt (discounted) *system txs* from min gas price check and gas fee deduction
* [1632](https://github.com/zeta-chain/node/pull/1632) - set keygen to `KeygenStatus_KeyGenSuccess` if its in `KeygenStatus_PendingKeygen`
* [1576](https://github.com/zeta-chain/node/pull/1576) - Fix zetaclient crash due to out of bound integer conversion and log prints
* [1575](https://github.com/zeta-chain/node/issues/1575) - Skip unsupported chain parameters by IsSupported flag

### CI

* [1580](https://github.com/zeta-chain/node/pull/1580) - Fix release pipelines cleanup step

### Chores

* [1585](https://github.com/zeta-chain/node/pull/1585) - Updated release instructions
* [1615](https://github.com/zeta-chain/node/pull/1615) - Add upgrade handler for version v12.1.0

### Features

* [1591](https://github.com/zeta-chain/node/pull/1591) - support lower gas limit for voting on inbound and outbound transactions
* [1592](https://github.com/zeta-chain/node/issues/1592) - check inbound tracker tx hash against Tss address and some refactor on inTx observation

### Refactoring

* [1628](https://github.com/zeta-chain/node/pull/1628) optimize return and simplify code
* [1640](https://github.com/zeta-chain/node/pull/1640) reorganize zetaclient into subpackages
* [1619](https://github.com/zeta-chain/node/pull/1619) - Add evm fee calculation to tss migration of evm chains

## Version: v12.0.0

### Breaking Changes

TSS and chain validation related queries have been moved from `crosschain` module to `observer` module:
* `PendingNonces` :Changed from `/zeta-chain/crosschain/pendingNonces/{chain_id}/{address}` to `/zeta-chain/observer/pendingNonces/{chain_id}/{address}` . It returns all the pending nonces for a chain id and address. This returns the current pending nonces for the chain
* `ChainNonces` : Changed from `/zeta-chain/crosschain/chainNonces/{chain_id}` to`/zeta-chain/observer/chainNonces/{chain_id}` . It returns all the chain nonces for a chain id. This returns the current nonce of the TSS address for the chain
* `ChainNoncesAll` :Changed from `/zeta-chain/crosschain/chainNonces` to `/zeta-chain/observer/chainNonces` . It returns all the chain nonces for all chains. This returns the current nonce of the TSS address for all chains

All chains now have the same observer set:
* `ObserversByChain`: `/zeta-chain/observer/observers_by_chain/{observation_chain}` has been removed and replaced with `/zeta-chain/observer/observer_set`. All chains have the same observer set
* `AllObserverMappers`: `/zeta-chain/observer/all_observer_mappers` has been removed. `/zeta-chain/observer/observer_set` should be used to get observers.

Observer params and core params have been merged into chain params:
* `Params`: `/zeta-chain/observer/params` no longer returns observer params. Observer params data have been moved to chain params described below.
* `GetCoreParams`: Renamed into `GetChainParams`. `/zeta-chain/observer/get_core_params` moved to `/zeta-chain/observer/get_chain_params`
* `GetCoreParamsByChain`: Renamed into `GetChainParamsForChain`. `/zeta-chain/observer/get_core_params_by_chain` moved to `/zeta-chain/observer/get_chain_params_by_chain`

Getting the correct TSS address for Bitcoin now requires proviidng the Bitcoin chain id:
* `GetTssAddress` : Changed from `/zeta-chain/observer/get_tss_address/` to `/zeta-chain/observer/getTssAddress/{bitcoin_chain_id}` . Optional bitcoin chain id can now be passed as a parameter to fetch the correct tss for required BTC chain. This parameter only affects the BTC tss address in the response

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

* [1446](https://github.com/zeta-chain/node/pull/1446) - renamed file `zetaclientd/aux.go` to `zetaclientd/utils.go` to avoid complaints from go package resolver
* [1499](https://github.com/zeta-chain/node/pull/1499) - Add scripts to localnet to help test gov proposals
* [1442](https://github.com/zeta-chain/node/pull/1442) - remove build types in `.goreleaser.yaml`
* [1504](https://github.com/zeta-chain/node/pull/1504) - remove `-race` in the `make install` commmand
* [1564](https://github.com/zeta-chain/node/pull/1564) - bump ti-actions/changed-files

### Tests

* [1538](https://github.com/zeta-chain/node/pull/1538) - improve stateful e2e testing

### CI

* Removed private runners and unused GitHub Action

## Version: v11.0.0

### Features

* [1387](https://github.com/zeta-chain/node/pull/1387) - Add HSM capability for zetaclient hot key
* add a new thread to zetaclient which checks zeta supply in all connected chains in every block
* add a new tx to update an observer, this can be either be run a tombstoned observer/validator or via admin_policy_group_2

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

* Add unit tests for adding votes to a ballot 

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
