# Messages

## MsgAddObserver

AddObserver adds an observer address to the observer set

```proto
message MsgAddObserver {
	string creator = 1;
	string observer_address = 2;
	string zetaclient_grantee_pubkey = 3;
	bool add_node_account_only = 4;
}
```

## MsgUpdateObserver

UpdateObserver handles updating an observer address
Authorized: admin policy (admin update), old observer address (if the
reason is that the observer was tombstoned).

```proto
message MsgUpdateObserver {
	string creator = 1;
	string old_observer_address = 2;
	string new_observer_address = 3;
	ObserverUpdateReason update_reason = 4;
}
```

## MsgUpdateChainParams

UpdateChainParams updates chain parameters for a specific chain, or add a new one.
Chain parameters include: confirmation count, outbound transaction schedule interval, ZETA token,
connector and ERC20 custody contract addresses, etc.
Only the admin policy account is authorized to broadcast this message.

```proto
message MsgUpdateChainParams {
	string creator = 1;
	ChainParams chainParams = 2;
}
```

## MsgRemoveChainParams

RemoveChainParams removes chain parameters for a specific chain.

```proto
message MsgRemoveChainParams {
	string creator = 1;
	int64 chain_id = 2;
}
```

## MsgAddBlameVote

```proto
message MsgAddBlameVote {
	string creator = 1;
	int64 chain_id = 2;
	Blame blame_info = 3;
}
```

## MsgUpdateKeygen

UpdateKeygen updates the block height of the keygen and sets the status to
"pending keygen".

Authorized: admin policy group 1.

```proto
message MsgUpdateKeygen {
	string creator = 1;
	int64 block = 2;
}
```

## MsgVoteBlockHeader

VoteBlockHeader vote for a new block header to the storers

```proto
message MsgVoteBlockHeader {
	string creator = 1;
	int64 chain_id = 2;
	bytes block_hash = 3;
	int64 height = 4;
	pkg.proofs.HeaderData header = 5;
}
```

## MsgResetChainNonces

ResetChainNonces handles resetting chain nonces

```proto
message MsgResetChainNonces {
	string creator = 1;
	int64 chain_id = 2;
	int64 chain_nonce_low = 3;
	int64 chain_nonce_high = 4;
}
```

## MsgVoteTSS

VoteTSS votes on creating a TSS key and recording the information about it (public
key, participant and operator addresses, finalized and keygen heights).

If the vote passes, the information about the TSS key is recorded on chain
and the status of the keygen is set to "success".

Fails if the keygen does not exist, the keygen has been already
completed, or the keygen has failed.

Only node accounts are authorized to broadcast this message.

```proto
message MsgVoteTSS {
	string creator = 1;
	string tss_pubkey = 2;
	int64 keygen_zeta_height = 3;
	pkg.chains.ReceiveStatus status = 4;
}
```

## MsgEnableCCTXFlags

```proto
message MsgEnableCCTXFlags {
	string creator = 1;
	bool enableInbound = 2;
	bool enableOutbound = 3;
}
```

## MsgDisableCCTXFlags

```proto
message MsgDisableCCTXFlags {
	string creator = 1;
	bool disableInbound = 2;
	bool disableOutbound = 3;
}
```

## MsgUpdateGasPriceIncreaseFlags

```proto
message MsgUpdateGasPriceIncreaseFlags {
	string creator = 1;
	GasPriceIncreaseFlags gasPriceIncreaseFlags = 2;
}
```

