# CHANGELOG

## Unreleased

### Breaking Changes

* EVM inbounds support multiple calls inside same tx. EVM Gateway contracts must be upgraded before node upgrade, and an additional action fee should be set,
by calling `updateAdditionalActionFee` admin function.
* `confirmation_count` in the chain params has been entirely removed. It was made deprecated in v28. Use `confirmation_params` instead.
* The message `MsgWhitelistERC20` has been renamed to `MsgWhitelistAsset`. The message parameters remain unchanged.
  * The event `EventERC20Whitelist` has been renamed to `EventAssetWhitelist`. The event parameters remain unchanged.

### Features

* [4064](https://github.com/zeta-chain/node/pull/4064) - add support for withdraws using the new v2 connector contract
* [4138](https://github.com/zeta-chain/node/pull/4138) - allow zetaclient to observe inbound events from Sui original gateway package
* [4153](https://github.com/zeta-chain/node/pull/4153) - make the gas limit used for gateway calls a configurable parameter
* [4157](https://github.com/zeta-chain/node/pull/4157) - multiple evm calls in single tx
* [4211](https://github.com/zeta-chain/node/pull/4211) - provide error information in cctx when Bitcoin deposit fail
* [4218](https://github.com/zeta-chain/node/pull/4218) - enable NoAssetCall from Bitcoin chain
* [3834](https://github.com/zeta-chain/node/pull/3734) - refund a portion of remaining unused tokens to user

### Refactor

* [4070](https://github.com/zeta-chain/node/pull/4070) - remove support for v1 revert address for BTC
* [4144](https://github.com/zeta-chain/node/pull/4144) - standardize structured logging for zetaclient
* [4192](https://github.com/zeta-chain/node/pull/4192) - remove deprecated code in observer module, including `confirmation_count`
* [4180](https://github.com/zeta-chain/node/pull/4180) - remove unused loggers and log fields
* [4203](https://github.com/zeta-chain/node/pull/4203) - rename `whitelistERC20` into `whitelistAsset`
* [4199](https://github.com/zeta-chain/node/pull/4199) - remove `MsgUpdateERC20CustodyPauseStatus` and `MsgMigrateERC20CustodyFunds`
* [4205](https://github.com/zeta-chain/node/pull/4205) - remove index field in ballot
* [4200](https://github.com/zeta-chain/node/pull/4200) - remove `LastBlockHeight` state variable
* [4174](https://github.com/zeta-chain/node/pull/4174) - add documentation for ZetaClient logging fields
* [4210](https://github.com/zeta-chain/node/pull/4210) - skip writing config file to the filesystem when updating consensus timeout deltas
* [4213](https://github.com/zeta-chain/node/pull/4213) - prepare the client interfaces of the observer-signers for dry mode

### Fixes

* [4090](https://github.com/zeta-chain/node/pull/4090) - print error message in detail if unable to decode Bitcoin memo
* [4116](https://github.com/zeta-chain/node/pull/4116) - remove confirmation mode from outbound and inbound digest
* [4111](https://github.com/zeta-chain/node/pull/4111) - cancel Solana outbound if transaction size is too large
* [4112](https://github.com/zeta-chain/node/pull/4112) - fix error when deploying contracts on testnet
* [4121](https://github.com/zeta-chain/node/pull/4121) - dbg trace block by number gas limit legacy
* [4169](https://github.com/zeta-chain/node/pull/4169) - unpack revert message from Bitcoin memo without considering `CallOnRevert` flag
* [4194](https://github.com/zeta-chain/node/pull/4194) - remove duplicate solana post-gas-price goroutine
* [4197](https://github.com/zeta-chain/node/pull/4197) - re-check for finalized ballot when executing inbound vote to create cctx
* [4217](https://github.com/zeta-chain/node/pull/4217) - remove ZetaChain chain ID from GasStabilityPoolBalances query

### Tests

* [4071](https://github.com/zeta-chain/node/pull/4071) - use v2 connector contract in admin e2e tests
* [4113](https://github.com/zeta-chain/node/pull/4113) - fix Solana flaky depositAndCall e2e tests in live networks
* [4142](https://github.com/zeta-chain/node/pull/4142) - fix Solana flaky SPL deposit e2e test in live networks
* [4158](https://github.com/zeta-chain/node/pull/4158) - have e2e tests interact with pre-deployed example dApp contract
* [4165](https://github.com/zeta-chain/node/pull/4165) - fix Sui flaky depositAndCall e2e test in live networks
* [4177](https://github.com/zeta-chain/node/pull/4177) - add an E2E test to verify depositAndCall with high gas consumption

## v36.0.2

### Fixes

* [4202](https://github.com/zeta-chain/node/pull/4202) - force rescan if inbound vote monitoring fails using a context that can timeout

## v36.0.0

### Features

* [4153](https://github.com/zeta-chain/node/pull/4153) - make the gas limit used for gateway calls a configurable parameter

## v33.0.0

### Features

* [3961](https://github.com/zeta-chain/node/pull/3961) - cosmos v53 upgrade
* [3977](https://github.com/zeta-chain/node/pull/3977) - integrate new ton features (call, increase_seqno, compliance)
* [3990](https://github.com/zeta-chain/node/pull/3990) - add support for deposits using the new v2 connector contract
* [4007](https://github.com/zeta-chain/node/pull/4007) - add support for Sui withdraw and authenticated call
* [4063](https://github.com/zeta-chain/node/pull/4063) - add a message to burn funds in the fungible module
* [3991](https://github.com/zeta-chain/node/pull/3991) - migrate from Ethermint to new Cosmos EVM module
* [4077](https://github.com/zeta-chain/node/pull/4077) - add Cosmos EVM default precompiles support
* [4068](https://github.com/zeta-chain/node/pull/4068) - add support for hostname in addition to public IP for zetaclient to connect to zetacore
* [4087](https://github.com/zeta-chain/node/pull/4087) - re-enable Sui authenticated call and adopt MessageContext ID as a gateway dynamic field

### Fixes

* [3953](https://github.com/zeta-chain/node/pull/3953) - skip Bitcoin outbound when scanning inbound transactions
* [3957](https://github.com/zeta-chain/node/pull/3957) - remove tx format assumption from solana parse gateway instruction
* [3956](https://github.com/zeta-chain/node/pull/3956) - use the latest nonce to perform pre-broadcast check to make evm tx replacement possible
* [3954](https://github.com/zeta-chain/node/pull/3954) - fail sui withdrawal event in the ZEVM if it carries invalid receiver address
* [3917](https://github.com/zeta-chain/node/pull/3917) - prevent jailed observers from voting
* [3971](https://github.com/zeta-chain/node/pull/3971) - zetaclient should load restricted addresses correctly from `zetaclient_restricted_addresses.json`
* [3964](https://github.com/zeta-chain/node/pull/3964) - use the inscription initiator address as Bitcoin inbound sender address
* [4018](https://github.com/zeta-chain/node/pull/4018) - Sui token accounting mismatch
* [4020](https://github.com/zeta-chain/node/pull/4020) - add a migration script to delete ZRC20 SUI gas tokens from stability pool
* [4067](https://github.com/zeta-chain/node/pull/4067) - disable sui authenticated call temporarily until gateway upgrade
* [4072](https://github.com/zeta-chain/node/pull/4072) - improve rate limiter robustness and code quality

### Refactor

* [3940](https://github.com/zeta-chain/node/pull/3940) - avoid pre-signing solana outbound by waiting for the exact PDA nonce to arrive
* [4033](https://github.com/zeta-chain/node/pull/4033) - improve error handling logic when decoding Bitcoin inscribed memo standard
* [4052](https://github.com/zeta-chain/node/pull/4052) - remove the deprecated 'BitcoinConfig' from zetaclient config
* [4054](https://github.com/zeta-chain/node/pull/4054) - refactor deposit logic to prevent minting surplus zeta
* [4060](https://github.com/zeta-chain/node/pull/4060) - cleanup forked mempool code
* [4097](https://github.com/zeta-chain/node/pull/4097) - refactor zetacore to overwrite consensus config values for all timeout deltas on startup.

### Tests

* [3972](https://github.com/zeta-chain/node/pull/3972) - add `gasLimit` argument to erc20 withdrawAndCall e2e test
* [3983](https://github.com/zeta-chain/node/pull/3983) - fix simulation tests and modify `sim.yml` workflow to run sim tests on changes to `app/` and `cmd/` directories in addition to `x/`
* [3975](https://github.com/zeta-chain/node/pull/3975) - add v2 connector deployment to the e2e test.
* [3976](https://github.com/zeta-chain/node/pull/3976) - add connector fund migration e2e test using contracts only.
* [3999](https://github.com/zeta-chain/node/pull/3999) - run simulation tests nightly
* [3985](https://github.com/zeta-chain/node/pull/3985) - add e2e tests for deposit and withdraw with big payload
* [4022](https://github.com/zeta-chain/node/pull/4022) - configure Solana e2e test connected program IDs
* [4040](https://github.com/zeta-chain/node/pull/4040) - make Bitcoin inscription e2e test working for live networks
* [4053](https://github.com/zeta-chain/node/pull/4053) - fix flaky e2e tests that failed on zrc20 balance assertion in live networks

## v32.0.0

### Chores

* [4003](https://github.com/zeta-chain/node/pull/4003) - upgrade Cosmos SDK to v0.50.14

## v31.0.1

### Features

* [3929](https://github.com/zeta-chain/node/pull/3929) - add support for TON http-rpc
* [3958](https://github.com/zeta-chain/node/pull/3958) - integrate TON http-rpc

## v31.0.0

### Breaking Changes

* All Solana inbounds have new optional param `revert_options`. Solana Gateway program must be upgraded after node upgrade.

### Features

* [3672](https://github.com/zeta-chain/node/pull/3672) - zetaclient: cache tss signatures for performance.
* [3671](https://github.com/zeta-chain/node/pull/3671) - use gas budget argument to refund TSS for Sui withdraw cost
* [3699](https://github.com/zeta-chain/node/pull/3699) - use real gas usage for TON withdrawals
* [3710](https://github.com/zeta-chain/node/pull/3710) - support preflight RPC health metrics before fully enable a chain
* [3377](https://github.com/zeta-chain/node/pull/3377) - have zetacore increase the gas price in pending Bitcoin cctxs and burns additional gas fees
* [3750](https://github.com/zeta-chain/node/pull/3750) - support simple call from solana
* [3764](https://github.com/zeta-chain/node/pull/3764) - add payload parsing for Sui WithdrawAndCall
* [3756](https://github.com/zeta-chain/node/pull/3756) - parse revert options in solana inbounds
* [3765](https://github.com/zeta-chain/node/pull/3765) - support cancelling Sui rejected withdrawal
* [3792](https://github.com/zeta-chain/node/pull/3792) - add compliance check for Sui inbound and outbound
* [3790](https://github.com/zeta-chain/node/pull/3790) - integrate execute revert
* [3797](https://github.com/zeta-chain/node/pull/3797) - integrate execute SPL revert
* [3807](https://github.com/zeta-chain/node/pull/3807) - integrate ZEVM to Solana call
* [3826](https://github.com/zeta-chain/node/pull/3826) - add global tss signature rate-limiter to zetaclient
* [3793](https://github.com/zeta-chain/node/pull/3793) - support Sui withdrawAndCall using the PTB transaction
* [3839](https://github.com/zeta-chain/node/pull/3839) - parse Solana inbounds from inner instructions
* [3837](https://github.com/zeta-chain/node/pull/3837) - cancel Sui withdrawAndCall if tx cannot go through, e.g. on_call fails due to invalid data
* [3396](https://github.com/zeta-chain/node/pull/3396) - add support for Bitcoin RBF (Replace-By-Fee) in zetaclient
* [3864](https://github.com/zeta-chain/node/pull/3864) - add compliance checks for TON inbounds
* [3881](https://github.com/zeta-chain/node/pull/3881) - add zetatool cmd to analyze size of application.db
* [3906](https://github.com/zeta-chain/node/pull/3906) - revert restricted cctx for EVM, bitcoin and solana chains
* [3918](https://github.com/zeta-chain/node/pull/3918) - attach failure reason to solana increment nonce
* [3932](https://github.com/zeta-chain/node/pull/3932) - hardcoded block time related params

### Refactor

* [3709](https://github.com/zeta-chain/node/pull/3709) - improve cctx error message for out of gas errors when creating outbound
* [3777](https://github.com/zeta-chain/node/pull/3777) - use SignBatch keysign for solana outbound tx and fallback tx
* [3813](https://github.com/zeta-chain/node/pull/3813) - set ZETA protocol fee to 0
* [3848](https://github.com/zeta-chain/node/pull/3848) - extend min gas limit check to prevent intrinsic low gas limit

### Fixes

* [3711](https://github.com/zeta-chain/node/pull/3711) - fix TON call_data parsing
* [3717](https://github.com/zeta-chain/node/pull/3717) - fix solana withdraw and call panic
* [3770](https://github.com/zeta-chain/node/pull/3770) - improve fallback tx error handling
* [3802](https://github.com/zeta-chain/node/pull/3802) - prevent Sui withdraw with invalid address
* [3786](https://github.com/zeta-chain/node/pull/3786) - reorder end block order to allow gov changes to be added before staking.
* [3821](https://github.com/zeta-chain/node/pull/3821) - set retry gas limit if outbound is successful
* [3847](https://github.com/zeta-chain/node/pull/3847) - have EVM chain tracker reporter monitor `nonce too low` outbound hashes
* [3863](https://github.com/zeta-chain/node/pull/3863) - give enough timeout to the EVM chain transaction broadcasting
* [3850](https://github.com/zeta-chain/node/pull/3850) - broadcast single sui withdraw tx at a time to avoid nonce mismatch failure
* [3877](https://github.com/zeta-chain/node/pull/3877) - use multiple SUI coin objects to pay PTB transaction gas fee
* [3890](https://github.com/zeta-chain/node/pull/3890) - solana abort address format
* [3901](https://github.com/zeta-chain/node/pull/3901) - prevent cctx being set as abortRefunded if the abort processing failed before the refund
* [3872](https://github.com/zeta-chain/node/pull/3872) - delete testnet ballots for creation height 0 and add a query to list all ballots created at a height.
* [3916](https://github.com/zeta-chain/node/pull/3916) - infinite scan in filter solana inbound events
* [3914](https://github.com/zeta-chain/node/pull/3914) - check tx result err in filter inbound events
* [3904](https://github.com/zeta-chain/node/pull/3904) - improve observer emissions distribution to maximise pool utilisation
* [3895](https://github.com/zeta-chain/node/pull/3895) - solana call required accounts number condition
* [3896](https://github.com/zeta-chain/node/pull/3896) - add sender to solana execute message hash
* [3920](https://github.com/zeta-chain/node/pull/3920) - show correct gas limit for synthetic txs
* [3934](https://github.com/zeta-chain/node/pull/3934) - post zero priority fee for EVM chains to avoid gas price pump failure in the zetacore

### Tests

* [3692](https://github.com/zeta-chain/node/pull/3692) - e2e staking test for `MsgUndelegate` tx, to test observer staking hooks
* [3831](https://github.com/zeta-chain/node/pull/3831) - e2e tests for sui fungible token withdraw and call
* [3582](https://github.com/zeta-chain/node/pull/3852) - add solana to tss migration e2e tests
* [3866](https://github.com/zeta-chain/node/pull/3866) - add e2e test for upgrading sui gateway package
* [3417](https://github.com/zeta-chain/node/pull/3417) - add e2e test for the Bitcoin RBF (Replace-By-Fee) feature
* [3885](https://github.com/zeta-chain/node/pull/3885) - add e2e test for MsgAddObserver
* [3893](https://github.com/zeta-chain/node/pull/3893) - add e2e performance tests for sui deposit and withdrawal


### Refactor

* [3700](https://github.com/zeta-chain/node/pull/3700) - use sender and senderEVM in cross-chain call message context

## v29.0.0

### Breaking Changes

* The CCTX List RPC (`/zeta-chain/crosschain/cctx`) will now return CCTXs ordered by creation time. CCTXs from before the upgrade will not be displayed. Use the `?unordered=true` parameter to revert to the old behavior.

### Features

* [3414](https://github.com/zeta-chain/node/pull/3414) - support advanced abort workflow (onAbort)
* [3461](https://github.com/zeta-chain/node/pull/3461) - add new `ConfirmationParams` field to chain params to enable multiple confirmation count values, deprecating `confirmation_count`
* [3489](https://github.com/zeta-chain/node/pull/3489) - add Sui chain info
* [3455](https://github.com/zeta-chain/node/pull/3455) - add `track-cctx` command to zetatools
* [3506](https://github.com/zeta-chain/node/pull/3506) - define `ConfirmationMode` enum and add it to `InboundParams`, `OutboundParams`, `MsgVoteInbound` and `MsgVoteOutbound`
* [3469](https://github.com/zeta-chain/node/pull/3469) - add `MsgRemoveInboundTracker` to remove inbound trackers. This message can be triggered by the emergency policy.
* [3450](https://github.com/zeta-chain/node/pull/3450) - integrate SOL withdraw and call
* [3538](https://github.com/zeta-chain/node/pull/3538) - implement `MsgUpdateOperationalChainParams` for updating operational-related chain params with operational policy
* [3534](https://github.com/zeta-chain/node/pull/3534) - Add Sui deposit & depositAndCall
* [3541](https://github.com/zeta-chain/node/pull/3541) - implement `MsgUpdateZRC20Name` to update the name or symbol of a ZRC20 token
* [3439](https://github.com/zeta-chain/node/pull/3439) - use protocol contracts V2 with TON deposits
* [3520](https://github.com/zeta-chain/node/pull/3520) - integrate SPL withdraw and call
* [3527](https://github.com/zeta-chain/node/pull/3527) - integrate SOL/SPL withdraw and call revert
* [3522](https://github.com/zeta-chain/node/pull/3522) - add `MsgDisableFastConfirmation` to disable fast confirmation. This message can be triggered by the emergency policy.
* [3548](https://github.com/zeta-chain/node/pull/3548) - ensure cctx list is sorted by creation time
* [3562](https://github.com/zeta-chain/node/pull/3562) - add Sui withdrawals
* [3600](https://github.com/zeta-chain/node/pull/3600) - add dedicated zetaclient restricted addresses config. This file will be automatically reloaded when it changes without needing to restart zetaclient.
* [3578](https://github.com/zeta-chain/node/pull/3578) - Add disable_tss_block_scan parameter. This parameter will be used to disable expensive block scanning actions on non-ethereum EVM Chains.
* [3551](https://github.com/zeta-chain/node/pull/3551) - support for EVM chain and Bitcoin chain inbound fast confirmation
* [3615](https://github.com/zeta-chain/node/pull/3615) - make Bitcoin deposit with invalid memo reverting

### Refactor

* [3381](https://github.com/zeta-chain/node/pull/3381) - split Bitcoin observer and signer into small files and organize outbound logic into reusable/testable functions; renaming, type unification, etc.
* [3496](https://github.com/zeta-chain/node/pull/3496) - zetaclient uses `ConfirmationParams` instead of old `ConfirmationCount`; use block ranged based observation for btc and evm chain.
* [3594](https://github.com/zeta-chain/node/pull/3594) - set outbound hash in cctx when adding outbound tracker
* [3553](https://github.com/zeta-chain/node/pull/3553) — add a new buffer blocks param to delay deletion of pending ballots

### Fixes

* [3501](https://github.com/zeta-chain/node/pull/3501) - fix E2E test failure caused by nil `ConfirmationParams` for Solana and TON
* [3509](https://github.com/zeta-chain/node/pull/3509) - schedule Bitcoin TSS keysign on interval to avoid TSS keysign spam
* [3517](https://github.com/zeta-chain/node/pull/3517) - remove duplicate gateway event appending to fix false positive on multiple events in same tx
* [3602](https://github.com/zeta-chain/node/pull/3602) - hardcode gas limits to avoid estimate gas calls
* [3622](https://github.com/zeta-chain/node/pull/3622) - allow object for tracerConfig in `debug_traceTransaction` RPC
* [3634](https://github.com/zeta-chain/node/pull/3634) - return proper synthetic tx in `eth_getBlockByNumber` RPC
* [3754](https://github.com/zeta-chain/node/pull/3754) - make the Bitcoin deposit to revert when the memo output is missing

### Tests

* [3430](https://github.com/zeta-chain/node/pull/3430) - add simulation test for MsgWithDrawEmission
* [3503](https://github.com/zeta-chain/node/pull/3503) - add check in e2e test to ensure deletion of stale ballots
* [3536](https://github.com/zeta-chain/node/pull/3536) - add e2e test for upgrading solana gateway program
* [3560](https://github.com/zeta-chain/node/pull/3560) - initialize Sui E2E deposit tests
* [3595](https://github.com/zeta-chain/node/pull/3595) - add E2E tests for Sui withdraws
* [3591](https://github.com/zeta-chain/node/pull/3591) - add a runner for gov proposals in the e2e test.
* [3612](https://github.com/zeta-chain/node/pull/3612) - add support for TON live e2e tests

## v28.0.0

v28 is based on the release/v27 branch rather than develop

### Fixes
* [3563](https://github.com/zeta-chain/node/pull/3563) - upgrade cosmos-sdk to v0.50.12 to resolve GHSA-x5vx-95h7-rv4p

## v27.0.5

### Fixes
* [3554](https://github.com/zeta-chain/node/pull/3554) - disable observation of direct to TSS Address deposits on Arbitrum and Avalanch networks

## v27.0.4

### Fixes
* [3508](https://github.com/zeta-chain/node/pull/3508) - fix empty from field in eth receipt rpc method

## v27.0.1

### Fixes

* [3460](https://github.com/zeta-chain/node/pull/3460) - add `group`,`gov`,`params`,`consensus`,`feemarket` ,`crisis`,`vesting` modules to the cosmos interface registry to enable parsing of tx results.

## v27.0.0

### Breaking Changes

* Universal contract calls from Bitcoin and Solana now follow the Protocol Contract V2 workflow.
  * For `depositAndCall` and `call` operations, the `onCall` method is invoked on the Universal Contract from the gateway, replacing the previous behavior where `onCrossChainCall` was triggered by the `systemContract`.
  * The interfaces of both functions remain the same.

### Features

* [3353](https://github.com/zeta-chain/node/pull/3353) - add liquidity cap parameter to ZRC20 creation
* [3357](https://github.com/zeta-chain/node/pull/3357) - cosmos-sdk v.50.x upgrade
* [3358](https://github.com/zeta-chain/node/pull/3358) - register aborted CCTX for Bitcoin inbound that carries insufficient depositor fee
* [3368](https://github.com/zeta-chain/node/pull/3368) - cli command to fetch inbound ballot from inbound hash added to zetatools.
* [3425](https://github.com/zeta-chain/node/pull/3425) - enable inscription parsing on Bitcoin mainnet
* [3332](https://github.com/zeta-chain/node/pull/3332) - implement orchestrator V2. Move BTC observer-signer to V2
* [3360](https://github.com/zeta-chain/node/pull/3360) - update protocol contract imports using consolidated path
* [3349](https://github.com/zeta-chain/node/pull/3349) - implement new bitcoin rpc in zetaclient with improved performance and observability
* [3390](https://github.com/zeta-chain/node/pull/3390) - orchestrator V2: EVM observer-signer
* [3426](https://github.com/zeta-chain/node/pull/3426) - use protocol contracts V2 with Bitcoin deposits
* [3326](https://github.com/zeta-chain/node/pull/3326) - improve error messages for cctx status object
* [3418](https://github.com/zeta-chain/node/pull/3418) - orchestrator V2: TON observer-signer
* [3432](https://github.com/zeta-chain/node/pull/3432) - use protocol contracts V2 with Solana deposits
* [3438](https://github.com/zeta-chain/node/pull/3438) - orchestrator V2: SOl observer-signer. Drop V1.
* [3440](https://github.com/zeta-chain/node/pull/3440) - remove unused method `FilterSolanaInboundEvents`
* [3428](https://github.com/zeta-chain/node/pull/3428) - zetaclient: converge EVM clients.
* [2863](https://github.com/zeta-chain/node/pull/2863) - refactor zetacore to delete matured ballots and add a migration script to remove all old ballots.

### Fixes

* [3416](https://github.com/zeta-chain/node/pull/3416) - add a check for nil gas price in the CheckTxFee function

## v26.0.0

### Features

* [3379](https://github.com/zeta-chain/node/pull/3379) - add Avalanche, Arbitrum and World Chain in chain info

### Fixes

* [3374](https://github.com/zeta-chain/node/pull/3374) - remove minimum rent exempt check for SPL token withdrawals
* [3348](https://github.com/zeta-chain/node/pull/3348) - add support to perform withdraws in ZetaChain `onRevert` call

## v25.0.0

### Features

* [3235](https://github.com/zeta-chain/node/pull/3235) - add /systemtime telemetry endpoint (zetaclient)
* [3317](https://github.com/zeta-chain/node/pull/3317) - add configurable signer latency correction (zetaclient)
* [3320](https://github.com/zeta-chain/node/pull/3320) - add zetaclient minimum version check

### Tests

* [3205](https://github.com/zeta-chain/node/issues/3205) - move Bitcoin revert address test to advanced group to avoid upgrade test failure
* [3254](https://github.com/zeta-chain/node/pull/3254) - rename v2 E2E tests as evm tests and rename old evm tests as legacy
* [3095](https://github.com/zeta-chain/node/pull/3095) - initialize simulation tests for custom zetachain modules
* [3276](https://github.com/zeta-chain/node/pull/3276) - add Solana E2E performance tests and improve Solana outbounds performance
* [3207](https://github.com/zeta-chain/node/pull/3207) - add simulation test operations for all messages in crosschain and observer module

### Refactor

* [3170](https://github.com/zeta-chain/node/pull/3170) - revamp TSS package in zetaclient
* [3291](https://github.com/zeta-chain/node/pull/3291) - revamp zetaclient initialization (+ graceful shutdown)
* [3319](https://github.com/zeta-chain/node/pull/3319) - implement scheduler for zetaclient

### Fixes

* [3206](https://github.com/zeta-chain/node/pull/3206) - skip Solana unsupported transaction version to not block inbound observation
* [3184](https://github.com/zeta-chain/node/pull/3184) - zetaclient should not retry if inbound vote message validation fails
* [3230](https://github.com/zeta-chain/node/pull/3230) - update pending nonces when aborting a cctx through MsgAbortStuckCCTX
* [3225](https://github.com/zeta-chain/node/pull/3225) - use separate database file names for btc signet and testnet4
* [3242](https://github.com/zeta-chain/node/pull/3242) - set the `Receiver` of `MsgVoteInbound` to the address pulled from solana memo
* [3253](https://github.com/zeta-chain/node/pull/3253) - fix solana inbound version 0 queries and move tss keysign prior to relayer key checking
* [3278](https://github.com/zeta-chain/node/pull/3278) - enforce checksum format for asset address in ZRC20
* [3289](https://github.com/zeta-chain/node/pull/3289) - remove all dynamic peer discovery (zetaclient)
* [3314](https://github.com/zeta-chain/node/pull/3314) - update `last_scanned_block_number` metrics more frequently for Solana chain
* [3321](https://github.com/zeta-chain/node/pull/3321) - make crosschain-call with invalid withdraw revert

## v24.0.0

* [3323](https://github.com/zeta-chain/node/pull/3323) - upgrade cosmos sdk to 0.47.15

## v23.0.0

### Features

* [2984](https://github.com/zeta-chain/node/pull/2984) - add Whitelist message ability to whitelist SPL tokens on Solana
* [3091](https://github.com/zeta-chain/node/pull/3091) - improve build reproducability. `make release{,-build-only}` checksums should now be stable.
* [3124](https://github.com/zeta-chain/node/pull/3124) - integrate SPL deposits
* [3134](https://github.com/zeta-chain/node/pull/3134) - integrate SPL tokens withdraw to Solana
* [3088](https://github.com/zeta-chain/node/pull/3088) - add functions to check and withdraw zrc20 as delegation rewards
* [3182](https://github.com/zeta-chain/node/pull/3182) - enable zetaclient pprof server on port 6061

### Tests

* [3075](https://github.com/zeta-chain/node/pull/3075) - ton: withdraw concurrent, deposit & revert.
* [3105](https://github.com/zeta-chain/node/pull/3105) - split Bitcoin E2E tests into two runners for deposit and withdraw
* [3154](https://github.com/zeta-chain/node/pull/3154) - configure Solana gateway program id for E2E tests
* [3188](https://github.com/zeta-chain/node/pull/3188) - add e2e test for v2 deposit and call with swap
* [3151](https://github.com/zeta-chain/node/pull/3151) - add withdraw emissions to e2e tests

### Refactor

* [3118](https://github.com/zeta-chain/node/pull/3118) - zetaclient: remove hsm signer
* [3122](https://github.com/zeta-chain/node/pull/3122) - improve & refactor zetaclientd cli
* [3125](https://github.com/zeta-chain/node/pull/3125) - drop support for header proofs
* [3131](https://github.com/zeta-chain/node/pull/3131) - move app context update from zetacore client
* [3137](https://github.com/zeta-chain/node/pull/3137) - remove chain.Chain from zetaclientd config

### Fixes

* [3117](https://github.com/zeta-chain/node/pull/3117) - register messages for emissions module to legacy amino codec.
* [3041](https://github.com/zeta-chain/node/pull/3041) - replace libp2p public DHT with private gossip peer discovery and connection gater for inbound connections
* [3106](https://github.com/zeta-chain/node/pull/3106) - prevent blocked CCTX on out of gas during omnichain calls
* [3139](https://github.com/zeta-chain/node/pull/3139) - fix config resolution in orchestrator
* [3149](https://github.com/zeta-chain/node/pull/3149) - abort the cctx if dust amount is detected in the revert outbound
* [3155](https://github.com/zeta-chain/node/pull/3155) - fix potential panic in the Bitcoin inscription parsing
* [3162](https://github.com/zeta-chain/node/pull/3162) - skip depositor fee calculation if transaction does not involve TSS address
* [3179](https://github.com/zeta-chain/node/pull/3179) - support inbound trackers for v2 cctx
* [3192](https://github.com/zeta-chain/node/pull/3192) - fix incorrect zContext origin caused by the replacement of 'sender' with 'revertAddress'

## v22.1.2

## Fixes

- [3181](https://github.com/zeta-chain/node/pull/3181) - add lock around pingRTT to prevent crash

## v22.1.1

## Fixes

- [3171](https://github.com/zeta-chain/node/pull/3171) - infinite discovery address leak

## v22.1.0

## Features

- [3028](https://github.com/zeta-chain/node/pull/3028) - whitelist connection gater

## Fixes

- [3041](https://github.com/zeta-chain/node/pull/3041) - replace DHT with private peer discovery
- [3162](https://github.com/zeta-chain/node/pull/3162) - skip depositor fee calculation on irrelevant transactions

## v22.0.2

## Fixes

- [3144](https://github.com/zeta-chain/node/pull/3145) - out of gas on ZetaClient during `onRevert`

## v22.0.1

## Fixes

- [3140](https://github.com/zeta-chain/node/pull/3140) - allow BTC revert with dust amount

## v22.0.0

## Refactor

* [3073](https://github.com/zeta-chain/node/pull/3073) - improve ZETA deposit check with max supply check

## v21.0.0

### Features

* [2633](https://github.com/zeta-chain/node/pull/2633) - support for stateful precompiled contracts
* [2788](https://github.com/zeta-chain/node/pull/2788) - add common importable zetacored rpc package
* [2784](https://github.com/zeta-chain/node/pull/2784) - staking precompiled contract
* [2795](https://github.com/zeta-chain/node/pull/2795) - support restricted address in Solana
* [2861](https://github.com/zeta-chain/node/pull/2861) - emit events from staking precompile
* [2860](https://github.com/zeta-chain/node/pull/2860) - bank precompiled contract
* [2870](https://github.com/zeta-chain/node/pull/2870) - support for multiple Bitcoin chains in the zetaclient
* [2883](https://github.com/zeta-chain/node/pull/2883) - add chain static information for btc signet testnet
* [2907](https://github.com/zeta-chain/node/pull/2907) - derive Bitcoin tss address by chain id and added more Signet static info
* [2911](https://github.com/zeta-chain/node/pull/2911) - add chain static information for btc testnet4
* [2904](https://github.com/zeta-chain/node/pull/2904) - integrate authenticated calls smart contract functionality into protocol
* [2919](https://github.com/zeta-chain/node/pull/2919) - add inbound sender to revert context
* [2957](https://github.com/zeta-chain/node/pull/2957) - enable Bitcoin inscription support on testnet
* [2896](https://github.com/zeta-chain/node/pull/2896) - add TON inbound observation
* [2987](https://github.com/zeta-chain/node/pull/2987) - add non-EVM standard inbound memo package
* [2979](https://github.com/zeta-chain/node/pull/2979) - add fungible keeper ability to lock/unlock ZRC20 tokens
* [3012](https://github.com/zeta-chain/node/pull/3012) - integrate authenticated calls erc20 smart contract functionality into protocol
* [3025](https://github.com/zeta-chain/node/pull/3025) - standard memo for Bitcoin inbound
* [3028](https://github.com/zeta-chain/node/pull/3028) - whitelist connection gater
* [3019](https://github.com/zeta-chain/node/pull/3019) - add ditribute functions to staking precompile
* [3020](https://github.com/zeta-chain/node/pull/3020) - add support for TON withdrawals

### Refactor

* [2749](https://github.com/zeta-chain/node/pull/2749) - fix all lint errors from govet
* [2725](https://github.com/zeta-chain/node/pull/2725) - refactor SetCctxAndNonceToCctxAndInboundHashToCctx to receive tsspubkey as an argument
* [2802](https://github.com/zeta-chain/node/pull/2802) - set default liquidity cap for new ZRC20s
* [2826](https://github.com/zeta-chain/node/pull/2826) - remove unused code from emissions module and add new parameter for fixed block reward amount
* [2890](https://github.com/zeta-chain/node/pull/2890) - refactor `MsgUpdateChainInfo` to accept a single chain, and add `MsgRemoveChainInfo` to remove a chain
* [2899](https://github.com/zeta-chain/node/pull/2899) - remove btc deposit fee v1 and improve unit tests
* [2952](https://github.com/zeta-chain/node/pull/2952) - add error_message to cctx.status
* [3039](https://github.com/zeta-chain/node/pull/3039) - use `btcd` native APIs to handle Bitcoin Taproot address
* [3082](https://github.com/zeta-chain/node/pull/3082) - replace docker-based bitcoin sidecar inscription build with Golang implementation

### Tests

* [2661](https://github.com/zeta-chain/node/pull/2661) - update connector and erc20Custody addresses in tss migration e2e tests
* [2703](https://github.com/zeta-chain/node/pull/2703) - add e2e tests for stateful precompiled contracts
* [2830](https://github.com/zeta-chain/node/pull/2830) - extend staking precompile tests
* [2867](https://github.com/zeta-chain/node/pull/2867) - skip precompiles test for tss migration
* [2833](https://github.com/zeta-chain/node/pull/2833) - add e2e framework for TON blockchain
* [2874](https://github.com/zeta-chain/node/pull/2874) - add support for multiple runs for precompile tests
* [2895](https://github.com/zeta-chain/node/pull/2895) - add e2e test for bitcoin deposit and call
* [2894](https://github.com/zeta-chain/node/pull/2894) - increase gas limit for TSS vote tx
* [2932](https://github.com/zeta-chain/node/pull/2932) - add gateway upgrade as part of the upgrade test
* [2947](https://github.com/zeta-chain/node/pull/2947) - initialize simulation tests
* [3033](https://github.com/zeta-chain/node/pull/3033) - initialize simulation tests for import and export

### Fixes

* [2674](https://github.com/zeta-chain/node/pull/2674) - allow operators to vote on ballots associated with discarded keygen without affecting the status of the current keygen.
* [2672](https://github.com/zeta-chain/node/pull/2672) - check observer set for duplicates when adding a new observer or updating an existing one
* [2735](https://github.com/zeta-chain/node/pull/2735) - fix the outbound tracker blocking confirmation and outbound processing on EVM chains by locally index outbound txs in zetaclient
* [2944](https://github.com/zeta-chain/node/pull/2844) - add tsspubkey to index for tss keygen voting
* [2842](https://github.com/zeta-chain/node/pull/2842) - fix: move interval assignment out of cctx loop in EVM outbound tx scheduler
* [2853](https://github.com/zeta-chain/node/pull/2853) - calling precompile through sc with sc state update
* [2925](https://github.com/zeta-chain/node/pull/2925) - add recover to init chainer to diplay informative message when starting a node from block 1
* [2909](https://github.com/zeta-chain/node/pull/2909) - add legacy messages back to codec for querier backward compatibility
* [3018](https://github.com/zeta-chain/node/pull/3018) - support `DepositAndCall` and `WithdrawAndCall` with empty payload
* [3030](https://github.com/zeta-chain/node/pull/3030) - Avoid storing invalid Solana gateway address in the `SetGatewayAddress`
* [3047](https://github.com/zeta-chain/node/pull/3047) - wrong block hash in subscribe new heads

## v20.0.0

### Features

* [2578](https://github.com/zeta-chain/node/pull/2578) - add Gateway address in protocol contract list
* [2630](https://github.com/zeta-chain/node/pull/2630) - implement `MsgMigrateERC20CustodyFunds` to migrate the funds from the ERC20Custody to a new contracts (to be used for the new ERC20Custody contract for smart contract V2)
* [2578](https://github.com/zeta-chain/node/pull/2578) - Add Gateway address in protocol contract list
* [2594](https://github.com/zeta-chain/node/pull/2594) - Integrate Protocol Contracts V2 in the protocol
* [2634](https://github.com/zeta-chain/node/pull/2634) - add support for EIP-1559 gas fees
* [2597](https://github.com/zeta-chain/node/pull/2597) - Add generic rpc metrics to zetaclient
* [2538](https://github.com/zeta-chain/node/pull/2538) - add background worker routines to shutdown zetaclientd when needed for tss migration
* [2681](https://github.com/zeta-chain/node/pull/2681) - implement `MsgUpdateERC20CustodyPauseStatus` to pause or unpause ERC20 Custody contract (to be used for the migration process for smart contract V2)
* [2644](https://github.com/zeta-chain/node/pull/2644) - add created_timestamp to cctx status
* [2673](https://github.com/zeta-chain/node/pull/2673) - add relayer key importer, encryption and decryption
* [2633](https://github.com/zeta-chain/node/pull/2633) - support for stateful precompiled contracts
* [2751](https://github.com/zeta-chain/node/pull/2751) - add RPC status check for Solana chain
* [2788](https://github.com/zeta-chain/node/pull/2788) - add common importable zetacored rpc package
* [2784](https://github.com/zeta-chain/node/pull/2784) - staking precompiled contract
* [2795](https://github.com/zeta-chain/node/pull/2795) - support restricted address in Solana
* [2825](https://github.com/zeta-chain/node/pull/2825) - add Bitcoin inscriptions support

### Refactor

* [2615](https://github.com/zeta-chain/node/pull/2615) - Refactor cleanup of outbound trackers
* [2855](https://github.com/zeta-chain/node/pull/2855) - disable Bitcoin witness support for mainnet

### Tests

* [2726](https://github.com/zeta-chain/node/pull/2726) - add e2e tests for deposit and call, deposit and revert
* [2821](https://github.com/zeta-chain/node/pull/2821) - V2 protocol contracts migration e2e tests

### Fixes

* [2654](https://github.com/zeta-chain/node/pull/2654) - add validation for authorization list in when validating genesis state for authorization module
* [2672](https://github.com/zeta-chain/node/pull/2672) - check observer set for duplicates when adding a new observer or updating an existing one
* [2824](https://github.com/zeta-chain/node/pull/2824) - fix Solana deposit number

## v19.0.0

### Breaking Changes

* [2460](https://github.com/zeta-chain/node/pull/2460) - Upgrade to go 1.22. This required us to temporarily remove the QUIC backend from [go-libp2p](https://github.com/libp2p/go-libp2p). If you are a zetaclient operator and have configured quic peers, you need to switch to tcp peers.
* [List of the other breaking changes can be found in this document](docs/releases/v19_breaking_changes.md)

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
* [2416](https://github.com/zeta-chain/node/pull/2416) - add Solana chain information
* [2465](https://github.com/zeta-chain/node/pull/2465) - add Solana inbound SOL token observation
* [2497](https://github.com/zeta-chain/node/pull/2416) - support for runtime chain (de)provisioning
* [2518](https://github.com/zeta-chain/node/pull/2518) - add support for Solana address in zetacore
* [2483](https://github.com/zeta-chain/node/pull/2483) - add priorityFee (gasTipCap) gas to the state
* [2567](https://github.com/zeta-chain/node/pull/2567) - add sign latency metric to zetaclient (zetaclient_sign_latency)
* [2524](https://github.com/zeta-chain/node/pull/2524) - add inscription envelop parsing
* [2560](https://github.com/zeta-chain/node/pull/2560) - add support for Solana SOL token withdraw
* [2533](https://github.com/zeta-chain/node/pull/2533) - parse memo from both OP_RETURN and inscription
* [2765](https://github.com/zeta-chain/node/pull/2765) - bitcoin depositor fee improvement

### Refactor

* [2094](https://github.com/zeta-chain/node/pull/2094) - upgrade go-tss to use cosmos v0.47
* [2110](https://github.com/zeta-chain/node/pull/2110) - move non-query rate limiter logic to zetaclient side and code refactor
* [2032](https://github.com/zeta-chain/node/pull/2032) - improve some general structure of the ZetaClient codebase
* [2097](https://github.com/zeta-chain/node/pull/2097) - refactor lightclient verification flags to account for individual chains
* [2071](https://github.com/zeta-chain/node/pull/2071) - Modify chains struct to add all chain related information
* [2118](https://github.com/zeta-chain/node/pull/2118) - consolidate inbound and outbound naming
* [2124](https://github.com/zeta-chain/node/pull/2124) - removed unused variables and method
* [2150](https://github.com/zeta-chain/node/pull/2150) - created `chains` `zetacore` `orchestrator` packages in zetaclient and reorganized source files accordingly
* [2210](https://github.com/zeta-chain/node/pull/2210) - removed unnecessary panics in the zetaclientd process
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
* [2464](https://github.com/zeta-chain/node/pull/2464) - move common voting logic to voting.go and add new function VoteOnBallot
* [2515](https://github.com/zeta-chain/node/pull/2515) - replace chainName by chainID for ChainNonces indexing
* [2541](https://github.com/zeta-chain/node/pull/2541) - deprecate ChainName field in Chain object
* [2542](https://github.com/zeta-chain/node/pull/2542) - adjust permissions to be more restrictive
* [2572](https://github.com/zeta-chain/node/pull/2572) - turn off IBC modules
* [2556](https://github.com/zeta-chain/node/pull/2556) - refactor migrator length check to use consensus type
* [2568](https://github.com/zeta-chain/node/pull/2568) - improve AppContext by converging chains, chainParams, enabledChains, and additionalChains into a single zctx.Chain

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
* [2440](https://github.com/zeta-chain/node/pull/2440) - Add e2e test for TSS migration
* [2473](https://github.com/zeta-chain/node/pull/2473) - add e2e tests for most used admin transactions

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
* [2481](https://github.com/zeta-chain/node/pull/2481) - increase gas limit inbound and outbound vote message to 500k
* [2545](https://github.com/zeta-chain/node/pull/2545) - check solana minimum rent exempt to avoid outbound failure
* [2547](https://github.com/zeta-chain/node/pull/2547) - limit max txs in priority mempool
* [2628](https://github.com/zeta-chain/node/pull/2628) - avoid submitting invalid hashes to outbound tracker

### CI

* [2388](https://github.com/zeta-chain/node/pull/2388) - added GitHub attestations of binaries produced in the release workflow.
* [2285](https://github.com/zeta-chain/node/pull/2285) - added nightly EVM performance testing pipeline, modified localnet testing docker image to utilize debian:bookworm, removed build-jet runners where applicable, removed deprecated/removed upgrade path testing pipeline
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

### Performance

* [2482](https://github.com/zeta-chain/node/pull/2482) - increase the outbound tracker buffer length from 2 to 5

## v18.0.0

* [2470](https://github.com/zeta-chain/node/pull/2470) - add Polygon, Base and Base Sepolia in static chain info

## v17.0.1

### Fixes

* hotfix/v17.0.1 - modify the amount field in CCTXs that carry dust BTC amounts to avoid dust output error

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
* [2076](https://github.com/zeta-chain/node/pull/2076) - automatically deposit native zeta to an address if it doesn't exist on ZEVM
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

* [1817](https://github.com/zeta-chain/node/pull/1817) - Add migration script to fix pending and chain nonces on testnet

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
* [1762](https://github.com/zeta-chain/node/pull/1762) - improve coverage for fungible module
* [1782](https://github.com/zeta-chain/node/pull/1782) - improve coverage for fungible module system contract

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

Getting the correct TSS address for Bitcoin now requires providing the Bitcoin chain id:
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
* fix Athens-3 log print issue - avoid posting unnecessary outtx confirmation
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
* [1504](https://github.com/zeta-chain/node/pull/1504) - remove `-race` in the `make install` command
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
* [1261](https://github.com/zeta-chain/node/pull/1261) - Ethereum comparison checksum/non-checksum format
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
