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
	TestETHDepositName                         = "eth_deposit"
	TestETHDepositAndCallName                  = "eth_deposit_and_call"
	TestETHDepositAndCallNoMessageName         = "eth_deposit_and_call_no_message"
	TestETHDepositAndCallRevertName            = "eth_deposit_and_call_revert"
	TestETHDepositAndCallRevertWithCallName    = "eth_deposit_and_call_revert_with_call"
	TestETHWithdrawName                        = "eth_withdraw"
	TestETHWithdrawAndArbitraryCallName        = "eth_withdraw_and_arbitrary_call"
	TestETHWithdrawAndCallName                 = "eth_withdraw_and_call"
	TestETHWithdrawAndCallNoMessageName        = "eth_withdraw_and_call_no_message"
	TestETHWithdrawAndCallThroughContractName  = "eth_withdraw_and_call_through_contract"
	TestETHWithdrawAndCallRevertName           = "eth_withdraw_and_call_revert"
	TestETHWithdrawAndCallRevertWithCallName   = "eth_withdraw_and_call_revert_with_call"
	TestDepositAndCallOutOfGasName             = "deposit_and_call_out_of_gas"
	TestERC20DepositName                       = "erc20_deposit"
	TestERC20DepositAndCallName                = "erc20_deposit_and_call"
	TestERC20DepositAndCallNoMessageName       = "erc20_deposit_and_call_no_message"
	TestERC20DepositAndCallRevertName          = "erc20_deposit_and_call_revert"
	TestERC20DepositAndCallRevertWithCallName  = "erc20_deposit_and_call_revert_with_call"
	TestERC20WithdrawName                      = "erc20_withdraw"
	TestERC20WithdrawAndArbitraryCallName      = "erc20_withdraw_and_arbitrary_call"
	TestERC20WithdrawAndCallName               = "erc20_withdraw_and_call"
	TestERC20WithdrawAndCallNoMessageName      = "erc20_withdraw_and_call_no_message"
	TestERC20WithdrawAndCallRevertName         = "erc20_withdraw_and_call_revert"
	TestERC20WithdrawAndCallRevertWithCallName = "erc20_withdraw_and_call_revert_with_call"
	TestZEVMToEVMArbitraryCallName             = "zevm_to_evm_arbitrary_call"
	TestZEVMToEVMCallName                      = "zevm_to_evm_call"
	TestZEVMToEVMCallThroughContractName       = "zevm_to_evm_call_through_contract"
	TestEVMToZEVMCallName                      = "evm_to_zevm_call"
	TestDepositAndCallSwapName                 = "deposit_and_call_swap"
	TestEtherWithdrawRestrictedName            = "eth_withdraw_restricted"
	TestERC20DepositRestrictedName             = "erc20_deposit_restricted" // #nosec G101: Potential hardcoded credentials (gosec), not a credential

	/*
	 * Solana tests
	 */
	TestSolanaDepositName                      = "solana_deposit"
	TestSolanaWithdrawName                     = "solana_withdraw"
	TestSolanaDepositAndCallName               = "solana_deposit_and_call"
	TestSolanaDepositAndCallRevertName         = "solana_deposit_and_call_revert"
	TestSolanaDepositAndCallRevertWithDustName = "solana_deposit_and_call_revert_with_dust"
	TestSolanaDepositRestrictedName            = "solana_deposit_restricted"
	TestSolanaWithdrawRestrictedName           = "solana_withdraw_restricted"
	TestSPLDepositName                         = "spl_deposit"
	TestSPLDepositAndCallName                  = "spl_deposit_and_call"
	TestSPLWithdrawName                        = "spl_withdraw"
	TestSPLWithdrawAndCreateReceiverAtaName    = "spl_withdraw_and_create_receiver_ata"

	/**
	 * TON tests
	 */
	TestTONDepositName              = "ton_deposit"
	TestTONDepositAndCallName       = "ton_deposit_and_call"
	TestTONDepositAndCallRefundName = "ton_deposit_refund"
	TestTONWithdrawName             = "ton_withdraw"
	TestTONWithdrawConcurrentName   = "ton_withdraw_concurrent"

	/*
	 Bitcoin tests
	 Test transfer of Bitcoin asset across chains
	*/
	TestBitcoinDepositName                                 = "bitcoin_deposit"
	TestBitcoinDepositAndCallName                          = "bitcoin_deposit_and_call"
	TestBitcoinDepositAndCallRevertName                    = "bitcoin_deposit_and_call_revert"
	TestBitcoinDepositAndCallRevertWithDustName            = "bitcoin_deposit_and_call_revert_with_dust"
	TestBitcoinDonationName                                = "bitcoin_donation"
	TestBitcoinStdMemoDepositName                          = "bitcoin_std_memo_deposit"
	TestBitcoinStdMemoDepositAndCallName                   = "bitcoin_std_memo_deposit_and_call"
	TestBitcoinStdMemoDepositAndCallRevertName             = "bitcoin_std_memo_deposit_and_call_revert"
	TestBitcoinStdMemoDepositAndCallRevertOtherAddressName = "bitcoin_std_memo_deposit_and_call_revert_other_address"
	TestBitcoinStdMemoInscribedDepositAndCallName          = "bitcoin_std_memo_inscribed_deposit_and_call"
	TestBitcoinWithdrawSegWitName                          = "bitcoin_withdraw_segwit"
	TestBitcoinWithdrawTaprootName                         = "bitcoin_withdraw_taproot"
	TestBitcoinWithdrawMultipleName                        = "bitcoin_withdraw_multiple"
	TestBitcoinWithdrawLegacyName                          = "bitcoin_withdraw_legacy"
	TestBitcoinWithdrawP2WSHName                           = "bitcoin_withdraw_p2wsh"
	TestBitcoinWithdrawP2SHName                            = "bitcoin_withdraw_p2sh"
	TestBitcoinWithdrawInvalidAddressName                  = "bitcoin_withdraw_invalid"
	TestBitcoinWithdrawRestrictedName                      = "bitcoin_withdraw_restricted"

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
	TestMigrateTSSName                = "migrate_tss"
	TestSolanaWhitelistSPLName        = "solana_whitelist_spl"
	TestZetaclientRestartHeightName   = "zetaclient_restart_height"
	TestZetaclientSignerOffsetName    = "zetaclient_signer_offset"

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
	TestPrecompilesPrototypeName                 = "precompile_contracts_prototype"
	TestPrecompilesPrototypeThroughContractName  = "precompile_contracts_prototype_through_contract"
	TestPrecompilesStakingName                   = "precompile_contracts_staking"
	TestPrecompilesStakingThroughContractName    = "precompile_contracts_staking_through_contract"
	TestPrecompilesBankName                      = "precompile_contracts_bank"
	TestPrecompilesBankFailName                  = "precompile_contracts_bank_fail"
	TestPrecompilesBankThroughContractName       = "precompile_contracts_bank_through_contract"
	TestPrecompilesDistributeName                = "precompile_contracts_distribute"
	TestPrecompilesDistributeNonZRC20Name        = "precompile_contracts_distribute_non_zrc20"
	TestPrecompilesDistributeThroughContractName = "precompile_contracts_distribute_through_contract"

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
	TestLegacyZetaDepositName           = "legacy_zeta_deposit"
	TestLegacyZetaDepositNewAddressName = "legacy_zeta_deposit_new_address"
	TestLegacyZetaDepositRestrictedName = "legacy_zeta_deposit_restricted"
	TestLegacyZetaWithdrawName          = "legacy_zeta_withdraw"
	TestLegacyZetaWithdrawBTCRevertName = "legacy_zeta_withdraw_btc_revert" // #nosec G101 - not a hardcoded password

)

// AllE2ETests is an ordered list of all e2e tests
var AllE2ETests = []runner.E2ETest{
	/*
	 EVM chain tests
	*/
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
		TestETHWithdrawName,
		"withdraw Ether from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestETHWithdraw,
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
		},
		TestETHWithdrawAndCall,
	),
	runner.NewE2ETest(
		TestETHWithdrawAndCallNoMessageName,
		"withdraw Ether from ZEVM call a contract with no message content",
		[]runner.ArgDefinition{
			{Description: "amount in wei", DefaultValue: "100000"},
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
		},
		TestERC20WithdrawAndCall,
	),
	runner.NewE2ETest(
		TestERC20WithdrawAndCallNoMessageName,
		"withdraw ERC20 from ZEVM and authenticated call a contract with no message",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "1000"},
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
		TestZEVMToEVMArbitraryCallName,
		"zevm -> evm call",
		[]runner.ArgDefinition{},
		TestZEVMToEVMArbitraryCall,
	),
	runner.NewE2ETest(
		TestZEVMToEVMCallName,
		"zevm -> evm call",
		[]runner.ArgDefinition{},
		TestZEVMToEVMCall,
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
		TestSPLWithdrawName,
		"withdraw SPL from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in spl tokens", DefaultValue: "1000000"},
		},
		TestSPLWithdraw,
	),
	runner.NewE2ETest(
		TestSPLWithdrawAndCreateReceiverAtaName,
		"withdraw SPL from ZEVM and create receiver ata",
		[]runner.ArgDefinition{
			{Description: "amount in spl tokens", DefaultValue: "1000000"},
		},
		TestSPLWithdrawAndCreateReceiverAta,
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallRevertName,
		"deposit SOL into ZEVM and call a contract that reverts",
		[]runner.ArgDefinition{
			{Description: "amount in lamport", DefaultValue: "1200000"},
		},
		TestSolanaDepositAndCallRevert,
	),
	runner.NewE2ETest(
		TestSolanaDepositAndCallRevertWithDustName,
		"deposit SOL into ZEVM; revert with dust amount that aborts the CCTX",
		[]runner.ArgDefinition{},
		TestSolanaDepositAndCallRevertWithDust,
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
	runner.NewE2ETest(
		TestSolanaWhitelistSPLName,
		"whitelist SPL",
		[]runner.ArgDefinition{},
		TestSolanaWhitelistSPL,
	),
	runner.NewE2ETest(
		TestSPLDepositName,
		"deposit SPL into ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount of spl tokens", DefaultValue: "12000000"},
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
		TestTONWithdrawName,
		"withdraw TON from ZEVM",
		[]runner.ArgDefinition{
			{Description: "amount in nano tons", DefaultValue: "2000000000"}, // 2.0 TON
		},
		TestTONWithdraw,
	),
	runner.NewE2ETest(
		TestTONWithdrawConcurrentName,
		"withdraw TON from ZEVM for several recipients simultaneously",
		[]runner.ArgDefinition{},
		TestTONWithdrawConcurrent,
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
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinDeposit,
	),
	runner.NewE2ETest(
		TestBitcoinDepositAndCallName,
		"deposit Bitcoin into ZEVM and call a contract",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.001"},
		},
		TestBitcoinDepositAndCall,
	),
	runner.NewE2ETest(
		TestBitcoinDepositAndCallRevertName,
		"deposit Bitcoin into ZEVM; expect refund", []runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.1"},
		},
		TestBitcoinDepositAndCallRevert,
	),
	runner.NewE2ETest(
		TestBitcoinDepositAndCallRevertWithDustName,
		"deposit Bitcoin into ZEVM; revert with dust amount that aborts the CCTX", []runner.ArgDefinition{},
		TestBitcoinDepositAndCallRevertWithDust,
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
		TestBitcoinStdMemoInscribedDepositAndCallName,
		"deposit Bitcoin into ZEVM and call a contract with inscribed standard memo",
		[]runner.ArgDefinition{
			{Description: "amount in btc", DefaultValue: "0.1"},
			{Description: "fee rate", DefaultValue: "10"},
		},
		TestBitcoinStdMemoInscribedDepositAndCall,
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
	runner.NewE2ETest(
		TestStressSolanaDepositName,
		"stress test SOL deposit",
		[]runner.ArgDefinition{
			{Description: "amount in lamports", DefaultValue: "1200000"},
			{Description: "count of SOL deposits", DefaultValue: "50"},
		},
		TestStressSolanaDeposit,
	),
	runner.NewE2ETest(
		TestStressSPLDepositName,
		"stress test SPL deposit",
		[]runner.ArgDefinition{
			{Description: "amount in SPL tokens", DefaultValue: "1200000"},
			{Description: "count of SPL deposits", DefaultValue: "50"},
		},
		TestStressSPLDeposit,
	),
	runner.NewE2ETest(
		TestStressSolanaWithdrawName,
		"stress test SOL withdrawals",
		[]runner.ArgDefinition{
			{Description: "amount in lamports", DefaultValue: "1000000"},
			{Description: "count of SOL withdrawals", DefaultValue: "50"},
		},
		TestStressSolanaWithdraw,
	),
	runner.NewE2ETest(
		TestStressSPLWithdrawName,
		"stress test SPL withdrawals",
		[]runner.ArgDefinition{
			{Description: "amount in SPL tokens", DefaultValue: "1000000"},
			{Description: "count of SPL withdrawals", DefaultValue: "50"},
		},
		TestStressSPLWithdraw,
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
		TestPrecompilesStakingIsDisabled,
	),
	runner.NewE2ETest(
		TestPrecompilesStakingThroughContractName,
		"test stateful precompiled contracts staking through contract",
		[]runner.ArgDefinition{},
		TestPrecompilesStakingThroughContract,
	),
	runner.NewE2ETest(
		TestPrecompilesBankName,
		"test stateful precompiled contracts bank with ZRC20 tokens",
		[]runner.ArgDefinition{},
		TestPrecompilesBank,
	),
	runner.NewE2ETest(
		TestPrecompilesBankFailName,
		"test stateful precompiled contracts bank with non ZRC20 tokens",
		[]runner.ArgDefinition{},
		TestPrecompilesBankNonZRC20,
	),
	runner.NewE2ETest(
		TestPrecompilesBankThroughContractName,
		"test stateful precompiled contracts bank through contract",
		[]runner.ArgDefinition{},
		TestPrecompilesBankThroughContract,
	),
	runner.NewE2ETest(
		TestPrecompilesDistributeName,
		"test stateful precompiled contracts distribute",
		[]runner.ArgDefinition{},
		TestPrecompilesDistributeAndClaim,
	),
	runner.NewE2ETest(
		TestPrecompilesDistributeNonZRC20Name,
		"test stateful precompiled contracts distribute with non ZRC20 tokens",
		[]runner.ArgDefinition{},
		TestPrecompilesDistributeNonZRC20,
	),
	runner.NewE2ETest(
		TestPrecompilesDistributeThroughContractName,
		"test stateful precompiled contracts distribute through contract",
		[]runner.ArgDefinition{},
		TestPrecompilesDistributeAndClaimThroughContract,
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
			{Description: "amount in wei", DefaultValue: "100000"},
		},
		TestEtherWithdrawRestricted,
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
			{Description: "count", DefaultValue: "3"},
		},
		legacy.TestMultipleERC20Deposit,
	),
	runner.NewE2ETest(
		TestLegacyMultipleERC20WithdrawsName,
		"withdraw ERC20 from ZEVM in multiple withdrawals (v1 protocol contracts)",
		[]runner.ArgDefinition{
			{Description: "amount", DefaultValue: "100"},
			{Description: "count", DefaultValue: "3"},
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
}
