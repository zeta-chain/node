package smoketests

import "github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"

// AllSmokeTests is an ordered list of all smoke tests
var AllSmokeTests = []runner.SmokeTest{
	TestContextUpgrade,
	TestDepositAndCallRefund,
	TestERC20Deposit,
	TestERC20Withdraw,
	TestSendZetaOut,
	TestSendZetaOutBTCRevert,
	TestMessagePassing,
	TestZRC20Swap,
	TestBitcoinWithdraw,
	TestCrosschainSwap,
	TestMessagePassingRevertFail,
	TestMessagePassingRevertSuccess,
	TestPauseZRC20,
	TestERC20DepositAndCallRefund,
	TestUpdateBytecode,
	TestEtherDepositAndCall,
	TestDepositEtherLiquidityCap,
	TestBlockHeaders,
	TestWhitelistERC20,
	TestMyTest,
}
