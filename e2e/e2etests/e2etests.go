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
	runner.NewE2ETest(TestContextUpgradeName, "tests sending ETH on ZEVM and check context data using ContextApp", "amount (0.001ETH)", []string{"1000000000000000"}, TestContextUpgrade),
	runner.NewE2ETest(TestDepositAndCallRefundName, "deposit ZRC20 into ZEVM and call a contract that reverts; should refund", "amount (10ETH)", []string{"10000000000000000000"}, TestDepositAndCallRefund),
	runner.NewE2ETest(TestMultipleERC20DepositName, "deposit USDT ERC20 into ZEVM in multiple deposits", "amount (1e9), count (3)", []string{"1000000000", "3"}, TestMultipleERC20Deposit),
	runner.NewE2ETest(TestERC20WithdrawName, "withdraw ERC20 from ZEVM", "amount (1000)", []string{"1000"}, TestERC20Withdraw),
	runner.NewE2ETest(TestMultipleWithdrawsName, "withdraw ERC20 from ZEVM in multiple withdrawals", "amount (100), count (3)", []string{"100", "3"}, TestMultipleWithdraws),
	runner.NewE2ETest(TestZetaWithdrawName, "withdraw ZETA from ZEVM to Ethereum", "amount (10zeta)", []string{"10000000000000000000"}, TestZetaWithdraw),
	runner.NewE2ETest(TestZetaDepositName, "deposit ZETA from Ethereum to ZEVM", "amount (1zeta)", []string{"1000000000000000000"}, TestZetaDeposit),
	runner.NewE2ETest(TestZetaWithdrawBTCRevertName, "sending ZETA from ZEVM to Bitcoin with a message that should revert cctxs", "amount (1zeta)", []string{"1000000000000000000"}, TestZetaWithdrawBTCRevert),
	runner.NewE2ETest(TestMessagePassingName, "goerli->goerli message passing (sending ZETA only)", "amount (10zeta)", []string{"10000000000000000000"}, TestMessagePassing),
	runner.NewE2ETest(TestZRC20SwapName, "swap ZRC20 USDT for ZRC20 ETH", "", []string{}, TestZRC20Swap),
	runner.NewE2ETest(TestBitcoinWithdrawName, "withdraw BTC from ZEVM", "amount (0.1btc)", []string{"0.1"}, TestBitcoinWithdraw),
	runner.NewE2ETest(TestCrosschainSwapName, "testing Bitcoin ERC20 cross-chain swap", "", []string{}, TestCrosschainSwap),
	runner.NewE2ETest(TestMessagePassingRevertFailName, "goerli->goerli message passing (revert fail)", "amount (10zeta)", []string{"10000000000000000000"}, TestMessagePassingRevertFail),
	runner.NewE2ETest(TestMessagePassingRevertSuccessName, "goerli->goerli message passing (revert success)", "amount (10zeta)", []string{"10000000000000000000"}, TestMessagePassingRevertSuccess),
	runner.NewE2ETest(TestPauseZRC20Name, "pausing ZRC20 on ZetaChain", "", []string{}, TestPauseZRC20),
	runner.NewE2ETest(TestERC20DepositAndCallRefundName, "deposit a non-gas ZRC20 into ZEVM and call a contract that reverts", "", []string{}, TestERC20DepositAndCallRefund),
	runner.NewE2ETest(TestUpdateBytecodeName, "update ZRC20 bytecode swap", "", []string{}, TestUpdateBytecode),
	runner.NewE2ETest(TestEtherDepositAndCallName, "deposit ZRC20 into ZEVM and call a contract", "amount (1ETH)", []string{"1000000000000000000"}, TestEtherDepositAndCall),
	runner.NewE2ETest(TestDepositEtherLiquidityCapName, "deposit Ethers into ZEVM with a liquidity cap", "liquidity cap (0.01ETH)", []string{"100000000000000"}, TestDepositEtherLiquidityCap),
	runner.NewE2ETest(TestMyTestName, "performing custom test", "", []string{}, TestMyTest),
	runner.NewE2ETest(TestERC20DepositName, "deposit ERC20 into ZEVM", "amount (100000)", []string{"100000"}, TestERC20Deposit),
	runner.NewE2ETest(TestEtherDepositName, "deposit Ether into ZEVM", "amount (0.01ETH)", []string{"10000000000000000"}, TestEtherDeposit),
	runner.NewE2ETest(TestEtherWithdrawName, "withdraw Ether from ZEVM", "amount (100000wei)", []string{"100000"}, TestEtherWithdraw),
	runner.NewE2ETest(TestBitcoinDepositName, "deposit Bitcoin into ZEVM", "amount (0.001btc)", []string{"0.001"}, TestBitcoinDeposit),
	runner.NewE2ETest(TestDonationEtherName, "donate Ether to the TSS", "amount (0.1ETH)", []string{"100000000000000000"}, TestDonationEther),
	runner.NewE2ETest(TestStressEtherWithdrawName, "stress test Ether withdrawal", "amount (100000wei), count (100)", []string{"100000", "100"}, TestStressEtherWithdraw),
	runner.NewE2ETest(TestStressBTCWithdrawName, "stress test BTC withdrawal", "amount (0.01), count (100)", []string{}, TestStressBTCWithdraw),
	runner.NewE2ETest(TestStressEtherDepositName, "stress test Ether deposit", "amount (100000wei), count (100)", []string{"100000", "100"}, TestStressEtherDeposit),
	runner.NewE2ETest(TestStressBTCDepositName, "stress test BTC deposit", "amount (0.001btc), count (100)", []string{"0.001", "100"}, TestStressBTCDeposit),
}
