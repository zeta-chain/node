
# V19 Breaking Changes

### Verification Flags update

* `MsgUpdateVerificationFlags` has been removed, and replaced with `MsgEnableHeaderVerification` and `MsgDisableHeaderVerification` messages.
    * `MsgEnableHeaderVerification` message enables block header verification for a list of chains and can be triggered via `PolicyType_groupOperational`
    * `MsgDisableHeaderVerification` message disables block header verification for a list of chains and can be triggered via `PolicyType_emergency`

### Crosschain Flags update

* `MsgUpdateCrosschainFlags` has been removed,and replaced with `MsgEnableCCTX`, `MsgDisableCCTX` and `MsgUpdateGasPriceIncreaseFlags` messages.
    * `MsgEnableCCTX` message enables either the IsInboundEnabled flag,or the IsOutboundEnabled flag or both `PolicyType_groupOperational`
    * `MsgDisableCCTX` message disables either the IsInboundEnabled flag,or the IsOutboundEnabled flag or both `PolicyType_emergency`
    * `MsgUpdateGasPriceIncreaseFlags` message updates the gas price increase flags and can be triggered via `PolicyType_groupOperational`

### `BallotMaturityBlocks` moved to `emissions` module

* Observer param `ballot_maturity_blocks` is part of `emissions` module now. Observer `params` are deprecated and removed from `observer` module.

### `InTx` and `OutTx` renaming

* All references to inTx and outTx have been replaced with `inbound` and `outbound` respectively. In consequence several structures, messages and queries have been renamed to reflect this change.
    * Structure renaming:
        * `InTxHashToCctx` has been renamed to `InboundHashToCctx`
            * Field `InTxHash` has been renamed to `InboundHash`
        * `InTxTracker` has been renamed to `InboundTracker`
        * `OutTxTracker` has been renamed to `OutboundTracker`
        * In `ChainParams`:
            * `InTxTracker` has been renamed to `InboundTracker`
            * `OutTxTracker` has been renamed to `OutboundTracker`
            * `OutboundTxScheduleInterval` has been renamed to `OutboundScheduleInterval`
            * `OutboundTxScheduleLookahead` has been renamed to `OutboundScheduleLookahead`
    * Messages
        * `AddToOutTxTracker` has been renamed to `AddOutboundTracker`
        * `AddToInTxTracker` has been renamed to `AddInboundTracker`
        * `RemoveFromOutTxTracker` has been renamed to `RemoveOutboundTracker`
        * `VoteOnObservedOutboundTx` has been renamed to `VoteOutbound`
        * `VoteOnObservedInboundTx` has been renamed to `VoteInbound`
    * The previous queries have not been removed but have been deprecated and replaced with new queries:
        * `OutTxTracker` has been renamed to `OutboundTracker`
            * `/zeta-chain/crosschain/outTxTracker/{chainID}/{nonce}` endpoint is now `/zeta-chain/crosschain/outboundTracker/{chainID}/{nonce}`
        * `OutTxTrackerAll` has been renamed to `OutboundTrackerAll`
            * `/zeta-chain/crosschain/outTxTracker` endpoint is now `/zeta-chain/crosschain/outboundTracker`
        * `OutTxTrackerAllByChain` has been renamed to `OutboundTrackerAllByChain`
            * `/zeta-chain/crosschain/outTxTrackerByChain/{chainID}" endpoint is now /zeta-chain/crosschain/outboundTrackerByChain/{chainID}`
        * `InTxTrackerAllByChain` has been renamed to `InboundTrackerAllByChain`
            * `/zeta-chain/crosschain/inTxTrackerByChain/{chainID}` endpoint is now `/zeta-chain/crosschain/inboundTrackerByChain/{chainID}`
        * `InTxTrackerAll` has been renamed to `InboundTrackerAll`
            * `/zeta-chain/crosschain/inTxTracker` endpoint is now `/zeta-chain/crosschain/inboundTracker`
        * `InTxHashToCctx` has been renamed to `InboundHashToCctx`
            * `/zeta-chain/crosschain/inTxHashToCctx/{hash}` endpoint is now `/zeta-chain/crosschain/inboundHashToCctx/{hash}`
        * `InTxHashToCctxData` has been renamed to `InboundHashToCctxData`
            * `/zeta-chain/crosschain/inTxHashToCctxData/{hash}` endpoint is now `/zeta-chain/crosschain/inboundHashToCctxData/{hash}`
        * `InTxHashToCctxAll` has been renamed to `InboundHashToCctxAll`
            * `/zeta-chain/crosschain/inTxHashToCctx` endpoint is now `/zeta-chain/crosschain/inboundHashToCctx`
          
* `MsgUpdateZRC20` has been removed, and replaced with `MsgPauseZRC20` and `MsgUnpauseZRC20` messages.
    * `MsgPauseZRC20` message pauses a ZRC20 token and can be triggered via `PolicyType_groupEmergency`
    * `MsgUnpauseZRC20` message unpauses a ZRC20 token and can be triggered via `PolicyType_groupOperational`

### `MsgAddBlameVote` renaming

* `MsgAddBlameVote` has been renamed to `MsgVoteBlame` to maintain consistency with other voting messages

### `Chain.ChainName` deprecated

* `Chain.ChainName` has been deprecated and will be removed from the `Chain` structure. The `Chain.Name` should be used instead.