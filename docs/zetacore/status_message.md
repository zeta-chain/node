# CCTX status message 

The cctx object contains a field  `cctx_status` , which has the following structure 
```go
type Status struct {
	Status CctxStatus `protobuf:"varint,1,opt,name=status,proto3,enum=zetachain.zetacore.crosschain.CctxStatus" json:"status,omitempty"`
	StatusMessage string `protobuf:"bytes,2,opt,name=status_message,json=statusMessage,proto3" json:"status_message,omitempty"`
	ErrorMessage        string `protobuf:"bytes,6,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	LastUpdateTimestamp int64  `protobuf:"varint,3,opt,name=lastUpdate_timestamp,json=lastUpdateTimestamp,proto3" json:"lastUpdate_timestamp,omitempty"`
	IsAbortRefunded     bool   `protobuf:"varint,4,opt,name=isAbortRefunded,proto3" json:"isAbortRefunded,omitempty"`
	CreatedTimestamp int64 `protobuf:"varint,5,opt,name=created_timestamp,json=createdTimestamp,proto3" json:"created_timestamp,omitempty"`
	ErrorMessageRevert string `protobuf:"bytes,7,opt,name=error_message_revert,json=errorMessageRevert,proto3" json:"error_message_revert,omitempty"`
}
```

## Status 
This is the most updated status for the cctx . This can be one of the following values
- `PendingInbound` : The cctx is pending for the inbound to be finalized , this is an intermediately status used by the protocol only
- `PendingOutbound` : This means that the inbound has been finalzied, and the outbound is pending
- `OutboundMined` : The outbound has been successfully mined. This is a terminal status
- `Aborted` : The cctx has been aborted. This is a terminal status
- `PendingRevert` : The the cctx failed at some step and is pending for the revert to be finalized
- `Reverted` : The cctx has been successfully reverted. This is a terminal status

### StatusMessage
The status message provides a some details about the current status.This is primiary meant for the user to quickly understand the status of the cctx.
### LastUpdateTimestamp
The last time the status was updated
### IsAbortRefunded
This is a boolean value which is true if the cctx has been refunded after being aborted or not .
### CreatedTimestamp
The time when the cctx was created
### ErrorMessage and ErrorMessageRevert
A cctx can have a maximum of two outbound params. We can refer to the first outbound as `outbound` and the second as `revert`.
- A normal flow for a cctx is to go from `PendingOutbound` -> `OutboundMined` , which creates a single outbound
- A cctx where the outbound fails has the transition `PendingOutbound` -> `PendingRevert` -> `Reverted` , which creates two outbounds
- Any of the above two flows can abort the cctx at some point that can create either one or two outbounds

  - The `ErrorMessage` field only contains a value if the original outbound failed. It contains details about the error that caused the outbound to fail
  - The `ErrorMessageRevert` field only contains a value if the revert outbound failed. It contains details about the error that caused the revert outbound to fail.

### Example values for StatusMessage field and how to interpret them
- `initiating outbound` : The inbound votes have been successfully finalized, and the protocol is starting the outbound process
- `outbound successfully mined` : The outbound was successfully mined
- `revert successfully processed` : The outbound failed , but the revert was successful
- `revert failed to be processed` : The revert failed. This message also means that the initial outbound has failed.

- `outbound failed` : The outbound failed, The protocol would try to create a revert either in the same block or schedule one to be picked up by zetaclient
- `outbound failed for admin tx` : The outbound failed for an admin transaction, in this case we do not revert the cctx
- `outbound failed unable to process` : The outbound processing failed at the protocol level. When this happens, the protocol sets the cctx to aborted.
- `outbound failed but the universal contract did not revert` :  The outbound/deposit failed, but the contract did not revert,
   this is most likely caused by an internal error in the protocol.The CCTX is this case is aborted. Users can try connecting with the zetachain team to get a refund
- `cctx aborted through MsgAbortStuckCCTX` : The cctx was aborted manually by an admin command


### Example values for ErrorMessage and ErrorMessageRevert fields and how to interpret them

- For a failed deposit, the ErrorMessage would contain the following fields. The protocol generates the fields tagged as internal.
```
  - method: The method that was called by the protocol
  - contract: The contract that his method was called on
  - args:The argumets that were used for this call
  - errorMessage[Internal]: Error message from the ZEVM call
  - revertReason: Revert reason from the smart contract
```

- `outbound tx failed to be executed on connected chain` : `revert tx failed to be executed on connected chain` : The outbound/revert transaction failed to be executed on the connected chain.
- `coin type [CoinType] not supported for revert when source chain is Zetachain` : The coin type is not supported for revert when the source chain is Zetachain.
- `error from EVMDeposit: [Error_String]` : Error returned by the protocol when trying to deposit tokens( and optionally call a contract) on ZEVM. The error string should explain the cause 
