package testutil

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	observerCli "github.com/zeta-chain/zetacore/x/observer/client/cli"
)

func MsgAddObserverExec(clientCtx client.Context, chain, obsType fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{chain.String(), obsType.String()}
	args = append(args, extraArgs...)
	return clitestutil.ExecTestCLICmd(clientCtx, observerCli.CmdAddObserver(), args)
}
