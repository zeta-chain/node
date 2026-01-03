package constant

import "time"

const (
	// ZetaBlockTime is the block time of the ZetaChain network
	// It's a rough estimate that can be used in non-critical path to estimate the time of a block
	ZetaBlockTime = 2 * time.Second

	// DonationMessage is the message for donation transactions
	// Transaction sent to the TSS or ERC20 Custody address containing this message are considered as a donation
	DonationMessage = "I am rich!"

	// CmdWhitelistAsset is used for CCTX of type cmd to give the instruction to the TSS to whitelist an ERC20 on an exeternal chain
	CmdWhitelistAsset = "cmd_whitelist_asset"

	// CmdMigrateTSSFunds is used for CCTX of type cmd to give the instruction to the TSS to transfer its funds on a new address
	CmdMigrateTSSFunds = "cmd_migrate_tss_funds"
	// BTCWithdrawalDustAmount is the minimum satoshis that can be withdrawn from zEVM to avoid outbound dust output
	// The Bitcoin protocol sets a minimum output value to 546 satoshis (dust limit) but we set it to 1000 satoshis
	BTCWithdrawalDustAmount = 1000

	// SolanaWalletRentExempt is the minimum balance for a Solana wallet account to become rent exempt
	// The Solana protocol sets minimum rent exempt to 890880 lamports but we set it to 1_000_000 lamports (0.001 SOL)
	// The number 890880 comes from CLI command `solana rent 0` and has been verified on devnet gateway program
	SolanaWalletRentExempt = 1_000_000

	// EVMZeroAddress is the zero address for EVM address format
	EVMZeroAddress = "0x0000000000000000000000000000000000000000"

	// OptionPause is the argument used in CmdUpdateERC20CustodyPauseStatus to pause the ERC20 custody contract
	OptionPause = "pause"

	// OptionUnpause is the argument used in CmdUpdateERC20CustodyPauseStatus to unpause the ERC20 custody contract
	OptionUnpause = "unpause"

	// DefaultAppMempoolSize is the default size of ZetaChain mempool
	DefaultAppMempoolSize = 3000
	// TODO : Check if they are used
	// CmdWhitelistERC20 is used for CCTX of type cmd to give the instruction to the TSS to whitelist an ERC20 on an exeternal chain

	CmdWhitelistERC20 = "cmd_whitelist_erc20"

	// CmdMigrateERC20CustodyFunds is used for CCTX of type cmd to give the instruction to the TSS to transfer its funds on a new address
	CmdMigrateERC20CustodyFunds = "cmd_migrate_erc20_custody_funds"

	// CmdUpdateERC20CustodyPauseStatus is used for CCTX of type cmd to give the instruction to the TSS to update the pause status of the ERC20 custody contract
	CmdUpdateERC20CustodyPauseStatus = "cmd_update_erc20_custody_pause_status"
)
