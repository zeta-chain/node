package testutils

import ethcommon "github.com/ethereum/go-ethereum/common"

const (
	// TSSAddressEVMMainnet the EVM TSS address for test purposes
	// Note: public key is zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc
	TSSAddressEVMMainnet = "0x70e967acFcC17c3941E87562161406d41676FD83"

	// TSSAddressBTCMainnet the BTC TSS address for test purposes
	TSSAddressBTCMainnet = "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y"

	// TSSAddressEVMAthens3 the EVM TSS address for test purposes
	// Note: public key is zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p
	TSSAddressEVMAthens3 = "0x8531a5aB847ff5B22D855633C25ED1DA3255247e"

	// TSSAddressBTCAthens3 the BTC TSS address for test purposes
	TSSAddressBTCAthens3 = "tb1qy9pqmk2pd9sv63g27jt8r657wy0d9ueeh0nqur"

	OtherAddress1 = "0x21248Decd0B7EcB0F30186297766b8AB6496265b"
	OtherAddress2 = "0x33A351C90aF486AebC35042Bb0544123cAed26AB"
	OtherAddress3 = "0x86B77E4fBd07CFdCc486cAe4F2787fB5C5a62cd3"
)

// ConnectorAddresses contains constants ERC20 connector addresses for testing
var ConnectorAddresses = map[int64]ethcommon.Address{
	// Connector address on Ethereum mainnet
	1: ethcommon.HexToAddress("0x000007Cf399229b2f5A4D043F20E90C9C98B7C6a"),

	// Connector address on Binance Smart Chain mainnet
	56: ethcommon.HexToAddress("0x000063A6e758D9e2f438d430108377564cf4077D"),

	// Connector address on Goerli testnet
	5: ethcommon.HexToAddress("0x00005E3125aBA53C5652f9F0CE1a4Cf91D8B15eA"),

	// Connector address on Binance Smart Chain testnet
	97: ethcommon.HexToAddress("0x0000ecb8cdd25a18F12DAA23f6422e07fBf8B9E1"),
}

// CustodyAddresses contains constants ERC20 custody addresses for testing
var CustodyAddresses = map[int64]ethcommon.Address{
	// ERC20 custody address on Ethereum mainnet
	1: ethcommon.HexToAddress("0x0000030Ec64DF25301d8414eE5a29588C4B0dE10"),

	// ERC20 custody address on Binance Smart Chain mainnet
	56: ethcommon.HexToAddress("0x00000fF8fA992424957F97688015814e707A0115"),

	// ERC20 custody address on Goerli testnet
	5: ethcommon.HexToAddress("0x000047f11C6E42293F433C82473532E869Ce4Ec5"),

	// ERC20 custody address on Binance Smart Chain testnet
	97: ethcommon.HexToAddress("0x0000a7Db254145767262C6A81a7eE1650684258e"),
}
