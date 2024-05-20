
# V17 Breaking Changes

### Verification Flags update

* `MsgUpdateVerificationFlags` has been removed, and replaced with `MsgEnableHeaderVerification` and `MsgDisableHeaderVerification` messages.
    * `MsgEnableHeaderVerification` message enables block header verification for a list of chains and can be triggered via `PolicyType_groupOperational`
    * `MsgDisableHeaderVerification` message disables block header verification for a list of chains and can be triggered via `PolicyType_emergency`

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
