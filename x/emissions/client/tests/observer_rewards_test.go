package querytests

import (
	"fmt"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	emmisonscli "github.com/zeta-chain/zetacore/x/emissions/client/cli"
)

func (s *CliTestSuite) TestObserverRewards() {
	val := s.network.Validators[0]
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cli.GetBalancesCmd(), []string{val.Address.String(), "--output", "json"})
	s.Require().NoError(err)
	fmt.Println(out.String())
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, emmisonscli.CmdQueryParams(), []string{"--output", "json"})
	s.Require().NoError(err)
	fmt.Println(out.String())

}
