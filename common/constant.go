package common

const (
	// EVMOuttxGasPriceMultiplier is the default gas price multiplier for EVM-chain outbond txs
	EVMOuttxGasPriceMultiplier = 1.2

	// BTCOuttxGasPriceMultiplier is the default gas price multiplier for BTC outbond txs
	BTCOuttxGasPriceMultiplier = 2.0

	// DonationMessage is the message for donation transactions
	// Transaction sent to the TSS or ERC20 Custody address containing this message are considered as a donation
	DonationMessage = "I am rich!"
)
