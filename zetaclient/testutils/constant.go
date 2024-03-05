package testutils

import ethcommon "github.com/ethereum/go-ethereum/common"

const (
	// tss addresses
	TSSAddressEVMMainnet = "0x70e967acFcC17c3941E87562161406d41676FD83"
	TSSAddressBTCMainnet = "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y"

	TSSAddressEVMAthens3 = "0x8531a5aB847ff5B22D855633C25ED1DA3255247e"
	TSSAddressBTCAthens3 = "tb1qy9pqmk2pd9sv63g27jt8r657wy0d9ueeh0nqur"

	// some other address
	OtherAddress = "0x21248Decd0B7EcB0F30186297766b8AB6496265b"
)

// ConnectorAddresses contains constants ERC20 connector addresses for testing
var ConnectorAddresses = map[int64]ethcommon.Address{
	// mainnet
	1:  ethcommon.HexToAddress("0x000007Cf399229b2f5A4D043F20E90C9C98B7C6a"),
	56: ethcommon.HexToAddress("0x000063A6e758D9e2f438d430108377564cf4077D"),

	// testnet
	5:  ethcommon.HexToAddress("0x00005E3125aBA53C5652f9F0CE1a4Cf91D8B15eA"),
	97: ethcommon.HexToAddress("0x0000ecb8cdd25a18F12DAA23f6422e07fBf8B9E1"),
}

// CustodyAddresses contains constants ERC20 custody addresses for testing
var CustodyAddresses = map[int64]ethcommon.Address{
	// mainnet
	1:  ethcommon.HexToAddress("0x0000030Ec64DF25301d8414eE5a29588C4B0dE10"),
	56: ethcommon.HexToAddress("0x00000fF8fA992424957F97688015814e707A0115"),

	// testnet
	5:  ethcommon.HexToAddress("0x000047f11C6E42293F433C82473532E869Ce4Ec5"),
	97: ethcommon.HexToAddress("0x0000a7Db254145767262C6A81a7eE1650684258e"),
}
