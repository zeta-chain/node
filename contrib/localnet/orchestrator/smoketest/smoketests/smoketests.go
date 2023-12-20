package smoketests

import "github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"

// AllSmokeTests is an ordered list of all smoke tests
var AllSmokeTests = []runner.SmokeTest{
	{
		"context_upgrade",
		"tests sending ETH on ZEVM and check context data using ContextApp",
		TestContextUpgrade,
	},
	{
		"deposit_and_call_refund",
		"deposit ZRC20 into ZEVM and call a contract that reverts; should refund",
		TestDepositAndCallRefund,
	},
	{
		"erc20_multiple_deposit",
		"deposit USDT ERC20 into ZEVM",
		TestMultipleERC20Deposit,
	},
	{
		"erc20_withdraw",
		"withdraw USDT ERC20 from ZEVM",
		TestWithdrawERC20,
	},
	{
		"erc20_multiple_withdraw",
		"withdraw USDT ERC20 from ZEVM in multiple deposits",
		TestMultipleWithdraws,
	},
	{
		"send_zeta_out",
		"sending ZETA from ZEVM to Ethereum",
		TestSendZetaOut,
	},
	{
		"send_zeta_out_btc_revert",
		"sending ZETA from ZEVM to Bitcoin; should revert when ",
		TestSendZetaOutBTCRevert,
	},
	{
		"message_passing",
		"goerli->goerli message passing (sending ZETA only)",
		TestMessagePassing,
	},
	{
		"zrc20_swap",
		"swap ZRC20 USDT for ZRC20 ETH",
		TestZRC20Swap,
	},
	{
		"bitcoin_withdraw",
		"withdraw BTC from ZEVM",
		TestBitcoinWithdraw,
	},
	{
		"crosschain_swap",
		"testing Bitcoin ERC20 cross-chain swap",
		TestCrosschainSwap,
	},
	{
		"message_passing_revert_fail",
		"goerli->goerli message passing (revert fail)",
		TestMessagePassingRevertFail,
	},
	{
		"message_passing_revert_success",
		"goerli->goerli message passing (revert success)",
		TestMessagePassingRevertSuccess,
	},
	{
		"pause_zrc20",
		"pausing ZRC20 on ZetaChain",
		TestPauseZRC20,
	},
	{
		"erc20_deposit_and_call_refund",
		"deposit a non-gas ZRC20 into ZEVM and call a contract that reverts; should refund on ZetaChain if no liquidity pool, should refund on origin if liquidity pool",
		TestERC20DepositAndCallRefund,
	},
	{
		"update_bytecode",
		"update ZRC20 bytecode swap",
		TestUpdateBytecode,
	},
	{
		"eth_deposit_and_call",
		"deposit ZRC20 into ZEVM and call a contract",
		TestEtherDepositAndCall,
	},
	{
		"deposit_eth_liquidity_cap",
		"deposit Ethers into ZEVM with a liquidity cap",
		TestDepositEtherLiquidityCap,
	},
	{
		"block_headers",
		"fetch block headers of EVM on ZetaChain",
		TestBlockHeaders,
	},
	{
		"whitelist_erc20",
		"whitelist ERC20",
		TestWhitelistERC20,
	},
	{
		"my_test",
		"performing custom test",
		TestMyTest,
	},
}

// FindSmokeTest finds a smoke test by name
func FindSmokeTest(name string) (runner.SmokeTest, bool) {
	for _, test := range AllSmokeTests {
		if test.Name == name {
			return test, true
		}
	}
	return runner.SmokeTest{}, false
}
