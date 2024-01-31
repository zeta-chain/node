package smoketests

import "github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"

const (
	TestContextUpgradeName              = "context_upgrade"
	TestDepositAndCallRefundName        = "deposit_and_call_refund"
	TestMultipleERC20DepositName        = "erc20_multiple_deposit"
	TestMultipleWithdrawsName           = "erc20_multiple_withdraw"
	TestZetaWithdrawName                = "zeta_withdraw"
	TestZetaWithdrawBTCRevertName       = "zeta_withdraw_btc_revert" // #nosec G101 - not a hardcoded password
	TestMessagePassingName              = "message_passing"
	TestZRC20SwapName                   = "zrc20_swap"
	TestBitcoinWithdrawName             = "bitcoin_withdraw"
	TestCrosschainSwapName              = "crosschain_swap"
	TestMessagePassingRevertFailName    = "message_passing_revert_fail"
	TestMessagePassingRevertSuccessName = "message_passing_revert_success"
	TestPauseZRC20Name                  = "pause_zrc20"
	TestERC20DepositAndCallRefundName   = "erc20_deposit_and_call_refund"
	TestUpdateBytecodeName              = "update_bytecode"
	TestEtherDepositAndCallName         = "eth_deposit_and_call"
	TestDepositEtherLiquidityCapName    = "deposit_eth_liquidity_cap"
	TestMyTestName                      = "my_test"

	TestERC20WithdrawName  = "erc20_withdraw"
	TestERC20DepositName   = "erc20_deposit"
	TestEtherDepositName   = "eth_deposit"
	TestEtherWithdrawName  = "eth_withdraw"
	TestBitcoinDepositName = "bitcoin_deposit"
	TestZetaDepositName    = "zeta_deposit"

	TestDonationEtherName = "donation_ether"

	TestStressEtherWithdrawName = "stress_eth_withdraw"
	TestStressBTCWithdrawName   = "stress_btc_withdraw"
)

// AllSmokeTests is an ordered list of all smoke tests
var AllSmokeTests = []runner.SmokeTest{
	{
		TestContextUpgradeName,
		"tests sending ETH on ZEVM and check context data using ContextApp",
		TestContextUpgrade,
	},
	{
		TestDepositAndCallRefundName,
		"deposit ZRC20 into ZEVM and call a contract that reverts; should refund",
		TestDepositAndCallRefund,
	},
	{
		TestMultipleERC20DepositName,
		"deposit USDT ERC20 into ZEVM in multiple deposits",
		TestMultipleERC20Deposit,
	},
	{
		TestERC20WithdrawName,
		"withdraw ERC20 from ZEVM",
		TestERC20Withdraw,
	},
	{
		TestMultipleWithdrawsName,
		"withdraw ERC20 from ZEVM in multiple deposits",
		TestMultipleWithdraws,
	},
	{
		TestZetaWithdrawName,
		"withdraw ZETA from ZEVM to Ethereum",
		TestZetaWithdraw,
	},
	{
		TestZetaDepositName,
		"deposit ZETA from Ethereum to ZEVM",
		TestZetaDeposit,
	},
	{
		TestZetaWithdrawBTCRevertName,
		"sending ZETA from ZEVM to Bitcoin with a message that should revert cctxs",
		TestZetaWithdrawBTCRevert,
	},
	{
		TestMessagePassingName,
		"goerli->goerli message passing (sending ZETA only)",
		TestMessagePassing,
	},
	{
		TestZRC20SwapName,
		"swap ZRC20 USDT for ZRC20 ETH",
		TestZRC20Swap,
	},
	{
		TestBitcoinWithdrawName,
		"withdraw BTC from ZEVM",
		TestBitcoinWithdraw,
	},
	{
		TestCrosschainSwapName,
		"testing Bitcoin ERC20 cross-chain swap",
		TestCrosschainSwap,
	},
	{
		TestMessagePassingRevertFailName,
		"goerli->goerli message passing (revert fail)",
		TestMessagePassingRevertFail,
	},
	{
		TestMessagePassingRevertSuccessName,
		"goerli->goerli message passing (revert success)",
		TestMessagePassingRevertSuccess,
	},
	{
		TestPauseZRC20Name,
		"pausing ZRC20 on ZetaChain",
		TestPauseZRC20,
	},
	{
		TestERC20DepositAndCallRefundName,
		"deposit a non-gas ZRC20 into ZEVM and call a contract that reverts",
		TestERC20DepositAndCallRefund,
	},
	{
		TestUpdateBytecodeName,
		"update ZRC20 bytecode swap",
		TestUpdateBytecode,
	},
	{
		TestEtherDepositAndCallName,
		"deposit ZRC20 into ZEVM and call a contract",
		TestEtherDepositAndCall,
	},
	{
		TestDepositEtherLiquidityCapName,
		"deposit Ethers into ZEVM with a liquidity cap",
		TestDepositEtherLiquidityCap,
	},
	{
		TestMyTestName,
		"performing custom test",
		TestMyTest,
	},
	{
		TestERC20DepositName,
		"deposit ERC20 into ZEVM",
		TestERC20Deposit,
	},
	{
		TestEtherDepositName,
		"deposit Ether into ZEVM",
		TestEtherDeposit,
	},
	{
		TestEtherWithdrawName,
		"withdraw Ether from ZEVM",
		TestEtherWithdraw,
	},
	{
		TestBitcoinDepositName,
		"deposit Bitcoin into ZEVM",
		TestBitcoinDeposit,
	},
	{
		TestDonationEtherName,
		"donate Ether to the TSS",
		TestDonationEther,
	},
	{
		TestStressEtherWithdrawName,
		"stress test Ether withdrawal",
		TestStressEtherWithdraw,
	},
	{
		TestStressBTCWithdrawName,
		"stress test BTC withdrawal",
		TestStressBTCWithdraw,
	},
}
