package e2etests

import (
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/testutil/sample"
)

// List of all e2e test names to be used in zetae2e
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
	TestERC20DepositRestrictedName    = "erc20_deposit_restricted" // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestERC20DepositAndCallRefundName = "erc20_deposit_and_call_refund"

	/*
	 * Solana tests
	 */
	TestSolanaDepositName              = "solana_deposit"
	TestSolanaWithdrawName             = "solana_withdraw"
	TestSolanaDepositAndCallName       = "solana_deposit_and_call"
	TestSolanaDepositAndCallRefundName = "solana_deposit_and_call_refund"
	TestSolanaDepositRestrictedName    = "solana_deposit_restricted"
	TestSolanaWithdrawRestrictedName   = "solana_withdraw_restricted"

	/**
	 * TON tests
	 */
	TestTONDepositName = "ton_deposit"

	/*
	 Bitcoin tests
	 Test transfer of Bitcoin asset across chains
	*/
	TestBitcoinDepositName                = "bitcoin_deposit"
	TestBitcoinDepositRefundName          = "bitcoin_deposit_refund"
	TestBitcoinWithdrawSegWitName         = "bitcoin_withdraw_segwit"
	TestBitcoinWithdrawTaprootName        = "bitcoin_withdraw_taproot"
	TestBitcoinWithdrawMultipleName       = "bitcoin_withdraw_multiple"
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
	TestWhitelistERC20Name            = "whitelist_erc20"
	TestDepositEtherLiquidityCapName  = "deposit_eth_liquidity_cap"
	TestMigrateChainSupportName       = "migrate_chain_support"
	TestPauseZRC20Name                = "pause_zrc20"
	TestUpdateBytecodeZRC20Name       = "update_bytecode_zrc20"
	TestUpdateBytecodeConnectorName   = "update_bytecode_connector"
	TestRateLimiterName               = "rate_limiter"
	TestCriticalAdminTransactionsName = "critical_admin_transactions"
	TestPauseERC20CustodyName         = "pause_erc20_custody"
	TestMigrateERC20CustodyFundsName  = "migrate_erc20_custody_funds"
	TestMigrateTSSName                = "migrate_TSS"

	/*
	 V2 smart contract tests
	*/
	TestV2ETHDepositName                         = "v2_eth_deposit"
	TestV2ETHDepositAndCallName                  = "v2_eth_deposit_and_call"
	TestV2ETHDepositAndCallRevertName            = "v2_eth_deposit_and_call_revert"
	TestV2ETHDepositAndCallRevertWithCallName    = "v2_eth_deposit_and_call_revert_with_call"
	TestV2ETHWithdrawName                        = "v2_eth_withdraw"
	TestV2ETHWithdrawAndCallName                 = "v2_eth_withdraw_and_call"
	TestV2ETHWithdrawAndCallRevertName           = "v2_eth_withdraw_and_call_revert"
	TestV2ETHWithdrawAndCallRevertWithCallName   = "v2_eth_withdraw_and_call_revert_with_call"
	TestV2ERC20DepositName                       = "v2_erc20_deposit"
	TestV2ERC20DepositAndCallName                = "v2_erc20_deposit_and_call"
	TestV2ERC20DepositAndCallRevertName          = "v2_erc20_deposit_and_call_revert"
	TestV2ERC20DepositAndCallRevertWithCallName  = "v2_erc20_deposit_and_call_revert_with_call"
	TestV2ERC20WithdrawName                      = "v2_erc20_withdraw"
	TestV2ERC20WithdrawAndCallName               = "v2_erc20_withdraw_and_call"
	TestV2ERC20WithdrawAndCallRevertName         = "v2_erc20_withdraw_and_call_revert"
	TestV2ERC20WithdrawAndCallRevertWithCallName = "v2_erc20_withdraw_and_call_revert_with_call"
	TestV2ZEVMToEVMCallName                      = "v2_zevm_to_evm_call"
	TestV2EVMToZEVMCallName                      = "v2_evm_to_zevm_call"

	/*
	 Operational tests
	 Not used to test functionalities but do various interactions with the netwoks
	*/
	TestDeploy                         = "deploy"
	TestOperationAddLiquidityETHName   = "add_liquidity_eth"
	TestOperationAddLiquidityERC20Name = "add_liquidity_erc20"

	/*
	 Stateful precompiled contracts tests
	*/
	TestPrecompilesPrototypeName                = "precompile_contracts_prototype"
	TestPrecompilesPrototypeThroughContractName = "precompile_contracts_prototype_through_contract"
	TestPrecompilesStakingName                  = "precompile_contracts_staking"
	TestPrecompilesStakingThroughContractName   = "precompile_contracts_staking_through_contract"
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
		TestERC20DepositRestrictedName,
		"deposit ERC20 into ZEVM restricted address",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000"},
		},
		TestERC20DepositRestricted,
	),
	runner.NewE2ETest(
		TestERC20DepositAndCallRefundName,
		"deposit a non-gas ZRC20 into ZEVM and call a contract that reverts",
		[]runner.ArgDefinition{},
		TestERC20DepositAndCallRefund,
	),
	/*
	 Solana tests
	*/
	runner.NewE2ETest(
		TestSolanaDepositName,
		"deposit SOL into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "12000000"},
		},
		TestSolanaDeposit,
	),
	runner.NewE2ETest(
		TestSolanaWithdrawName,
		"withdraw SOL from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdraw,
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallName,
		"deposit SOL into ZEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1200000"},
		},
		TestSolanaDepositAndCall,
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallRefundName,
		"deposit SOL into ZEVM and call a contract that reverts; should refund",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1200000"},
		},
		TestSolanaDepositAndCallRefund,
	),
	runner.NewE2ETest(
		TestSolanaDepositRestrictedName,
		"deposit SOL into ZEVM restricted address",
		[]runner.ArgDefinition{
			{Description: "receiver", DefaultValue: sample.RestrictedEVMAddressTest},
			{Description: "amount in lamport", DefaultValue: "1200000"},
		},
		TestSolanaDepositRestricted,
	),
	runner.NewE2ETest(
		TestSolanaWithdrawRestrictedName,
		"withdraw SOL from ZEVM to restricted address",
		[]runner.ArgDefinition{
			{Description: "receiver", DefaultValue: sample.RestrictedSolAddressTest},
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdrawRestricted,
	),
	/*
	 TON tests
	*/
	runner.NewE2ETest(
		TestTONDepositName,
		"deposit TON into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "900000000"}, // 0.9 TON
		},
		TestTONDeposit,
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
		TestBitcoinDepositRefundName,
		"deposit Bitcoin into ZEVM; expect refund", []runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.1"},
		},
		TestBitcoinDepositRefund,
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
		TestBitcoinWithdrawMultipleName,
		"withdraw BTC from ZEVM multiple times",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "0.01"},
			{Description: "times", DefaultValue: "2"},
		},
		WithdrawBitcoinMultipleTimes,
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
		TestWhitelistERC20Name,
		"whitelist a new ERC20 token",
		[]runner.ArgDefinition{},
		TestWhitelistERC20,
	),
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
	runner.NewE2ETest(
		TestCriticalAdminTransactionsName,
		"test critical admin transactions",
		[]runner.ArgDefinition{},
		TestCriticalAdminTransactions,
	),
	runner.NewE2ETest(
		TestMigrateTSSName,
		"migrate TSS funds",
		[]runner.ArgDefinition{},
		TestMigrateTSS,
	),
	runner.NewE2ETest(
		TestPauseERC20CustodyName,
		"pausing ERC20 custody on ZetaChain",
		[]runner.ArgDefinition{},
		TestPauseERC20Custody,
	),
	runner.NewE2ETest(
		TestMigrateERC20CustodyFundsName,
		"migrate ERC20 custody funds",
		[]runner.ArgDefinition{},
		TestMigrateERC20CustodyFunds,
	),
	/*
	 V2 smart contract tests
	*/
	runner.NewE2ETest(
		TestV2ETHDepositName,
		"deposit Ether into ZEVM using V2 contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000000000000000000"},
		},
		TestV2ETHDeposit,
	),
	runner.NewE2ETest(
		TestV2ETHDepositAndCallName,
		"deposit Ether into ZEVM and call a contract using V2 contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestV2ETHDepositAndCall,
	),
	runner.NewE2ETest(
		TestV2ETHDepositAndCallRevertName,
		"deposit Ether into ZEVM and call a contract using V2 contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestV2ETHDepositAndCallRevert,
	),
	runner.NewE2ETest(
		TestV2ETHDepositAndCallRevertWithCallName,
		"deposit Ether into ZEVM and call a contract using V2 contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestV2ETHDepositAndCallRevertWithCall,
	),
	runner.NewE2ETest(
		TestV2ETHWithdrawName,
		"withdraw Ether from ZEVM using V2 contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestV2ETHWithdraw,
	),
	runner.NewE2ETest(
		TestV2ETHWithdrawAndCallName,
		"withdraw Ether from ZEVM and call a contract using V2 contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestV2ETHWithdrawAndCall,
	),
	runner.NewE2ETest(
		TestV2ETHWithdrawAndCallRevertName,
		"withdraw Ether from ZEVM and call a contract using V2 contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestV2ETHWithdrawAndCallRevert,
	),
	runner.NewE2ETest(
		TestV2ETHWithdrawAndCallRevertWithCallName,
		"withdraw Ether from ZEVM and call a contract using V2 contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestV2ETHWithdrawAndCallRevertWithCall,
	),
	runner.NewE2ETest(
		TestV2ERC20DepositName,
		"deposit ERC20 into ZEVM using V2 contract",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000000000000000000"},
		},
		TestV2ERC20Deposit,
	),
	runner.NewE2ETest(
		TestV2ERC20DepositAndCallName,
		"deposit ERC20 into ZEVM and call a contract using V2 contract",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000"},
		},
		TestV2ERC20DepositAndCall,
	),
	runner.NewE2ETest(
		TestV2ERC20DepositAndCallRevertName,
		"deposit ERC20 into ZEVM and call a contract using V2 contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "10000000000000000000"},
		},
		TestV2ERC20DepositAndCallRevert,
	),
	runner.NewE2ETest(
		TestV2ERC20DepositAndCallRevertWithCallName,
		"deposit ERC20 into ZEVM and call a contract using V2 contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "10000000000000000000"},
		},
		TestV2ERC20DepositAndCallRevertWithCall,
	),
	runner.NewE2ETest(
		TestV2ERC20WithdrawName,
		"withdraw ERC20 from ZEVM using V2 contract",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestV2ERC20Withdraw,
	),
	runner.NewE2ETest(
		TestV2ERC20WithdrawAndCallName,
		"withdraw ERC20 from ZEVM and call a contract using V2 contract",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestV2ERC20WithdrawAndCall,
	),
	runner.NewE2ETest(
		TestV2ERC20WithdrawAndCallRevertName,
		"withdraw ERC20 from ZEVM and call a contract using V2 contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestV2ERC20WithdrawAndCallRevert,
	),
	runner.NewE2ETest(
		TestV2ERC20WithdrawAndCallRevertWithCallName,
		"withdraw ERC20 from ZEVM and call a contract using V2 contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestV2ERC20WithdrawAndCallRevertWithCall,
	),
	runner.NewE2ETest(
		TestV2ZEVMToEVMCallName,
		"zevm -> evm call using V2 contract",
		[]runner.ArgDefinition{},
		TestV2ZEVMToEVMCall,
	),
	runner.NewE2ETest(
		TestV2EVMToZEVMCallName,
		"evm -> zevm call using V2 contract",
		[]runner.ArgDefinition{},
		TestV2EVMToZEVMCall,
	),
	/*
	 Special tests
	*/
	runner.NewE2ETest(
		TestDeploy,
		"deploy a contract",
		[]runner.ArgDefinition{
			{Description: "contract name", DefaultValue: ""},
		},
		TestDeployContract,
	),
	runner.NewE2ETest(
		TestOperationAddLiquidityETHName,
		"add liquidity to the ZETA/ETH pool",
		[]runner.ArgDefinition{
			{Description: "amountZETA", DefaultValue: "50000000000000000000"},
			{Description: "amountETH", DefaultValue: "50000000000000000000"},
		},
		TestOperationAddLiquidityETH,
	),
	runner.NewE2ETest(
		TestOperationAddLiquidityERC20Name,
		"add liquidity to the ZETA/ERC20 pool",
		[]runner.ArgDefinition{
			{Description: "amountZETA", DefaultValue: "50000000000000000000"},
			{Description: "amountERC20", DefaultValue: "50000000000000000000"},
		},
		TestOperationAddLiquidityERC20,
	),
	/*
	 Stateful precompiled contracts tests
	*/
	runner.NewE2ETest(
		TestPrecompilesPrototypeName,
		"test stateful precompiled contracts prototype",
		[]runner.ArgDefinition{},
		TestPrecompilesPrototype,
	),
	runner.NewE2ETest(
		TestPrecompilesPrototypeThroughContractName,
		"test stateful precompiled contracts prototype through contract",
		[]runner.ArgDefinition{},
		TestPrecompilesPrototypeThroughContract,
	),
	runner.NewE2ETest(
		TestPrecompilesStakingName,
		"test stateful precompiled contracts staking",
		[]runner.ArgDefinition{},
		TestPrecompilesStaking,
	),
	runner.NewE2ETest(
		TestPrecompilesStakingThroughContractName,
		"test stateful precompiled contracts staking through contract",
		[]runner.ArgDefinition{},
		TestPrecompilesStakingThroughContract,
	),
}
