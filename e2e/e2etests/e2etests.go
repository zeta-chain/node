package e2etests

import (
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// TODO : Add smoke test for abort refund
// https://github.com/zeta-chain/node/issues/1745
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
	TestStressEtherDepositName  = "stress_eth_deposit"
	TestStressBTCDepositName    = "stress_btc_deposit"
)

// AllE2ETests is an ordered list of all e2e tests
var AllE2ETests = []runner.E2ETest{
	{
		TestContextUpgradeName,
		[]string{},
		"tests sending ETH on ZEVM and check context data using ContextApp",
		"",
		TestContextUpgrade,
	},
	{
		TestDepositAndCallRefundName,
		[]string{},
		"deposit ZRC20 into ZEVM and call a contract that reverts; should refund",
		"",
		TestDepositAndCallRefund,
	},
	{
		TestMultipleERC20DepositName,
		[]string{},
		"deposit USDT ERC20 into ZEVM in multiple deposits",
		"",
		TestMultipleERC20Deposit,
	},
	{
		TestERC20WithdrawName,
		[]string{},
		"withdraw ERC20 from ZEVM",
		"",
		TestERC20Withdraw,
	},
	{
		TestMultipleWithdrawsName,
		[]string{},
		"withdraw ERC20 from ZEVM in multiple deposits",
		"",
		TestMultipleWithdraws,
	},
	{
		TestZetaWithdrawName,
		[]string{},
		"withdraw ZETA from ZEVM to Ethereum",
		"",
		TestZetaWithdraw,
	},
	{
		TestZetaDepositName,
		[]string{},
		"deposit ZETA from Ethereum to ZEVM",
		"",
		TestZetaDeposit,
	},
	{
		TestZetaWithdrawBTCRevertName,
		[]string{},
		"sending ZETA from ZEVM to Bitcoin with a message that should revert cctxs",
		"",
		TestZetaWithdrawBTCRevert,
	},
	{
		TestMessagePassingName,
		[]string{},
		"goerli->goerli message passing (sending ZETA only)",
		"",
		TestMessagePassing,
	},
	{
		TestZRC20SwapName,
		[]string{},
		"swap ZRC20 USDT for ZRC20 ETH",
		"",
		TestZRC20Swap,
	},
	{
		TestBitcoinWithdrawName,
		[]string{},
		"withdraw BTC from ZEVM",
		"",
		TestBitcoinWithdraw,
	},
	{
		TestCrosschainSwapName,
		[]string{},
		"testing Bitcoin ERC20 cross-chain swap",
		"",
		TestCrosschainSwap,
	},
	{
		TestMessagePassingRevertFailName,
		[]string{},
		"goerli->goerli message passing (revert fail)",
		"",
		TestMessagePassingRevertFail,
	},
	{
		TestMessagePassingRevertSuccessName,
		[]string{},
		"goerli->goerli message passing (revert success)",
		"",
		TestMessagePassingRevertSuccess,
	},
	{
		TestPauseZRC20Name,
		[]string{},
		"pausing ZRC20 on ZetaChain",
		"",
		TestPauseZRC20,
	},
	{
		TestERC20DepositAndCallRefundName,
		[]string{},
		"deposit a non-gas ZRC20 into ZEVM and call a contract that reverts",
		"",
		TestERC20DepositAndCallRefund,
	},
	{
		TestUpdateBytecodeName,
		[]string{},
		"update ZRC20 bytecode swap",
		"",
		TestUpdateBytecode,
	},
	{
		TestEtherDepositAndCallName,
		[]string{},
		"deposit ZRC20 into ZEVM and call a contract",
		"",
		TestEtherDepositAndCall,
	},
	{
		TestDepositEtherLiquidityCapName,
		[]string{},
		"deposit Ethers into ZEVM with a liquidity cap",
		"",
		TestDepositEtherLiquidityCap,
	},
	{
		TestMyTestName,
		[]string{},
		"performing custom test",
		"",
		TestMyTest,
	},
	{
		TestERC20DepositName,
		[]string{},
		"deposit ERC20 into ZEVM",
		"",
		TestERC20Deposit,
	},
	{
		TestEtherDepositName,
		[]string{},
		"deposit Ether into ZEVM",
		"amount in wei (default 0.01ETH)",
		TestEtherDeposit,
	},
	{
		TestEtherWithdrawName,
		[]string{},
		"withdraw Ether from ZEVM",
		"",
		TestEtherWithdraw,
	},
	{
		TestBitcoinDepositName,
		[]string{},
		"deposit Bitcoin into ZEVM",
		"",
		TestBitcoinDeposit,
	},
	{
		TestDonationEtherName,
		[]string{},
		"donate Ether to the TSS",
		"",
		TestDonationEther,
	},
	{
		TestStressEtherWithdrawName,
		[]string{},
		"stress test Ether withdrawal",
		"",
		TestStressEtherWithdraw,
	},
	{
		TestStressBTCWithdrawName,
		[]string{},
		"stress test BTC withdrawal",
		"",
		TestStressBTCWithdraw,
	},
	{
		TestStressEtherDepositName,
		[]string{},
		"stress test Ether deposit",
		"",
		TestStressEtherDeposit,
	},
	{
		TestStressBTCDepositName,
		[]string{},
		"stress test BTC deposit",
		"",
		TestStressBTCDeposit,
	},
}
