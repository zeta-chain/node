# Messages

## MsgAddToOutTxTracker

The AddToOutTxTracker function adds a transaction hash and signer to the
outgoing transaction tracker for a given chain and nonce. It first checks if
the message creator is a bonded validator or an admin key. It then retrieves
a chain based on the chain ID and checks if it is supported. After that, it
retrieves a tracker based on the chain ID and the nonce. If the tracker does
not exist, it creates a new tracker with the provided hash and signer. If the
tracker exists, it checks if the hash already exists in the tracker. If not,
it adds the hash and signer to the tracker. The function returns an empty
MsgAddToOutTxTrackerResponse and no error.

```proto
message MsgAddToOutTxTracker {
	string creator = 1;
	int64 chain_id = 2;
	uint64 nonce = 3;
	string tx_hash = 4;
}
```

## MsgRemoveFromOutTxTracker

The RemoveFromOutTxTracker function removes an outgoing transaction tracker
for a given chain and nonce. It first checks if the message creator is a
bonded validator or an admin key. It then removes the tracker with the
provided chain ID and nonce using the RemoveOutTxTracker function. The
function returns an empty MsgRemoveFromOutTxTrackerResponse and no error.

```proto
message MsgRemoveFromOutTxTracker {
	string creator = 1;
	int64 chain_id = 2;
	uint64 nonce = 3;
}
```

## MsgCreateTSSVoter

The CreateTSSVoter function creates a threshold signature scheme (TSS) voter
and adds it to the TSS voter store. It first checks if the message creator is
a bonded validator. It then calculates the sessionID based on the current
block height and creates an index using the message digest and the sessionID.
It retrieves a TSS voter based on the index and checks if the creator has
already signed. If the creator has not signed, the method adds the creator to
the Signers list in the TSS voter. If the TSS voter is not found, the method
creates a new TSS voter with the provided information and initializes the
Signers list with the creator. The method then sets the TSS voter in the
store using the SetTSSVoter function. If the Signers list in the TSS voter is
equal to the number of validators, the method creates a new TSS using the TSS
voter information and sets it in the TSS store using the SetTSS function. The
function returns an empty MsgCreateTSSVoterResponse and no error.

```proto
message MsgCreateTSSVoter {
	string creator = 1;
	string chain = 3;
	string address = 4;
	string pubkey = 5;
}
```

## MsgGasPriceVoter

```proto
message MsgGasPriceVoter {
	string creator = 1;
	int64 chain_id = 2;
	uint64 price = 3;
	uint64 block_number = 4;
	string supply = 5;
}
```

## MsgNonceVoter

```proto
message MsgNonceVoter {
	string creator = 1;
	int64 chain_id = 2;
	uint64 nonce = 3;
}
```

## MsgVoteOnObservedOutboundTx

```proto
message MsgVoteOnObservedOutboundTx {
	string creator = 1;
	string cctx_hash = 2;
	string observed_outTx_hash = 3;
	uint64 observed_outTx_blockHeight = 4;
	string zeta_minted = 5;
	common.ReceiveStatus status = 6;
	int64 outTx_chain = 7;
	uint64 outTx_tss_nonce = 8;
	common.CoinType coin_type = 9;
}
```

## MsgVoteOnObservedInboundTx

FIXME: use more specific error types & codes

```proto
message MsgVoteOnObservedInboundTx {
	string creator = 1;
	string sender = 2;
	int64 sender_chain_id = 3;
	string receiver = 4;
	int64 receiver_chain = 5;
	string amount = 6;
	string message = 8;
	string in_tx_hash = 9;
	uint64 in_block_height = 10;
	uint64 gas_limit = 11;
	common.CoinType coin_type = 12;
	string tx_origin = 13;
	string asset = 14;
}
```

## MsgSetNodeKeys

```proto
message MsgSetNodeKeys {
	string creator = 1;
	common.PubKeySet pubkeySet = 2;
	string tss_signer_Address = 3;
}
```

## MsgUpdatePermissionFlags

```proto
message MsgUpdatePermissionFlags {
	string creator = 1;
	bool isInboundEnabled = 3;
}
```

