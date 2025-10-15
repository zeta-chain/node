package logs

// A group of predefined field keys and module names for the zetaclient logs.
const (
	// The mode associated with the message.
	// Check mode.Mode for possible values.
	FieldMode = "mode"

	// The module in the ZetaClient that originated the message.
	// Check the ModName constants for possible values.
	FieldModule = "module"

	// Chain identifier.
	FieldChain = "chain"

	// Chain network.
	FieldNetwork = "network"

	// Height of a block.
	FieldBlock = "block"

	// Nonce of a transaction.
	FieldNonce = "nonce"

	// Hash of a transaction from an external chain.
	FieldTx = "tx"

	// Hash of a ZetaChain transaction.
	FieldZetaTx = "zeta_tx"

	// Unique identifier (index) of a cross-chain transaction.
	FieldCctxIndex = "cctx_index"

	// Unique identifier of an outbound transaction.
	// Usually contains the nonce of the outbound in the external chain plus the chain identifier.
	FieldOutboundID = "outbound_id"

	// Index of a ballot from voting on an inbound or outbound cross-chain transaction.
	FieldBallotIndex = "ballot_index"

	// Type of the asset in the cross-chain transaction.
	// This field can take one of the following values: Gas, ERC20, ZETA, and NoAssetCall.
	FieldCoinType = "coin_type"

	// TXID from a Bitcoin transaction.
	FieldBtcTxid = "txid"
)

// Range of values for FieldModule.
const (
	ModNameOrchestrator   = "orchestrator"
	ModNameInbound        = "inbound"
	ModNameOutbound       = "outbound"
	ModNameSigner         = "signer"
	ModNameGasPrice       = "gas_price"
	ModNameBtcClient      = "btc_client"
	ModNameZetaCoreClient = "zetacore_client"

	ModNameTssHealthCheck = "tss_healthcheck"
	ModNameTssKeyGen      = "tss_keygen"
	ModNameTssKeySign     = "tss_keysign"
	ModNameTssService     = "tss_service"
	ModNameTssSetup       = "tss_setup"
)
