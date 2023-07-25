package querytests

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	emmisonscli "github.com/zeta-chain/zetacore/x/emissions/client/cli"
	emmisonstypes "github.com/zeta-chain/zetacore/x/emissions/types"
)

func (s *CliTestSuite) TestObserverRewards() {
	val := s.network.Validators[0]
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.GetBalancesCmd(), []string{val.Address.String(), "--output", "json"})
	s.Require().NoError(err)
	fmt.Println(out.String())

	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, emmisonscli.CmdListPoolAddresses(), []string{"--output", "json"})
	s.Require().NoError(err)
	resPools := emmisonstypes.QueryListPoolAddressesResponse{}
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resPools))
	txArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin("azeta", sdk.NewInt(10))).String()),
	}
	sendArgs := []string{val.Address.String(),
		resPools.EmissionModuleAddress, "800000000000000000000azeta"}
	args := append(sendArgs, txArgs...)
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewSendTxCmd(), args)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, emmisonscli.CmdGetEmmisonsFactors(), []string{"--output", "json"})
	resFactors := emmisonstypes.QueryGetEmmisonsFactorsResponse{}
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resFactors))
	fmt.Println(resFactors)

	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, emmisonscli.CmdGetEmmisonsFactors(), []string{"--output", "json"})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &resFactors))
	fmt.Println(resFactors)
}
