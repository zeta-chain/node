package common

const (
	// DefaultGasPriceMultiplier is the default gas price multiplier for outbond txs
	DefaultGasPriceMultiplier = 2

	// DonationMessage is the message for donation transactions
	// Transaction sent to the TSS or ERC20 Custody address containing this message are considered as a donation
	DonationMessage = "I am rich!"
)
