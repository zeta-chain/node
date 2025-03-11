package common

const (
	// BlocksPerDay is the average number of Bitcoin blocks produced per day.
	BlocksPerDay = 144

	// TSSSignatureCacheSize is the size of the TSS signature cache for Bitcoin.
	// In Bitcoin, each UTXO requires a signature, and TSS may own hundreds of UTXOs,
	// so we need a bigger cache size than other chains.
	TSSSignatureCacheSize = 1000
)
