package pkg

const (
	// DonationMessage is the message for donation transactions
	// Transaction sent to the TSS or ERC20 Custody address containing this message are considered as a donation
	DonationMessage = "I am rich!"
	// CmdWhitelistERC20 is used for CCTX of type cmd to give the instruction to the TSS to whitelist an ERC20 on an exeternal chain
	CmdWhitelistERC20 = "cmd_whitelist_erc20"
	// CmdMigrateTssFunds is used for CCTX of type cmd to give the instruction to the TSS to transfer its funds on a new address
	CmdMigrateTssFunds = "cmd_migrate_tss_funds"
)
