package e2etests

import (
	"github.com/zeta-chain/node/e2e/e2etests/legacy"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/testutil/sample"
)

// List of all e2e test names to be used in zetae2e
const (
	/*
	  EVM chain tests
	*/
	TestETHDepositName                      = "eth_deposit"
	TestETHDepositAndCallBigPayloadName     = "eth_deposit_and_call_big_payload"
	TestETHDepositAndCallName               = "eth_deposit_and_call"
	TestETHDepositFastConfirmationName      = "eth_deposit_fast_confirmation"
	TestETHDepositAndCallNoMessageName      = "eth_deposit_and_call_no_message"
	TestETHDepositAndCallRevertName         = "eth_deposit_and_call_revert"
	TestETHDepositAndCallRevertWithCallName = "eth_deposit_and_call_revert_with_call"
	TestETHDepositRevertAndAbortName        = "eth_deposit_revert_and_abort"

	TestETHWithdrawName                          = "eth_withdraw"
	TestETHWithdrawCustomGasLimitName            = "eth_withdraw_custom_gas_limit"
	TestETHWithdrawAndArbitraryCallName          = "eth_withdraw_and_arbitrary_call"
	TestETHWithdrawAndCallName                   = "eth_withdraw_and_call"
	TestETHWithdrawAndCallBigPayloadName         = "eth_withdraw_and_call_big_payload"
	TestETHWithdrawAndCallNoMessageName          = "eth_withdraw_and_call_no_message"
	TestETHWithdrawAndCallThroughContractName    = "eth_withdraw_and_call_through_contract"
	TestETHWithdrawAndCallRevertName             = "eth_withdraw_and_call_revert"
	TestETHWithdrawAndCallRevertWithCallName     = "eth_withdraw_and_call_revert_with_call"
	TestETHWithdrawRevertAndAbortName            = "eth_withdraw_revert_and_abort"
	TestETHWithdrawAndCallRevertWithWithdrawName = "eth_withdraw_and_call_revert_with_withdraw"
	TestDepositAndCallOutOfGasName               = "deposit_and_call_out_of_gas"
	TestDepositAndCallHighGasUsageName           = "deposit_and_call_high_gas_usage"

	TestERC20DepositName                      = "erc20_deposit"
	TestERC20DepositAndCallName               = "erc20_deposit_and_call"
	TestERC20DepositAndCallNoMessageName      = "erc20_deposit_and_call_no_message"
	TestERC20DepositAndCallRevertName         = "erc20_deposit_and_call_revert"
	TestERC20DepositAndCallRevertWithCallName = "erc20_deposit_and_call_revert_with_call"
	TestERC20DepositRevertAndAbortName        = "erc20_deposit_revert_and_abort"

	TestERC20WithdrawName                      = "erc20_withdraw"
	TestERC20WithdrawAndArbitraryCallName      = "erc20_withdraw_and_arbitrary_call"
	TestERC20WithdrawAndCallName               = "erc20_withdraw_and_call"
	TestERC20WithdrawAndCallNoMessageName      = "erc20_withdraw_and_call_no_message"
	TestERC20WithdrawAndCallRevertName         = "erc20_withdraw_and_call_revert"
	TestERC20WithdrawAndCallRevertWithCallName = "erc20_withdraw_and_call_revert_with_call"
	TestERC20WithdrawRevertAndAbortName        = "erc20_withdraw_revert_and_abort"

	TestZEVMToEVMArbitraryCallName       = "zevm_to_evm_arbitrary_call"
	TestZEVMToEVMCallName                = "zevm_to_evm_call"
	TestZEVMToEVMCallRevertName          = "zevm_to_evm_call_revert"
	TestZEVMToEVMCallRevertAndAbortName  = "zevm_to_evm_call_revert_and_abort"
	TestZEVMToEVMCallThroughContractName = "zevm_to_evm_call_through_contract"
	TestEVMToZEVMCallName                = "evm_to_zevm_call"
	TestEVMToZEVMCallAbortName           = "evm_to_zevm_abort_call"

	TestDepositAndCallSwapName      = "deposit_and_call_swap"
	TestEtherWithdrawRestrictedName = "eth_withdraw_restricted"
	TestERC20DepositRestrictedName  = "erc20_deposit_restricted" // #nosec G101: Potential hardcoded credentials (gosec), not a credential

	/*
	 * Solana tests
	 */
	TestSolanaDepositName                                 = "solana_deposit"
	TestSolanaDepositThroughProgramName                   = "solana_deposit_through_program"
	TestSolanaWithdrawName                                = "solana_withdraw"
	TestSolanaWithdrawRevertExecutableReceiverName        = "solana_withdraw_revert_executable_receiver"
	TestSolanaWithdrawAndCallName                         = "solana_withdraw_and_call"
	TestSolanaWithdrawAndCallInvalidTxSizeName            = "solana_withdraw_and_call_invalid_tx_size"
	TestSolanaWithdrawAndCallInvalidMsgEncodingName       = "solana_withdraw_and_call_invalid_msg_encoding"
	TestZEVMToSolanaCallName                              = "zevm_to_solana_call"
	TestSolanaWithdrawAndCallRevertWithCallName           = "solana_withdraw_and_call_revert_with_call"
	TestSolanaDepositAndCallName                          = "solana_deposit_and_call"
	TestSolanaDepositAndCallRevertName                    = "solana_deposit_and_call_revert"
	TestSolanaDepositAndCallRevertWithCallName            = "solana_deposit_and_call_revert_with_call"
	TestSolanaDepositAndCallRevertWithCallThatRevertsName = "solana_deposit_and_call_revert_with_call_that_reverts"
	TestSolanaDepositAndCallRevertWithDustName            = "solana_deposit_and_call_revert_with_dust"
	TestSolanaToZEVMCallName                              = "solana_to_zevm_call"
	TestSolanaToZEVMCallAbortName                         = "solana_to_zevm_call_abort"
	TestSolanaDepositRestrictedName                       = "solana_deposit_restricted"
	TestSolanaWithdrawRestrictedName                      = "solana_withdraw_restricted"
	TestSPLDepositName                                    = "spl_deposit"
	TestSPLDepositAndCallName                             = "spl_deposit_and_call"
	TestSPLDepositAndCallRevertName                       = "spl_deposit_and_call_revert"
	TestSPLDepositAndCallRevertWithCallName               = "spl_deposit_and_call_revert_with_call"
	TestSPLDepositAndCallRevertWithCallThatRevertsName    = "spl_deposit_and_call_revert_with_call_that_reverts"
	TestSPLWithdrawName                                   = "spl_withdraw"
	TestSPLWithdrawAndCallName                            = "spl_withdraw_and_call"
	TestSPLWithdrawAndCallRevertName                      = "spl_withdraw_and_call_revert"
	TestSPLWithdrawAndCreateReceiverAtaName               = "spl_withdraw_and_create_receiver_ata"

	/**
	 * TON tests
	 */
	TestTONDepositName              = "ton_deposit"
	TestTONDepositAndCallName       = "ton_deposit_and_call"
	TestTONDepositAndCallRefundName = "ton_deposit_refund"
	TestTONDepositRestrictedName    = "ton_deposit_restricted"
	TestTONCallName                 = "ton_to_zevm_call"
	TestTONWithdrawName             = "ton_withdraw"
	TestTONWithdrawRestrictedName   = "ton_withdraw_restricted"
	TestTONWithdrawMasterchainName  = "ton_withdraw_masterchain"
	TestTONWithdrawConcurrentName   = "ton_withdraw_concurrent"

	/*
	 Sui tests
	*/
	TestSuiDepositName                            = "sui_deposit"
	TestSuiDepositAndCallName                     = "sui_deposit_and_call"
	TestSuiDepositAndCallRevertName               = "sui_deposit_and_call_revert"
	TestSuiTokenDepositName                       = "sui_token_deposit"                 // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiTokenDepositAndCallName                = "sui_token_deposit_and_call"        // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiTokenDepositAndCallRevertName          = "sui_token_deposit_and_call_revert" // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiWithdrawName                           = "sui_withdraw"
	TestSuiTokenWithdrawName                      = "sui_token_withdraw"                           // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiTokenWithdrawAndCallName               = "sui_token_withdraw_and_call"                  // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiTokenWithdrawAndCallRevertWithCallName = "sui_token_withdraw_and_call_revert_with_call" // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiWithdrawAndCallName                    = "sui_withdraw_and_call"
	TestSuiWithdrawRevertWithCallName             = "sui_withdraw_revert_with_call"          // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiWithdrawAndCallInvalidPayloadName      = "sui_withdraw_and_call_invalid_payload"  // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiWithdrawAndCallRevertWithCallName      = "sui_withdraw_and_call_revert_with_call" // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestSuiDepositRestrictedName                  = "sui_deposit_restricted"
	TestSuiWithdrawRestrictedName                 = "sui_withdraw_restricted"
	TestSuiWithdrawInvalidReceiverName            = "sui_withdraw_invalid_receiver"

	/*
	 Bitcoin tests
	 Test transfer of Bitcoin asset across chains
	*/
	TestBitcoinDepositName                                 = "bitcoin_deposit"
	TestBitcoinDepositAndCallName                          = "bitcoin_deposit_and_call"
	TestBitcoinDepositFastConfirmationName                 = "bitcoin_deposit_fast_confirmation"
	TestBitcoinDepositAndCallRevertName                    = "bitcoin_deposit_and_call_revert"
	TestBitcoinDepositAndCallRevertWithDustName            = "bitcoin_deposit_and_call_revert_with_dust"
	TestBitcoinDepositAndWithdrawWithDustName              = "bitcoin_deposit_and_withdraw_with_dust"
	TestBitcoinDonationName                                = "bitcoin_donation"
	TestBitcoinStdMemoDepositName                          = "bitcoin_std_memo_deposit"
	TestBitcoinStdMemoDepositAndCallName                   = "bitcoin_std_memo_deposit_and_call"
	TestBitcoinStdMemoDepositAndCallRevertName             = "bitcoin_std_memo_deposit_and_call_revert"
	TestBitcoinStdMemoDepositAndCallRevertOtherAddressName = "bitcoin_std_memo_deposit_and_call_revert_other_address"
	TestBitcoinStdMemoDepositAndCallRevertAndAbortName     = "bitcoin_std_memo_deposit_and_call_revert_and_abort"
	TestBitcoinStdMemoInscribedDepositAndCallName          = "bitcoin_std_memo_inscribed_deposit_and_call"
	TestBitcoinDepositAndAbortWithLowDepositFeeName        = "bitcoin_deposit_and_abort_with_low_deposit_fee"
	TestBitcoinWithdrawSegWitName                          = "bitcoin_withdraw_segwit"
	TestBitcoinWithdrawTaprootName                         = "bitcoin_withdraw_taproot"
	TestBitcoinWithdrawMultipleName                        = "bitcoin_withdraw_multiple"
	TestBitcoinWithdrawLegacyName                          = "bitcoin_withdraw_legacy"
	TestBitcoinWithdrawP2WSHName                           = "bitcoin_withdraw_p2wsh"
	TestBitcoinWithdrawP2SHName                            = "bitcoin_withdraw_p2sh"
	TestBitcoinWithdrawInvalidAddressName                  = "bitcoin_withdraw_invalid"
	TestBitcoinWithdrawRestrictedName                      = "bitcoin_withdraw_restricted"
	TestBitcoinDepositInvalidMemoRevertName                = "bitcoin_deposit_invalid_memo_revert"
	TestBitcoinWithdrawRBFName                             = "bitcoin_withdraw_rbf"

	/*
	 Application tests
	 Test various smart contract applications across chains
	*/
	TestCrosschainSwapName = "crosschain_swap"

	/*
	 Miscellaneous tests
	 Test various functionalities not related to assets
	*/
	TestDonationEtherName   = "donation_ether"
	TestInboundTrackersName = "inbound_trackers"
	TestPrecompilesName     = "precompiles"
	TestOpcodesName         = "opcodes"

	/*
	 Stress tests
	 Test stressing networks with many cross-chain transactions
	*/
	TestStressEtherWithdrawName  = "stress_eth_withdraw"
	TestStressBTCWithdrawName    = "stress_btc_withdraw"
	TestStressEtherDepositName   = "stress_eth_deposit"
	TestStressBTCDepositName     = "stress_btc_deposit"
	TestStressSolanaDepositName  = "stress_solana_deposit"
	TestStressSPLDepositName     = "stress_spl_deposit"
	TestStressSolanaWithdrawName = "stress_solana_withdraw"
	TestStressSPLWithdrawName    = "stress_spl_withdraw"
	TestStressSuiDepositName     = "stress_sui_deposit"
	TestStressSuiWithdrawName    = "stress_sui_withdraw"

	/*
		Staking tests
	*/

	TestUndelegateToBelowMinimumObserverDelegation = "undelegate_to_below_minimum_observer_delegation"

	/*
	 Admin tests
	 Test admin functionalities
	*/
	TestWhitelistERC20Name               = "whitelist_erc20"
	TestDepositEtherLiquidityCapName     = "deposit_eth_liquidity_cap"
	TestMigrateChainSupportName          = "migrate_chain_support"
	TestPauseZRC20Name                   = "pause_zrc20"
	TestUpdateBytecodeZRC20Name          = "update_bytecode_zrc20"
	TestUpdateBytecodeConnectorName      = "update_bytecode_connector"
	TestRateLimiterName                  = "rate_limiter"
	TestCriticalAdminTransactionsName    = "critical_admin_transactions"
	TestPauseERC20CustodyName            = "pause_erc20_custody"
	TestMigrateERC20CustodyFundsName     = "migrate_erc20_custody_funds"
	TestMigrateTSSName                   = "migrate_tss"
	TestSolanaWhitelistSPLName           = "solana_whitelist_spl"
	TestUpdateZRC20NameName              = "update_zrc20_name"
	TestZetaclientRestartHeightName      = "zetaclient_restart_height"
	TestZetaclientSignerOffsetName       = "zetaclient_signer_offset"
	TestUpdateOperationalChainParamsName = "update_operational_chain_params"
	TestMigrateConnectorFundsName        = "migrate_connector_funds"
	TestBurnFungibleModuleAssetName      = "burn_fungible_module_asset"

	/*
	 Operational tests
	 Not used to test functionalities but do various interactions with the netwoks
	*/
	TestDeploy                                    = "deploy"
	TestOperationAddLiquidityETHName              = "add_liquidity_eth"
	TestOperationAddLiquidityERC20Name            = "add_liquidity_erc20"
	TestOperationAddLiquidityBTCName              = "add_liquidity_btc"
	TestOperationAddLiquiditySOLName              = "add_liquidity_sol"
	TestOperationAddLiquiditySPLName              = "add_liquidity_spl"
	TestOperationAddLiquiditySUIName              = "add_liquidity_sui"
	TestOperationAddLiquiditySuiFungibleTokenName = "add_liquidity_sui_fungible_token" // #nosec G101: Potential hardcoded credentials (gosec), not a credential
	TestOperationAddLiquidityTONName              = "add_liquidity_ton"

	/*
	 Legacy tests (using v1 protocol contracts)
	*/
	TestLegacyMessagePassingExternalChainsName              = "legacy_message_passing_external_chains"
	TestLegacyMessagePassingRevertFailExternalChainsName    = "legacy_message_passing_revert_fail"
	TestLegacyMessagePassingRevertSuccessExternalChainsName = "legacy_message_passing_revert_success"
	TestLegacyMessagePassingEVMtoZEVMName                   = "legacy_message_passing_evm_to_zevm"
	TestLegacyMessagePassingZEVMToEVMName                   = "legacy_message_passing_zevm_to_evm"
	TestLegacyMessagePassingZEVMtoEVMRevertName             = "legacy_message_passing_zevm_to_evm_revert"
	TestLegacyMessagePassingEVMtoZEVMRevertName             = "legacy_message_passing_evm_to_zevm_revert"
	TestLegacyMessagePassingZEVMtoEVMRevertFailName         = "legacy_message_passing_zevm_to_evm_revert_fail"
	TestLegacyMessagePassingEVMtoZEVMRevertFailName         = "legacy_message_passing_evm_to_zevm_revert_fail"
	TestLegacyEtherDepositName                              = "legacy_eth_deposit"
	TestLegacyEtherWithdrawName                             = "legacy_eth_withdraw"
	TestLegacyEtherDepositAndCallRefundName                 = "legacy_eth_deposit_and_call_refund"
	TestLegacyEtherDepositAndCallName                       = "legacy_eth_deposit_and_call"
	TestLegacyERC20WithdrawName                             = "legacy_erc20_withdraw"
	TestLegacyERC20DepositName                              = "legacy_erc20_deposit"
	TestLegacyMultipleERC20DepositName                      = "legacy_erc20_multiple_deposit"
	TestLegacyMultipleERC20WithdrawsName                    = "legacy_erc20_multiple_withdraw"
	TestLegacyERC20DepositAndCallRefundName                 = "legacy_erc20_deposit_and_call_refund"

	/*
	 ZETA tests
	 Test transfer of ZETA asset across chains
	 Note: It is still the only way to transfer ZETA across chains. Work to integrate ZETA transfers as part of the gateway is in progress
	 These tests are marked as legacy because there is no longer active development on ZETA transfers, and we stopped integrating ZETA support on new mainnet chains
	*/
	TestLegacyZetaDepositName             = "legacy_zeta_deposit"
	TestLegacyZetaDepositAndCallAbortName = "legacy_zeta_deposit_and_call_abort"
	TestLegacyZetaDepositNewAddressName   = "legacy_zeta_deposit_new_address"
	TestLegacyZetaDepositRestrictedName   = "legacy_zeta_deposit_restricted"
	TestLegacyZetaWithdrawName            = "legacy_zeta_withdraw"
	TestLegacyZetaWithdrawBTCRevertName   = "legacy_zeta_withdraw_btc_revert" // #nosec G101 - not a hardcoded password

	TestZetaDepositName                       = "zeta_deposit"
	TestZetaDepositAndCallName                = "zeta_deposit_and_call"
	TestZetaDepositAndCallRevertName          = "zeta_deposit_and_call_revert"
	TestZetaDepositRevertAndAbortName         = "zeta_deposit_revert_and_abort"
	TestZetaDepositAndCallRevertWithCallName  = "zeta_deposit_and_call_revert_with_call"
	TestZetaDepositAndCallNoMessageName       = "zeta_deposit_and_call_no_message"
	TestZetaWithdrawName                      = "zeta_withdraw"
	TestZetaWithdrawAndCallName               = "zeta_withdraw_and_call"
	TestZetaWithdrawAndArbitraryCallName      = "zeta_withdraw_and_arbitrary_call"
	TestZetaWithdrawAndCallRevertName         = "zeta_withdraw_and_call_revert"
	TestZetaWithdrawAndCallRevertWithCallName = "zeta_withdraw_and_call_revert_with_call"
	TestZetaWithdrawRevertAndAbortName        = "zeta_withdraw_revert_and_abort"
)

const (
	CountArgDescription = "count"
)

// Here are all the dependencies for the e2e tests, add more dependencies here if needed
var (
	// DepdencyAllBitcoinDeposits is a dependency to wait for all bitcoin deposit tests to complete
	DepdencyAllBitcoinDeposits = runner.NewE2EDependency("all_bitcoin_deposits")
)

// v2ZetaVersion is the minimum version that supports ZETA transfers using the gateway (v2 protocol contracts)
const v2ZetaVersion = "v37.0.0"

// AllE2ETests is an ordered list of all e2e tests
var AllE2ETests = []runner.E2ETest{
	/*
	 EVM chain tests
	*/
	runner.NewE2ETest(
		TestZetaDepositName,
		"deposit ZETA into ZEVM using connector contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "1000000000000000000"},
		},
		TestZetaDeposit,
		runner.WithMinimumVersion(v2ZetaVersion),
	),
	runner.NewE2ETest(
		TestZetaDepositAndCallName,
		"deposit Zeta into ZEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000000000000000000"},
		},
		TestZetaDepositAndCall,
		runner.WithMinimumVersion(v2ZetaVersion),
	),
	runner.NewE2ETest(
		TestZetaDepositAndCallRevertName,
		"deposit Zeta into ZEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "10000000000000000000"},
		},
		TestZetaDepositAndCallRevert,
		runner.WithMinimumVersion(v2ZetaVersion),
	),
	runner.NewE2ETest(
		TestZetaDepositRevertAndAbortName,
		"deposit Zeta into ZEVM, revert, then abort with onAbort because revert fee cannot be paid",
		[]runner.ArgDefinition{},
		TestZetaDepositRevertAndAbort,
		runner.WithMinimumVersion(v2ZetaVersion),
	),
	runner.NewE2ETest(
		TestZetaDepositAndCallRevertWithCallName,
		"deposit Zeta into ZEVM and call a contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "10000000000000000000"},
		},
		TestZetaDepositAndCallRevertWithCall,
		runner.WithMinimumVersion(v2ZetaVersion),
	),
	runner.NewE2ETest(
		TestZetaDepositAndCallNoMessageName,
		"deposit Zeta into ZEVM and call a contract using no message content",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "10000000000000000000"},
		},
		TestZetaDepositAndCallNoMessage,
		runner.WithMinimumVersion(v2ZetaVersion),
	),
	runner.NewE2ETest(
		TestZetaWithdrawName,
		"withdraw Zeta from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestZetaWithdraw,
		runner.WithMinimumVersion(v2ZetaVersion),
	),
	runner.NewE2ETest(
		TestZetaWithdrawAndCallName,
		"withdraw zeta from ZEVM and call a contract on connected eth chain",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
			{Description: "gas limit for withdraw and call", DefaultValue: "350000"},
		},
		TestZetaWithdrawAndCall,
		runner.WithMinimumVersion(v2ZetaVersion),
	),
	runner.NewE2ETest(
		TestZetaWithdrawAndCallRevertName,
		"withdraw Zeta from ZEVM and call a contract on connected eth chain that reverts",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestZetaWithdrawAndCallRevert,
		runner.WithMinimumVersion(v2ZetaVersion),
	),
	runner.NewE2ETest(
		TestZetaWithdrawAndCallRevertWithCallName,
		"withdraw Zeta from ZEVM and call a contract on connected eth chain that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestZetaWithdrawAndCallRevertWithCall,
		runner.WithMinimumVersion(v2ZetaVersion),
	),
	runner.NewE2ETest(
		TestZetaWithdrawRevertAndAbortName,
		"withdraw Zeta from ZEVM, revert, then abort with onAbort",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
			{Description: "gas limit for withdraw and call", DefaultValue: "350000"},
		},
		TestZetaWithdrawRevertAndAbort,
		runner.WithMinimumVersion(v2ZetaVersion),
	),
	runner.NewE2ETest(
		TestZetaWithdrawAndArbitraryCallName,
		"withdraw Zeta from ZEVM and arbitrary call a contract",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestZetaWithdrawAndArbitraryCall,
		runner.WithMinimumVersion(v2ZetaVersion),
	),
	runner.NewE2ETest(
		TestETHDepositName,
		"deposit Ether into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000000000000000000"},
		},
		TestETHDeposit,
	),
	runner.NewE2ETest(
		TestETHDepositAndCallName,
		"deposit Ether into ZEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestETHDepositAndCall,
	),
	runner.NewE2ETest(
		TestETHDepositAndCallBigPayloadName,
		"deposit Ether to ZetaChain call a contract with a big payload",
		[]runner.ArgDefinition{},
		TestETHDepositAndCallBigPayload,
		runner.WithMinimumVersion("v32.0.0"),
	),
	runner.NewE2ETest(
		TestETHDepositFastConfirmationName,
		"deposit Ether into ZEVM using fast confirmation",
		[]runner.ArgDefinition{},
		TestETHDepositFastConfirmation,
		runner.WithMinimumVersion("v37.0.0"),
	),
	runner.NewE2ETest(
		TestETHDepositAndCallNoMessageName,
		"deposit Ether into ZEVM and call a contract using no message content",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestETHDepositAndCallNoMessage,
	),
	runner.NewE2ETest(
		TestETHDepositAndCallRevertName,
		"deposit Ether into ZEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestETHDepositAndCallRevert,
	),
	runner.NewE2ETest(
		TestETHDepositAndCallRevertWithCallName,
		"deposit Ether into ZEVM and call a contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestETHDepositAndCallRevertWithCall,
	),
	runner.NewE2ETest(
		TestETHDepositRevertAndAbortName,
		"deposit Ether into ZEVM, revert, then abort with onAbort",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestETHDepositRevertAndAbort,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestETHWithdrawName,
		"withdraw Ether from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestETHWithdraw,
	),
	runner.NewE2ETest(
		TestETHWithdrawCustomGasLimitName,
		"withdraw Ether from ZEVM using a custom gas limit",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
			{Description: "gas limit for withdraw", DefaultValue: "200000"},
		},
		TestETHWithdrawCustomGasLimit,
	),
	runner.NewE2ETest(
		TestETHWithdrawAndArbitraryCallName,
		"withdraw Ether from ZEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestETHWithdrawAndArbitraryCall,
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallName,
		"withdraw Ether from ZEVM call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
			{Description: "gas limit for withdraw", DefaultValue: "350000"},
		},
		TestETHWithdrawAndCall,
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallBigPayloadName,
		"withdraw Ether from ZEVM call a contract with a big payload",
		[]runner.ArgDefinition{},
		TestETHWithdrawAndCallBigPayload,
		runner.WithMinimumVersion("v32.0.0"),
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallNoMessageName,
		"withdraw Ether from ZEVM call a contract with no message content",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
			{Description: "gas limit for withdraw", DefaultValue: "350000"},
		},
		TestETHWithdrawAndCallNoMessage,
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallThroughContractName,
		"withdraw Ether from ZEVM call a contract through intermediary contract",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestETHWithdrawAndCallThroughContract,
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallRevertName,
		"withdraw Ether from ZEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestETHWithdrawAndCallRevert,
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallRevertWithCallName,
		"withdraw Ether from ZEVM and call a contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestETHWithdrawAndCallRevertWithCall,
	),
	runner.NewE2ETest(
		TestETHWithdrawRevertAndAbortName,
		"withdraw Ether from ZEVM, revert, then abort with onAbort, check onAbort can created cctx",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "1000000000000000000"},
			{Description: "gas limit for withdraw", DefaultValue: "350000"},
		},
		TestETHWithdrawRevertAndAbort,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallRevertWithWithdrawName,
		"withdraw Ether from ZEVM and call a contract that reverts with a onRevert call that triggers a withdraw",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestETHWithdrawAndCallRevertWithWithdraw,
		runner.WithMinimumVersion("v26.0.0"),
	),
	runner.NewE2ETest(
		TestDepositAndCallHighGasUsageName,
		"deposit Ether into ZEVM and call a contract that consumes high gas",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestDepositAndCallHighGasUsage,
	),
	runner.NewE2ETest(
		TestDepositAndCallOutOfGasName,
		"deposit Ether into ZEVM and call a contract that runs out of gas",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		TestDepositAndCallOutOfGas,
	),
	runner.NewE2ETest(
		TestERC20DepositName,
		"deposit ERC20 into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000000000000000000"},
		},
		TestERC20Deposit,
	),
	runner.NewE2ETest(
		TestERC20DepositAndCallName,
		"deposit ERC20 into ZEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000"},
		},
		TestERC20DepositAndCall,
	),
	runner.NewE2ETest(
		TestERC20DepositAndCallNoMessageName,
		"deposit ERC20 into ZEVM and call a contract with no message content",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000"},
		},
		TestERC20DepositAndCallNoMessage,
	),
	runner.NewE2ETest(
		TestERC20DepositAndCallRevertName,
		"deposit ERC20 into ZEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "10000000000000000000"},
		},
		TestERC20DepositAndCallRevert,
	),
	runner.NewE2ETest(
		TestERC20DepositAndCallRevertWithCallName,
		"deposit ERC20 into ZEVM and call a contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "10000000000000000000"},
		},
		TestERC20DepositAndCallRevertWithCall,
	),
	runner.NewE2ETest(
		TestERC20DepositRevertAndAbortName,
		"deposit ERC20 into ZEVM, revert, then abort with onAbort because revert fee cannot be paid",
		[]runner.ArgDefinition{},
		TestERC20DepositRevertAndAbort,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestERC20WithdrawName,
		"withdraw ERC20 from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestERC20Withdraw,
	),
	runner.NewE2ETest(
		TestERC20WithdrawAndArbitraryCallName,
		"withdraw ERC20 from ZEVM and arbitrary call a contract",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestERC20WithdrawAndArbitraryCall,
	),
	runner.NewE2ETest(
		TestERC20WithdrawAndCallName,
		"withdraw ERC20 from ZEVM and authenticated call a contract",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
			{Description: "gas limit for withdraw and call", DefaultValue: "350000"},
		},
		TestERC20WithdrawAndCall,
	),
	runner.NewE2ETest(
		TestERC20WithdrawAndCallNoMessageName,
		"withdraw ERC20 from ZEVM and authenticated call a contract with no message",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
			{Description: "gas limit for withdraw and call", DefaultValue: "350000"},
		},
		TestERC20WithdrawAndCallNoMessage,
	),
	runner.NewE2ETest(
		TestERC20WithdrawAndCallRevertName,
		"withdraw ERC20 from ZEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestERC20WithdrawAndCallRevert,
	),
	runner.NewE2ETest(
		TestERC20WithdrawAndCallRevertWithCallName,
		"withdraw ERC20 from ZEVM and call a contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		TestERC20WithdrawAndCallRevertWithCall,
	),
	runner.NewE2ETest(
		TestERC20WithdrawRevertAndAbortName,
		"withdraw ERC20 from ZEVM, revert, then abort with onAbort",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
			{Description: "gas limit for withdraw and call", DefaultValue: "250000"},
		},
		TestERC20WithdrawRevertAndAbort,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestZEVMToEVMArbitraryCallName,
		"zevm -> evm call",
		[]runner.ArgDefinition{},
		TestZEVMToEVMArbitraryCall,
	),
	runner.NewE2ETest(
		TestZEVMToEVMCallName,
		"zevm -> evm call",
		[]runner.ArgDefinition{
			{Description: "gas limit for call", DefaultValue: "250000"},
		},
		TestZEVMToEVMCall,
	),
	runner.NewE2ETest(
		TestZEVMToEVMCallRevertName,
		"zevm -> evm call that reverts and call onRevert",
		[]runner.ArgDefinition{
			{Description: "gas limit for call", DefaultValue: "250000"},
		},
		TestZEVMToEVMCallRevert,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestZEVMToEVMCallRevertAndAbortName,
		"zevm -> evm call that reverts and abort with onAbort",
		[]runner.ArgDefinition{
			{Description: "gas limit for call", DefaultValue: "250000"},
		},
		TestZEVMToEVMCallRevertAndAbort,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestZEVMToEVMCallThroughContractName,
		"zevm -> evm call through intermediary contract",
		[]runner.ArgDefinition{},
		TestZEVMToEVMCallThroughContract,
	),
	runner.NewE2ETest(
		TestEVMToZEVMCallName,
		"evm -> zevm call",
		[]runner.ArgDefinition{},
		TestEVMToZEVMCall,
	),
	runner.NewE2ETest(
		TestEVMToZEVMCallAbortName,
		"evm -> zevm call fails and abort with onAbort",
		[]runner.ArgDefinition{},
		TestEVMToZEVMCallAbort,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestDepositAndCallSwapName,
		"evm -> zevm deposit and call with swap and withdraw back to evm",
		[]runner.ArgDefinition{},
		TestDepositAndCallSwap,
	),

	/*
	 Solana tests
	*/
	runner.NewE2ETest(
		TestSolanaDepositName,
		"deposit SOL into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "24000000"},
		},
		TestSolanaDeposit,
	),
	runner.NewE2ETest(
		TestSolanaDepositThroughProgramName,
		"deposit SOL into ZEVM through example connected program",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "24000000"},
		},
		TestSolanaDepositThroughProgram,
	),
	runner.NewE2ETest(
		TestSolanaWithdrawName,
		"withdraw SOL from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdraw,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaWithdrawAndCallName,
		"withdraw SOL from ZEVM and call solana program",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdrawAndCall,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaWithdrawRevertExecutableReceiverName,
		"withdraw SOL from ZEVM reverts if executable receiver",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdrawRevertExecutableReceiver,
	),
	runner.NewE2ETest(
		TestZEVMToSolanaCallName,
		"call solana program from ZEVM",
		[]runner.ArgDefinition{},
		TestZEVMToSolanaCall,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaWithdrawAndCallInvalidTxSizeName,
		"withdraw SOL from ZEVM and call solana program with invalid tx size",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdrawAndCallInvalidTxSize,
	),
	runner.NewE2ETest(
		TestSolanaWithdrawAndCallInvalidMsgEncodingName,
		"withdraw SOL from ZEVM and call solana program with invalid msg encoding",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdrawAndCallInvalidMsgEncoding,
	),
	runner.NewE2ETest(
		TestSolanaWithdrawAndCallRevertWithCallName,
		"withdraw SOL from ZEVM and call solana program that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSolanaWithdrawAndCallRevertWithCall,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSPLWithdrawAndCallName,
		"withdraw SPL from ZEVM and call solana program",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSPLWithdrawAndCall,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSPLWithdrawAndCallRevertName,
		"withdraw SPL from ZEVM and call solana program that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1000000"},
		},
		TestSPLWithdrawAndCallRevert,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallName,
		"deposit SOL into ZEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1200000"},
		},
		TestSolanaDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaToZEVMCallName,
		"call a zevm contract",
		[]runner.ArgDefinition{},
		TestSolanaToZEVMCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaToZEVMCallAbortName,
		"call a zevm contract and abort",
		[]runner.ArgDefinition{},
		TestSolanaToZEVMCallAbort,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSPLWithdrawName,
		"withdraw SPL from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in spl tokens", DefaultValue: "100000"},
		},
		TestSPLWithdraw,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSPLWithdrawAndCreateReceiverAtaName,
		"withdraw SPL from ZEVM and create receiver ata",
		[]runner.ArgDefinition{
			{Description: "amount in spl tokens", DefaultValue: "1000000"},
		},
		TestSPLWithdrawAndCreateReceiverAta,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallRevertName,
		"deposit SOL into ZEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1200000"},
		},
		TestSolanaDepositAndCallRevert,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallRevertWithCallName,
		"deposit SOL into ZEVM and call a contract that reverts with call on revert",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1200000"},
		},
		TestSolanaDepositAndCallRevertWithCall,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallRevertWithCallThatRevertsName,
		"deposit SOL into ZEVM and call a contract that reverts with call on revert and connected program reverts",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1200000"},
		},
		TestSolanaDepositAndCallRevertWithCallThatReverts,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallRevertWithDustName,
		"deposit SOL into ZEVM; revert with dust amount that aborts the CCTX",
		[]runner.ArgDefinition{},
		TestSolanaDepositAndCallRevertWithDust,
		runner.WithMinimumVersion("v29.0.0"),
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
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSolanaWhitelistSPLName,
		"whitelist SPL",
		[]runner.ArgDefinition{},
		TestSolanaWhitelistSPL,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSPLDepositName,
		"deposit SPL into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount of spl tokens", DefaultValue: "24000000"},
		},
		TestSPLDeposit,
	),
	runner.NewE2ETest(
		TestSPLDepositAndCallName,
		"deposit SPL into ZEVM and call",
		[]runner.ArgDefinition{
			{Description: "amount of spl tokens", DefaultValue: "12000000"},
		},
		TestSPLDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSPLDepositAndCallRevertName,
		"deposit SPL into ZEVM and call which reverts",
		[]runner.ArgDefinition{
			{Description: "amount of spl tokens", DefaultValue: "12000000"},
		},
		TestSPLDepositAndCallRevert,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSPLDepositAndCallRevertWithCallName,
		"deposit SPL into ZEVM and call which reverts with call on revert",
		[]runner.ArgDefinition{
			{Description: "amount of spl tokens", DefaultValue: "12000000"},
		},
		TestSPLDepositAndCallRevertWithCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSPLDepositAndCallRevertWithCallThatRevertsName,
		"deposit SPL into ZEVM and call which reverts with call on revert that reverts",
		[]runner.ArgDefinition{
			{Description: "amount of spl tokens", DefaultValue: "12000000"},
		},
		TestSPLDepositAndCallRevertWithCallThatReverts,
		runner.WithMinimumVersion("v30.0.0"),
	),
	/*
	 TON tests
	*/
	runner.NewE2ETest(
		TestTONDepositName,
		"deposit TON into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "1000000000"}, // 1.0 TON
		},
		TestTONDeposit,
	),
	runner.NewE2ETest(
		TestTONDepositAndCallName,
		"deposit TON into ZEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "1000000000"}, // 1.0 TON
		},
		TestTONDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestTONDepositAndCallRefundName,
		"deposit TON into ZEVM and call a smart contract that reverts; expect refund",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "1000000000"}, // 1.0 TON
		},
		TestTONDepositAndCallRefund,
	),
	runner.NewE2ETest(
		TestTONDepositRestrictedName,
		"deposit TON into ZEVM restricted address",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "100000000"}, // 0.1 TON
		},
		TestTONDepositRestricted,
	),
	runner.NewE2ETest(
		TestTONCallName,
		"call TON into ZEVM",
		[]runner.ArgDefinition{},
		TestTONToZEVMCall,
	),
	runner.NewE2ETest(
		TestTONWithdrawName,
		"withdraw TON from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "2000000000"}, // 2.0 TON
		},
		TestTONWithdraw,
	),
	runner.NewE2ETest(
		TestTONWithdrawRestrictedName,
		"withdraw TON from ZEVM to restricted address (compliance check)",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "100000000"}, // 0.1 TON
		},
		TestTONWithdrawRestricted,
	),
	runner.NewE2ETest(
		// TON address starts with an chain index (0:... or -1:...)
		// Zetachain operates only on base chain (0:)
		TestTONWithdrawMasterchainName,
		"withdraw TON from ZEVM to masterchain that is a consensus chain rather than a base workchain",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "100000000"}, // 0.1 TON
		},
		TestTONWithdrawMasterchain,
	),
	runner.NewE2ETest(
		TestTONWithdrawConcurrentName,
		"withdraw TON from ZEVM for several recipients simultaneously",
		[]runner.ArgDefinition{},
		TestTONWithdrawConcurrent,
	),
	/*
	 Sui tests
	*/
	runner.NewE2ETest(
		TestSuiDepositName,
		"deposit SUI into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "10000000000"},
		},
		TestSuiDeposit,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSuiDepositAndCallName,
		"deposit SUI into ZEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestSuiDepositAndCallRevertName,
		"deposit SUI into ZEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "10000000000"},
		},
		TestSuiDepositAndCallRevert,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSuiTokenDepositName,
		"deposit fungible token SUI into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in base unit", DefaultValue: "10000000000"},
		},
		TestSuiTokenDeposit,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSuiTokenDepositAndCallName,
		"deposit fungible token into ZEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in base unit", DefaultValue: "1000000"},
		},
		TestSuiTokenDepositAndCall,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSuiTokenDepositAndCallRevertName,
		"deposit fungible token into ZEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in base unit", DefaultValue: "10000000000"},
		},
		TestSuiTokenDepositAndCallRevert,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestSuiWithdrawName,
		"withdraw SUI from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiWithdraw,
		runner.WithMinimumVersion("v33.0.0"),
	),
	runner.NewE2ETest(
		TestSuiWithdrawAndCallName,
		"withdraw SUI from ZEVM and makes an authenticated call to a contract",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "1000000"},
			{Description: "gas limit for withdraw and call", DefaultValue: "100000"},
		},
		TestSuiWithdrawAndCall,
		runner.WithMinimumVersion("v35.0.0"),
	),
	runner.NewE2ETest(
		TestSuiWithdrawRevertWithCallName,
		"withdraw SUI from ZEVM that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiWithdrawRevertWithCall,
		runner.WithMinimumVersion("v33.0.0"),
	),
	runner.NewE2ETest(
		TestSuiWithdrawAndCallInvalidPayloadName,
		"withdraw SUI from ZEVM and makes an authenticated call to a contract that reverts due to invalid payload",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiWithdrawAndCallInvalidPayload,
		runner.WithMinimumVersion("v35.0.0"),
	),
	runner.NewE2ETest(
		TestSuiWithdrawAndCallRevertWithCallName,
		"withdraw SUI from ZEVM and makes an authenticated call to a contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "1000000"},
			{Description: "gas limit for withdraw and call", DefaultValue: "100000"},
		},
		TestSuiWithdrawAndCallRevertWithCall,
		runner.WithMinimumVersion("v35.0.0"),
	),
	runner.NewE2ETest(
		TestSuiTokenWithdrawName,
		"withdraw fungible token from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in base unit", DefaultValue: "100000"},
		},
		TestSuiTokenWithdraw,
		runner.WithMinimumVersion("v33.0.0"),
	),
	runner.NewE2ETest(
		TestSuiTokenWithdrawAndCallName,
		"withdraw fungible token from ZEVM and makes an authenticated call to a contract",
		[]runner.ArgDefinition{
			{Description: "amount in base unit", DefaultValue: "100000"},
			{Description: "gas limit for withdraw and call", DefaultValue: "100000"},
		},
		TestSuiTokenWithdrawAndCall,
		runner.WithMinimumVersion("v35.0.0"),
	),
	runner.NewE2ETest(
		TestSuiTokenWithdrawAndCallRevertWithCallName,
		"withdraw fungible token from ZEVM and makes an authenticated call to a contract that reverts with a onRevert call",
		[]runner.ArgDefinition{
			{Description: "amount in base unit", DefaultValue: "100000"},
			{Description: "gas limit for withdraw and call", DefaultValue: "100000"},
		},
		TestSuiTokenWithdrawAndCallRevertWithCall,
		runner.WithMinimumVersion("v35.0.0"),
	),
	runner.NewE2ETest(
		TestSuiDepositRestrictedName,
		"deposit SUI into ZEVM restricted address",
		[]runner.ArgDefinition{
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiDepositRestrictedAddress,
	),
	runner.NewE2ETest(
		TestSuiWithdrawRestrictedName,
		"withdraw SUI from ZEVM to restricted address",
		[]runner.ArgDefinition{
			{Description: "receiver", DefaultValue: sample.RestrictedSuiAddressTest},
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiWithdrawRestrictedAddress,
		runner.WithMinimumVersion("v33.0.0"),
	),
	runner.NewE2ETest(
		TestSuiWithdrawInvalidReceiverName,
		"withdraw SUI from ZEVM to invalid receiver address",
		[]runner.ArgDefinition{
			{Description: "receiver", DefaultValue: "0x547a07f0564e0c8d48c4ae53305eabdef87e9610"},
			{Description: "amount in mist", DefaultValue: "1000000"},
		},
		TestSuiWithdrawInvalidReceiver,
		runner.WithMinimumVersion("v33.0.0"),
	),

	/*
	 Bitcoin tests
	*/
	runner.NewE2ETest(
		TestBitcoinDonationName,
		"donate Bitcoin to TSS address", []runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.1"},
		},
		TestBitcoinDonation,
	),
	runner.NewE2ETest(
		TestBitcoinDepositName,
		"deposit Bitcoin into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "1.0"},
		},
		TestBitcoinDeposit,
	),
	runner.NewE2ETest(
		TestBitcoinDepositFastConfirmationName,
		"deposit Bitcoin into ZEVM using fast confirmation",
		[]runner.ArgDefinition{},
		TestBitcoinDepositFastConfirmation,
		runner.WithMinimumVersion("v37.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinDepositAndCallName,
		"deposit Bitcoin into ZEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinDepositAndCallRevertName,
		"deposit Bitcoin into ZEVM; expect refund",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.1"},
		},
		TestBitcoinDepositAndCallRevert,
	),
	runner.NewE2ETest(
		TestBitcoinDepositAndCallRevertWithDustName,
		"deposit Bitcoin into ZEVM; revert with dust amount that aborts the CCTX",
		[]runner.ArgDefinition{},
		TestBitcoinDepositAndCallRevertWithDust,
	),
	runner.NewE2ETest(
		TestBitcoinDepositAndWithdrawWithDustName,
		"deposit Bitcoin into ZEVM and withdraw with dust amount that fails the CCTX",
		[]runner.ArgDefinition{},
		TestBitcoinDepositAndWithdrawWithDust,
	),
	runner.NewE2ETest(
		TestBitcoinStdMemoDepositName,
		"deposit Bitcoin into ZEVM with standard memo",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.2"},
		},
		TestBitcoinStdMemoDeposit,
	),
	runner.NewE2ETest(
		TestBitcoinStdMemoDepositAndCallName,
		"deposit Bitcoin into ZEVM and call a contract with standard memo",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.5"},
		},
		TestBitcoinStdMemoDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinStdMemoDepositAndCallRevertName,
		"deposit Bitcoin into ZEVM and call a contract with standard memo; expect revert",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.1"},
		},
		TestBitcoinStdMemoDepositAndCallRevert,
	),
	runner.NewE2ETest(
		TestBitcoinStdMemoDepositAndCallRevertOtherAddressName,
		"deposit Bitcoin into ZEVM and call a contract with standard memo; expect revert to other address",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.1"},
		},
		TestBitcoinStdMemoDepositAndCallRevertOtherAddress,
	),
	runner.NewE2ETest(
		TestBitcoinStdMemoDepositAndCallRevertAndAbortName,
		"deposit Bitcoin into ZEVM and call a contract with standard memo; revert and abort with onAbort",
		[]runner.ArgDefinition{},
		TestBitcoinStdMemoDepositAndCallRevertAndAbort,
		runner.WithMinimumVersion("v37.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinStdMemoInscribedDepositAndCallName,
		"deposit Bitcoin into ZEVM and call a contract with inscribed standard memo",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.1"},
		},
		TestBitcoinStdMemoInscribedDepositAndCall,
		runner.WithMinimumVersion("v32.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinDepositAndAbortWithLowDepositFeeName,
		"deposit Bitcoin into ZEVM that aborts due to insufficient deposit fee",
		[]runner.ArgDefinition{},
		TestBitcoinDepositAndAbortWithLowDepositFee,
		runner.WithMinimumVersion("v27.0.0"),
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
			{Description: "receiver", DefaultValue: sample.RestrictedBtcAddressTest},
			{Description: "amount in btc", DefaultValue: "0.001"},
			{Description: "revert address", DefaultValue: sample.RevertAddressZEVM},
		},
		TestBitcoinWithdrawRestricted,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinDepositInvalidMemoRevertName,
		"deposit Bitcoin with invalid memo; expect revert",
		[]runner.ArgDefinition{},
		TestBitcoinDepositInvalidMemoRevert,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestBitcoinWithdrawRBFName,
		"withdraw Bitcoin from ZEVM and replace the outbound using RBF",
		[]runner.ArgDefinition{
			{Description: "receiver address", DefaultValue: ""},
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinWithdrawRBF,
		runner.WithDependencies(DepdencyAllBitcoinDeposits),
	),
	/*
	 Application tests
	*/
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
		TestDonationEtherName,
		"donate Ether to the TSS",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000000000000000"},
		},
		TestDonationEther,
	),
	runner.NewE2ETest(
		TestInboundTrackersName,
		"test processing inbound trackers for observation",
		[]runner.ArgDefinition{},
		TestInboundTrackers,
	),
	runner.NewE2ETest(
		TestPrecompilesName,
		"test precompiles on ZEVM",
		[]runner.ArgDefinition{},
		TestPrecompiles,
		runner.WithMinimumVersion("v33.0.0"),
	),
	runner.NewE2ETest(

		TestOpcodesName,
		"test opcodes support in ZEVM",
		[]runner.ArgDefinition{},
		TestOpcodes,
	),
	/*
	 Stress tests
	*/
	runner.NewE2ETest(
		TestStressEtherWithdrawName,
		"stress test Ether withdrawal",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
			{Description: CountArgDescription, DefaultValue: "100"},
		},
		TestStressEtherWithdraw,
	),
	runner.NewE2ETest(
		TestStressBTCWithdrawName,
		"stress test BTC withdrawal",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.01"},
			{Description: CountArgDescription, DefaultValue: "100"},
		},
		TestStressBTCWithdraw,
	),
	runner.NewE2ETest(
		TestStressEtherDepositName,
		"stress test Ether deposit",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
			{Description: CountArgDescription, DefaultValue: "100"},
		},
		TestStressEtherDeposit,
	),
	runner.NewE2ETest(
		TestStressBTCDepositName,
		"stress test BTC deposit",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.001"},
			{Description: CountArgDescription, DefaultValue: "100"},
		},
		TestStressBTCDeposit,
	),
	runner.NewE2ETest(
		TestStressSolanaDepositName,
		"stress test SOL deposit",
		[]runner.ArgDefinition{
			{Description: "amount in lamports", DefaultValue: "1200000"},
			{Description: CountArgDescription, DefaultValue: "50"},
		},
		TestStressSolanaDeposit,
	),
	runner.NewE2ETest(
		TestStressSPLDepositName,
		"stress test SPL deposit",
		[]runner.ArgDefinition{
			{Description: "amount in SPL tokens", DefaultValue: "1200000"},
			{Description: CountArgDescription, DefaultValue: "50"},
		},
		TestStressSPLDeposit,
	),
	runner.NewE2ETest(
		TestStressSolanaWithdrawName,
		"stress test SOL withdrawals",
		[]runner.ArgDefinition{
			{Description: "amount in lamports", DefaultValue: "1000000"},
			{Description: CountArgDescription, DefaultValue: "50"},
		},
		TestStressSolanaWithdraw,
	),
	runner.NewE2ETest(
		TestStressSPLWithdrawName,
		"stress test SPL withdrawals",
		[]runner.ArgDefinition{
			{Description: "amount in SPL tokens", DefaultValue: "1000000"},
			{Description: CountArgDescription, DefaultValue: "50"},
		},
		TestStressSPLWithdraw,
	),
	runner.NewE2ETest(
		TestStressSuiDepositName,
		"stress test SUI deposits",
		[]runner.ArgDefinition{
			{Description: "amount in SUI", DefaultValue: "1000000"},
			{Description: CountArgDescription, DefaultValue: "50"},
		},
		TestStressSuiDeposit,
	),
	runner.NewE2ETest(
		TestStressSuiWithdrawName,
		"stress test SUI withdrawals",
		[]runner.ArgDefinition{
			{Description: "amount in SUI", DefaultValue: "1000000"},
			{Description: CountArgDescription, DefaultValue: "50"},
		},
		TestStressSuiWithdraw,
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
		legacy.TestRateLimiter,
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
	runner.NewE2ETest(
		TestMigrateConnectorFundsName,
		"migrate connector funds from V1 to V2 connector",
		[]runner.ArgDefinition{},
		TestMigrateConnectorFunds,
	),
	runner.NewE2ETest(
		TestUpdateZRC20NameName,
		"update ZRC20 name and symbol",
		[]runner.ArgDefinition{},
		TestUpdateZRC20Name,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestZetaclientRestartHeightName,
		"zetaclient scheduled restart height",
		[]runner.ArgDefinition{},
		TestZetaclientRestartHeight,
	),
	runner.NewE2ETest(
		TestZetaclientSignerOffsetName,
		"zetaclient signer offset",
		[]runner.ArgDefinition{},
		TestZetaclientSignerOffset,
	),
	runner.NewE2ETest(
		TestUpdateOperationalChainParamsName,
		"update operational chain params",
		[]runner.ArgDefinition{},
		TestUpdateOperationalChainParams,
		runner.WithMinimumVersion("v29.0.0"),
	),
	runner.NewE2ETest(
		TestBurnFungibleModuleAssetName,
		"burn fungible module asset",
		[]runner.ArgDefinition{},
		TestBurnFungibleModuleAsset,
		runner.WithMinimumVersion("v33.0.0"),
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
	runner.NewE2ETest(
		TestOperationAddLiquidityBTCName,
		"add liquidity to the ZETA/BTC pool",
		[]runner.ArgDefinition{
			{Description: "amountZETA", DefaultValue: "50000000000000000000"},
			{Description: "amountBTC", DefaultValue: "5000000000"},
		},
		TestOperationAddLiquidityBTC,
	),
	runner.NewE2ETest(
		TestOperationAddLiquiditySOLName,
		"add liquidity to the ZETA/SOL pool",
		[]runner.ArgDefinition{
			{Description: "amountZETA", DefaultValue: "50000000000000000000"},
			{Description: "amountSOL", DefaultValue: "50000000000"},
		},
		TestOperationAddLiquiditySOL,
	),
	runner.NewE2ETest(
		TestOperationAddLiquiditySPLName,
		"add liquidity to the ZETA/SPL pool",
		[]runner.ArgDefinition{
			{Description: "amountZETA", DefaultValue: "50000000000000000000"},
			{Description: "amountSPL", DefaultValue: "50000000000000000000"},
		},
		TestOperationAddLiquiditySPL,
	),
	runner.NewE2ETest(
		TestOperationAddLiquiditySUIName,
		"add liquidity to the ZETA/SUI pool",
		[]runner.ArgDefinition{
			{Description: "amountZETA", DefaultValue: "50000000000000000000"},
			{Description: "amountSUI", DefaultValue: "50000000000"},
		},
		TestOperationAddLiquiditySUI,
	),
	runner.NewE2ETest(
		TestOperationAddLiquiditySuiFungibleTokenName,
		"add liquidity to the ZETA/SuiFungibleToken pool",
		[]runner.ArgDefinition{
			{Description: "amountZETA", DefaultValue: "50000000000000000000"},
			{Description: "amountSuiFungibleToken", DefaultValue: "50000000"},
		},
		TestOperationAddLiquiditySuiFungibleToken,
	),
	runner.NewE2ETest(
		TestOperationAddLiquidityTONName,
		"add liquidity to the ZETA/TON pool",
		[]runner.ArgDefinition{
			{Description: "amountZETA", DefaultValue: "50000000000000000000"},
			{Description: "amountTON", DefaultValue: "50000000000"},
		},
		TestOperationAddLiquidityTON,
	),
	/*
	 Legacy tests
	*/
	runner.NewE2ETest(
		TestLegacyMessagePassingExternalChainsName,
		"evm->evm message passing (sending ZETA only) (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000000"},
		},
		legacy.TestMessagePassingExternalChains,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingRevertFailExternalChainsName,
		"message passing with failing revert between external EVM chains (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000000"},
		},
		legacy.TestMessagePassingRevertFailExternalChains,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingRevertSuccessExternalChainsName,
		"message passing with successful revert between external EVM chains (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000000"},
		},
		legacy.TestMessagePassingRevertSuccessExternalChains,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingEVMtoZEVMName,
		"evm -> zevm message passing contract call (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000009"},
		},
		legacy.TestMessagePassingEVMtoZEVM,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingZEVMToEVMName,
		"zevm -> evm message passing contract call (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000007"},
		},
		legacy.TestMessagePassingZEVMtoEVM,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingZEVMtoEVMRevertName,
		"zevm -> evm message passing contract call reverts (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000006"},
		},
		legacy.TestMessagePassingZEVMtoEVMRevert,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingEVMtoZEVMRevertName,
		"evm -> zevm message passing and revert back to evm (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000008"},
		},
		legacy.TestMessagePassingEVMtoZEVMRevert,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingZEVMtoEVMRevertFailName,
		"zevm -> evm message passing contract with failing revert (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000008"},
		},
		legacy.TestMessagePassingZEVMtoEVMRevertFail,
	),
	runner.NewE2ETest(
		TestLegacyMessagePassingEVMtoZEVMRevertFailName,
		"evm -> zevm message passing contract with failing revert (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000008"},
		},
		legacy.TestMessagePassingEVMtoZEVMRevertFail,
	),
	runner.NewE2ETest(
		TestLegacyEtherDepositName,
		"deposit Ether into ZEVM (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000"},
		},
		legacy.TestEtherDeposit,
	),
	runner.NewE2ETest(
		TestLegacyEtherWithdrawName,
		"withdraw Ether from ZEVM (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		legacy.TestEtherWithdraw,
	),
	runner.NewE2ETest(
		TestEtherWithdrawRestrictedName,
		"withdraw Ether from ZEVM to restricted address (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "receiver", DefaultValue: sample.RestrictedEVMAddressTest},
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestEtherWithdrawRestricted,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestLegacyEtherDepositAndCallRefundName,
		"deposit Ether into ZEVM and call a contract that reverts; should refund (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "10000000000000000000"},
		},
		legacy.TestEtherDepositAndCallRefund,
	),
	runner.NewE2ETest(
		TestLegacyEtherDepositAndCallName,
		"deposit ZRC20 into ZEVM and call a contract (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "1000000000000000000"},
		},
		legacy.TestEtherDepositAndCall,
		runner.WithMinimumVersion("v30.0.0"),
	),
	runner.NewE2ETest(
		TestLegacyERC20WithdrawName,
		"withdraw ERC20 from ZEVM (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
		},
		legacy.TestERC20Withdraw,
	),
	runner.NewE2ETest(
		TestLegacyERC20DepositName,
		"deposit ERC20 into ZEVM (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000"},
		},
		legacy.TestERC20Deposit,
	),
	runner.NewE2ETest(
		TestLegacyMultipleERC20DepositName,
		"deposit ERC20 into ZEVM in multiple deposits (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000000000"},
			{Description: CountArgDescription, DefaultValue: "3"},
		},
		legacy.TestMultipleERC20Deposit,
	),
	runner.NewE2ETest(
		TestLegacyMultipleERC20WithdrawsName,
		"withdraw ERC20 from ZEVM in multiple withdrawals (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100"},
			{Description: CountArgDescription, DefaultValue: "3"},
		},
		legacy.TestMultipleERC20Withdraws,
	),
	runner.NewE2ETest(
		TestERC20DepositRestrictedName,
		"deposit ERC20 into ZEVM restricted address (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100000"},
		},
		TestERC20DepositRestricted,
	),
	runner.NewE2ETest(
		TestLegacyERC20DepositAndCallRefundName,
		"deposit a non-gas ZRC20 into ZEVM and call a contract that reverts (v1 protocol contracts)",
		[]runner.ArgDefinition{},
		legacy.TestERC20DepositAndCallRefund,
	),

	/*
	 ZETA tests
	*/
	runner.NewE2ETest(
		TestLegacyZetaDepositName,
		"deposit ZETA from Ethereum to ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "1000000000000000000"},
		},
		legacy.TestZetaDeposit,
	),
	runner.NewE2ETest(
		TestLegacyZetaDepositAndCallAbortName,
		"deposit and ZETA from Ethereum to ZEVM and call a contract.The cctx reverts and then aborts",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "1000000000000000000"},
		},
		legacy.TestZetaDepositAndCallAbort,
	),
	runner.NewE2ETest(
		TestLegacyZetaDepositNewAddressName,
		"deposit ZETA from Ethereum to a new ZEVM address which does not exist yet",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "1000000000000000000"},
		},
		legacy.TestZetaDepositNewAddress,
	),
	runner.NewE2ETest(
		TestLegacyZetaDepositRestrictedName,
		"deposit ZETA from Ethereum to ZEVM restricted address",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "1000000000000000000"},
		},
		legacy.TestZetaDepositRestricted,
	),
	runner.NewE2ETest(
		TestLegacyZetaWithdrawName,
		"withdraw ZETA from ZEVM to Ethereum",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "10000000000000000000"},
		},
		legacy.TestZetaWithdraw,
	),
	runner.NewE2ETest(
		TestLegacyZetaWithdrawBTCRevertName,
		"sending ZETA from ZEVM to Bitcoin with a message that should revert cctxs",
		[]runner.ArgDefinition{
			{Description: "amount in azeta", DefaultValue: "1000000000000000000"},
		},
		legacy.TestZetaWithdrawBTCRevert,
	),
	runner.NewE2ETest(
		TestUndelegateToBelowMinimumObserverDelegation,
		"test undelegating to below minimum observer delegation",
		[]runner.ArgDefinition{},
		UndelegateToBelowMinimumObserverDelegation),
}
