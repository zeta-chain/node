package e2etests

import (
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// List of all e2e test names to be used in zetae2e
// TODO : E2E smoke test for abort refund
// https://github.com/zeta-chain/node/issues/1745
const (
	/*
	 ZETA tests
	 Test transfer of ZETA asset across chains
	*/
	TestZetaDepositName           = "zeta_deposit"
	TestZetaDepositNewAddressName = "zeta_deposit_new_address"
	TestZetaDepositRestrictedName = "zeta_deposit_restricted"
	TestZetaWithdrawName          = "zeta_withdraw"
	TestZetaWithdrawBTCRevertName = "zeta_withdraw_btc_revert" // #nosec G101 - not a hardcoded password

	/*
	 Message passing tests
	 Test message passing across chains
	*/
	TestMessagePassingExternalChainsName              = "message_passing_external_chains"
	TestMessagePassingRevertFailExternalChainsName    = "message_passing_revert_fail"
	TestMessagePassingRevertSuccessExternalChainsName = "message_passing_revert_success"
	TestMessagePassingEVMtoZEVMName                   = "message_passing_evm_to_zevm"
	TestMessagePassingZEVMToEVMName                   = "message_passing_zevm_to_evm"
	TestMessagePassingZEVMtoEVMRevertName             = "message_passing_zevm_to_evm_revert"
	TestMessagePassingEVMtoZEVMRevertName             = "message_passing_evm_to_zevm_revert"
	TestMessagePassingZEVMtoEVMRevertFailName         = "message_passing_zevm_to_evm_revert_fail"
	TestMessagePassingEVMtoZEVMRevertFailName         = "message_passing_evm_to_zevm_revert_fail"

	/*
	 EVM gas tests
	 Test transfer of EVM gas asset across chains
	*/
	TestEtherDepositName              = "eth_deposit"
	TestEtherWithdrawName             = "eth_withdraw"
	TestEtherWithdrawRestrictedName   = "eth_withdraw_restricted"
	TestEtherDepositAndCallRefundName = "eth_deposit_and_call_refund"
	TestEtherDepositAndCallName       = "eth_deposit_and_call"

	/*
	 EVM erc20 tests
	 Test transfer of EVM erc20 asset across chains
	*/
	TestERC20WithdrawName             = "erc20_withdraw"
	TestERC20DepositName              = "erc20_deposit"
	TestMultipleERC20DepositName      = "erc20_multiple_deposit"
	TestMultipleERC20WithdrawsName    = "erc20_multiple_withdraw"
	TestERC20DepositAndCallRefundName = "erc20_deposit_and_call_refund"
	TestERC20DepositRestrictedName    = "erc20_deposit_restricted" // #nosec G101: Potential hardcoded credentials (gosec), not a credential

	/*
	 Bitcoin tests
	 Test transfer of Bitcoin asset across chains
	*/
	TestBitcoinDepositName                = "bitcoin_deposit"
	TestBitcoinWithdrawSegWitName         = "bitcoin_withdraw_segwit"
	TestBitcoinWithdrawTaprootName        = "bitcoin_withdraw_taproot"
	TestBitcoinWithdrawLegacyName         = "bitcoin_withdraw_legacy"
	TestBitcoinWithdrawP2WSHName          = "bitcoin_withdraw_p2wsh"
	TestBitcoinWithdrawP2SHName           = "bitcoin_withdraw_p2sh"
	TestBitcoinWithdrawInvalidAddressName = "bitcoin_withdraw_invalid"
	TestBitcoinWithdrawRestrictedName     = "bitcoin_withdraw_restricted"

	/*
	 Application tests
	 Test various smart contract applications across chains
	*/
	TestZRC20SwapName      = "zrc20_swap"
	TestCrosschainSwapName = "crosschain_swap"

	/*
	 Miscellaneous tests
	 Test various functionalities not related to assets
	*/
	TestContextUpgradeName = "context_upgrade"
	TestMyTestName         = "my_test"
	TestDonationEtherName  = "donation_ether"

	/*
	 Stress tests
	 Test stressing networks with many cross-chain transactions
	*/
	TestStressEtherWithdrawName = "stress_eth_withdraw"
	TestStressBTCWithdrawName   = "stress_btc_withdraw"
	TestStressEtherDepositName  = "stress_eth_deposit"
	TestStressBTCDepositName    = "stress_btc_deposit"

	/*
	 Admin tests
	 Test admin functionalities
	*/
	TestDepositEtherLiquidityCapName = "deposit_eth_liquidity_cap"
	TestMigrateChainSupportName      = "migrate_chain_support"
	TestPauseZRC20Name               = "pause_zrc20"
	TestUpdateBytecodeZRC20Name      = "update_bytecode_zrc20"
	TestUpdateBytecodeConnectorName  = "update_bytecode_connector"
	TestRateLimiterName              = "rate_limiter"
)

// AllE2ETests is an ordered list of all e2e tests
var AllE2ETests = []runner.E2ETest{
	/*
	 ZETA tests
	*/
	runner.NewE2ETest(
		TestZetaDepositName,
		"deposit ZETA from Ethereum to ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "1000000000000000000"},
		},
		TestZetaDeposit,
	),
	runner.NewE2ETest(
		TestZetaDepositNewAddressName,
		"deposit ZETA from Ethereum to a new ZEVM address which does not exist yet",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "1000000000000000000"},
		},
		TestZetaDepositNewAddress,
	),
	runner.NewE2ETest(
		TestZetaDepositRestrictedName,
		"deposit ZETA from Ethereum to ZEVM restricted address",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "1000000000000000000"},
		},
		TestZetaDepositRestricted,
	),
	runner.NewE2ETest(
		TestZetaWithdrawName,
		"withdraw ZETA from ZEVM to Ethereum",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000000"},
		},
		TestZetaWithdraw,
	),
	runner.NewE2ETest(
		TestZetaWithdrawBTCRevertName,
		"sending ZETA from ZEVM to Bitcoin with a message that should revert cctxs",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "1000000000000000000"},
		},
		TestZetaWithdrawBTCRevert,
	),
	/*
	 Message passing tests
	*/
	runner.NewE2ETest(
		TestMessagePassingExternalChainsName,
		"evm->evm message passing (sending ZETA only)",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000000"},
		},
		TestMessagePassingExternalChains,
	),
	runner.NewE2ETest(
		TestMessagePassingRevertFailExternalChainsName,
		"message passing with failing revert between external EVM chains",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000000"},
		},
		TestMessagePassingRevertFailExternalChains,
	),
	runner.NewE2ETest(
		TestMessagePassingRevertSuccessExternalChainsName,
		"message passing with successful revert between external EVM chains",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000000"},
		},
		TestMessagePassingRevertSuccessExternalChains,
	),
	runner.NewE2ETest(
		TestMessagePassingEVMtoZEVMName,
		"evm -> zevm message passing contract call ",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000009"},
		},
		TestMessagePassingEVMtoZEVM,
	),
	runner.NewE2ETest(
		TestMessagePassingZEVMToEVMName,
		"zevm -> evm message passing contract call",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000007"},
		},
		TestMessagePassingZEVMtoEVM,
	),
	runner.NewE2ETest(
		TestMessagePassingZEVMtoEVMRevertName,
		"zevm -> evm message passing contract call reverts",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000006"},
		},
		TestMessagePassingZEVMtoEVMRevert,
	),
	runner.NewE2ETest(
		TestMessagePassingEVMtoZEVMRevertName,
		"evm -> zevm message passing and revert back to evm",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000008"},
		},
		TestMessagePassingEVMtoZEVMRevert,
	),
	runner.NewE2ETest(
		TestMessagePassingZEVMtoEVMRevertFailName,
		"zevm -> evm message passing contract with failing revert",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000008"},
		},
		TestMessagePassingZEVMtoEVMRevertFail,
	),
	runner.NewE2ETest(
		TestMessagePassingEVMtoZEVMRevertFailName,
		"evm -> zevm message passing contract with failing revert",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000008"},
		},
		TestMessagePassingEVMtoZEVMRevertFail,
	),

	/*
	 EVM gas tests
	*/
	runner.NewE2ETest(
		TestEtherDepositName,
		"deposit Ether into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestEtherDeposit,
	),
	runner.NewE2ETest(
		TestEtherWithdrawName,
		"withdraw Ether from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestEtherWithdraw,
	),
	runner.NewE2ETest(
		TestEtherWithdrawRestrictedName,
		"withdraw Ether from ZEVM to restricted address",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestEtherWithdrawRestricted,
	),
	runner.NewE2ETest(
		TestEtherDepositAndCallRefundName,
		"deposit Ether into ZEVM and call a contract that reverts; should refund",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000000"},
		},
		TestEtherDepositAndCallRefund,
	),
	runner.NewE2ETest(
		TestEtherDepositAndCallName,
		"deposit ZRC20 into ZEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "1000000000000000000"},
		},
		TestEtherDepositAndCall,
	),
	/*
	 EVM erc20 tests
	*/
	runner.NewE2ETest(
		TestERC20WithdrawName,
		"withdraw ERC20 from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestERC20Withdraw,
	),
	runner.NewE2ETest(
		TestERC20DepositName,
		"deposit ERC20 into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000"},
		},
		TestERC20Deposit,
	),
	runner.NewE2ETest(
		TestMultipleERC20DepositName,
		"deposit ERC20 into ZEVM in multiple deposits",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000000000"},
			{Description: "count", DefaultValue: "3"},
		},
		TestMultipleERC20Deposit,
	),
	runner.NewE2ETest(
		TestMultipleERC20WithdrawsName,
		"withdraw ERC20 from ZEVM in multiple withdrawals",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100"},
			{Description: "count", DefaultValue: "3"},
		},
		TestMultipleERC20Withdraws,
	),
	runner.NewE2ETest(
		TestERC20DepositAndCallRefundName,
		"deposit a non-gas ZRC20 into ZEVM and call a contract that reverts",
		[]runner.ArgDefinition{},
		TestERC20DepositAndCallRefund,
	),
	runner.NewE2ETest(
		TestERC20DepositRestrictedName,
		"deposit ERC20 into ZEVM restricted address",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000"},
		},
		TestERC20DepositRestricted,
	),
	/*
	 Bitcoin tests
	*/
	runner.NewE2ETest(
		TestBitcoinDepositName,
		"deposit Bitcoin into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinDeposit,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawSegWitName,
		"withdraw BTC from ZEVM to a SegWit address",
		[]runner.ArgDefinition{
			{Description: "receiver address", DefaultValue: ""},
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawSegWit,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawTaprootName,
		"withdraw BTC from ZEVM to a Taproot address",
		[]runner.ArgDefinition{
			{Description: "receiver address", DefaultValue: ""},
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawTaproot,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawLegacyName,
		"withdraw BTC from ZEVM to a legacy address",
		[]runner.ArgDefinition{
			{Description: "receiver address", DefaultValue: ""},
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawLegacy,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawP2WSHName,
		"withdraw BTC from ZEVM to a P2WSH address",
		[]runner.ArgDefinition{
			{Description: "receiver address", DefaultValue: ""},
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawP2WSH,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawP2SHName,
		"withdraw BTC from ZEVM to a P2SH address",
		[]runner.ArgDefinition{
			{Description: "receiver address", DefaultValue: ""},
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawP2SH,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawInvalidAddressName,
		"withdraw BTC from ZEVM to an unsupported btc address",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.00001"},
		},
		TestBitcoinWithdrawToInvalidAddress,
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawRestrictedName,
		"withdraw Bitcoin from ZEVM to restricted address",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawRestricted,
	),
	/*
	 Application tests
	*/
	runner.NewE2ETest(
		TestZRC20SwapName,
		"swap ZRC20 ERC20 for ZRC20 ETH",
		[]runner.ArgDefinition{},
		TestZRC20Swap,
	),
	runner.NewE2ETest(
		TestCrosschainSwapName,
		"testing Bitcoin ERC20 cross-chain swap",
		[]runner.ArgDefinition{},
		TestCrosschainSwap,
	),
	/*
	 Miscellaneous tests
	*/
	runner.NewE2ETest(
		TestContextUpgradeName,
		"tests sending ETH on ZEVM and check context data using ContextApp",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "1000000000000000"},
		},
		TestContextUpgrade,
	),
	runner.NewE2ETest(
		TestMyTestName,
		"performing custom test",
		[]runner.ArgDefinition{},
		TestMyTest,
	),
	runner.NewE2ETest(
		TestDonationEtherName,
		"donate Ether to the TSS",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000000000000000"},
		},
		TestDonationEther,
	),
	/*
	 Stress tests
	*/
	runner.NewE2ETest(
		TestStressEtherWithdrawName,
		"stress test Ether withdrawal",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
			{Description: "count", DefaultValue: "100"},
		},
		TestStressEtherWithdraw,
	),
	runner.NewE2ETest(
		TestStressBTCWithdrawName,
		"stress test BTC withdrawal",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.01"},
			{Description: "count", DefaultValue: "100"},
		},
		TestStressBTCWithdraw,
	),
	runner.NewE2ETest(
		TestStressEtherDepositName,
		"stress test Ether deposit",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
			{Description: "count", DefaultValue: "100"},
		},
		TestStressEtherDeposit,
	),
	runner.NewE2ETest(
		TestStressBTCDepositName,
		"stress test BTC deposit",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.001"},
			{Description: "count", DefaultValue: "100"},
		},
		TestStressBTCDeposit,
	),
	/*
	 Admin tests
	*/
	runner.NewE2ETest(
		TestDepositEtherLiquidityCapName,
		"deposit Ethers into ZEVM with a liquidity cap",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000000000000"},
		},
		TestDepositEtherLiquidityCap,
	),
	runner.NewE2ETest(
		TestMigrateChainSupportName,
		"migrate the evm chain from goerli to sepolia",
		[]runner.ArgDefinition{},
		TestMigrateChainSupport,
	),
	runner.NewE2ETest(
		TestPauseZRC20Name,
		"pausing ZRC20 on ZetaChain",
		[]runner.ArgDefinition{},
		TestPauseZRC20,
	),
	runner.NewE2ETest(
		TestUpdateBytecodeZRC20Name,
		"update ZRC20 bytecode swap",
		[]runner.ArgDefinition{},
		TestUpdateBytecodeZRC20,
	),
	runner.NewE2ETest(
		TestUpdateBytecodeConnectorName,
		"update zevm connector bytecode",
		[]runner.ArgDefinition{},
		TestUpdateBytecodeConnector,
	),
	runner.NewE2ETest(
		TestRateLimiterName,
		"test sending cctxs with rate limiter enabled and show logs when processing cctxs",
		[]runner.ArgDefinition{},
		TestRateLimiter,
	),
}
