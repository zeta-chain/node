package local

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
)

var (
	// DeployerAddress is the address of the account for deploying networks
	DeployerAddress    = ethcommon.HexToAddress("0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC")
	DeployerPrivateKey = "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263" // #nosec G101 - used for testing

	// UserERC20Address is the address of the account for testing ERC20
	UserERC20Address    = ethcommon.HexToAddress("0x6F57D5E7c6DBb75e59F1524a3dE38Fc389ec5Fd6")
	UserERC20PrivateKey = "fda3be1b1517bdf48615bdadacc1e6463d2865868dc8077d2cdcfa4709a16894" // #nosec G101 - used for testing

	// UserZetaTestAddress is the address of the account for testing Zeta
	UserZetaTestAddress    = ethcommon.HexToAddress("0x5cC2fBb200A929B372e3016F1925DcF988E081fd")
	UserZetaTestPrivateKey = "729a6cdc5c925242e7df92fdeeb94dadbf2d0b9950d4db8f034ab27a3b114ba7" // #nosec G101 - used for testing

	// UserBitcoinAddress is the address of the account for testing Bitcoin
	UserBitcoinAddress    = ethcommon.HexToAddress("0x283d810090EdF4043E75247eAeBcE848806237fD")
	UserBitcoinPrivateKey = "7bb523963ee2c78570fb6113d886a4184d42565e8847f1cb639f5f5e2ef5b37a" // #nosec G101 - used for testing

	// UserEtherAddress is the address of the account for testing Ether
	UserEtherAddress    = ethcommon.HexToAddress("0x8D47Db7390AC4D3D449Cc20D799ce4748F97619A")
	UserEtherPrivateKey = "098e74a1c2261fa3c1b8cfca8ef2b4ff96c73ce36710d208d1f6535aef42545d" // #nosec G101 - used for testing

	// UserMiscAddress is the address of the account for miscellaneous tests
	UserMiscAddress    = ethcommon.HexToAddress("0x90126d02E41c9eB2a10cfc43aAb3BD3460523Cdf")
	UserMiscPrivateKey = "853c0945b8035a501b1161df65a17a0a20fc848bda8975a8b4e9222cc6f84cd4" // #nosec G101 - used for testing

	// UserAdminAddress is the address of the account for testing admin function features
	UserAdminAddress    = ethcommon.HexToAddress("0xcC8487562AAc220ea4406196Ee902C7c076966af")
	UserAdminPrivateKey = "95409f1f0e974871cc26ba98ffd31f613aa1287d40c0aea6a87475fc3521d083" // #nosec G101 - used for testing

	FungibleAdminMnemonic = "snow grace federal cupboard arrive fancy gym lady uniform rotate exercise either leave alien grass" // #nosec G101 - used for testing
)
