
# V19 Observer Migration Guide

## Authorizations Update

Observers use the `authz` module to grant authorization to execute messages with their hotkey.

The following messages using `authz` have been updated:
- `AddToOutTxTracker` to `AddOutboundTracker`
- `AddToInTxTracker` to `AddInboundTracker`
- `VoteOnObservedOutboundTx` to `VoteOutbound`
- `VoteOnObservedInboundTx` to `VoteInbound`

In consequence, the observers must manually add the authorization for the new messages.

The authorizations can be added by interacting with the `authz` module. In this example with the `zetacored` CLI:

The current authorization grants can be listed with the following command:
```bash
zetacored q authz grants-by-grantee [operator_address]
```

To add the authorization for the new messages, use the following command:
```bash
zetacored tx authz grant [grantee_address] generic --msg-type=/zetachain.zetacore.crosschain.MsgVoteInbound
zetacored tx authz grant [grantee_address] generic --msg-type=/zetachain.zetacore.crosschain.MsgVoteOutbound
zetacored tx authz grant [grantee_address] generic --msg-type=/zetachain.zetacore.crosschain.MsgAddOutboundTracker
zetacored tx authz grant [grantee_address] generic --msg-type=/zetachain.zetacore.crosschain.MsgAddInboundTracker
```
