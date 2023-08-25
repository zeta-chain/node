//go:build TESTNET
// +build TESTNET

package integrationtests

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/network"
	"github.com/zeta-chain/zetacore/x/crosschain/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/zeta-chain/zetacore/x/crosschain/client/cli"
)

func MsgVoteOnObservedInboundTxExec(clientCtx client.Context, chain, obsType fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{chain.String(), obsType.String()}
	args = append(args, extraArgs...)
	return clitestutil.ExecTestCLICmd(clientCtx, cli.CmdCCTXInboundVoter(), args)
}

func TxSignExec(clientCtx client.Context, from fmt.Stringer, filename string, extraArgs ...string) (testutil.BufferWriter, error) {
	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, keyring.BackendTest),
		fmt.Sprintf("--from=%s", from.String()),
		fmt.Sprintf("--%s=%s", flags.FlagChainID, clientCtx.ChainID),
		filename,
	}

	cmd := authcli.GetSignCommand()
	tmcli.PrepareBaseCmd(cmd, "", "")

	return clitestutil.ExecTestCLICmd(clientCtx, cmd, append(args, extraArgs...))
}

func WriteToNewTempFile(t testing.TB, s string) *os.File {
	t.Helper()

	fp := TempFile(t)
	_, err := fp.WriteString(s)

	require.Nil(t, err)

	return fp
}

// TempFile returns a writable temporary file for the test to use.
func TempFile(t testing.TB) *os.File {
	t.Helper()

	fp, err := os.CreateTemp(GetTempDir(t), "")
	require.NoError(t, err)

	return fp
}

// GetTempDir returns a writable temporary director for the test to use.
func GetTempDir(t testing.TB) string {
	t.Helper()
	// os.MkDir() is used instead of testing.T.TempDir()
	// see https://github.com/cosmos/cosmos-sdk/pull/8475 and
	// https://github.com/cosmos/cosmos-sdk/pull/10341 for
	// this change's rationale.
	tempdir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.RemoveAll(tempdir) })
	return tempdir
}

func BuildSignedGasPriceVote(t testing.TB, val *network.Validator, denom string, account authtypes.AccountI) *os.File {
	cmd := cli.CmdGasPriceVoter()
	inboundVoterArgs := []string{
		strconv.FormatInt(common.GoerliChain().ChainId, 10),
		"10000000000",
		"100",
		"100",
	}
	txArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(10))).String()),
	}
	args := append(inboundVoterArgs, txArgs...)
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, args)
	require.NoError(t, err)
	unsignerdTx := WriteToNewTempFile(t, out.String())
	res, err := TxSignExec(val.ClientCtx, val.Address, unsignerdTx.Name(),
		"--offline", "--account-number", strconv.FormatUint(account.GetAccountNumber(), 10), "--sequence", strconv.FormatUint(account.GetSequence(), 10))
	require.NoError(t, err)
	return WriteToNewTempFile(t, res.String())
}

func BuildSignedTssVote(t testing.TB, val *network.Validator, denom string, account authtypes.AccountI) *os.File {
	cmd := cli.CmdCreateTSSVoter()
	inboundVoterArgs := []string{
		"tsspubkey",
		strconv.FormatInt(common.GoerliChain().ChainId, 10),
		"0",
	}
	txArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(10))).String()),
	}
	args := append(inboundVoterArgs, txArgs...)
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, args)
	require.NoError(t, err)
	unsignerdTx := WriteToNewTempFile(t, out.String())
	res, err := TxSignExec(val.ClientCtx, val.Address, unsignerdTx.Name(),
		"--offline", "--account-number", strconv.FormatUint(account.GetAccountNumber(), 10), "--sequence", strconv.FormatUint(account.GetSequence(), 10))
	require.NoError(t, err)
	return WriteToNewTempFile(t, res.String())
}

func BuildSignedOutboundVote(t testing.TB, val *network.Validator, denom string, account authtypes.AccountI, cctxIndex, outTxHash, zetaminted, status string) *os.File {
	cmd := cli.CmdCCTXOutboundVoter()
	outboundVoterArgs := []string{
		cctxIndex,
		outTxHash,
		"1",
		"0",
		zetaminted,
		status,
		strconv.FormatInt(common.GoerliChain().ChainId, 10),
		"1",
		"Gas",
	}
	txArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(10))).String()),
	}
	args := append(outboundVoterArgs, txArgs...)
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, args)
	require.NoError(t, err)

	unsignerdTx := WriteToNewTempFile(t, out.String())
	res, err := TxSignExec(val.ClientCtx, val.Address, unsignerdTx.Name(),
		"--offline", "--account-number", strconv.FormatUint(account.GetAccountNumber(), 10), "--sequence", strconv.FormatUint(account.GetSequence(), 10))
	require.NoError(t, err)
	return WriteToNewTempFile(t, res.String())
}

func BuildSignedInboundVote(t testing.TB, val *network.Validator, denom string, account authtypes.AccountI, message string) *os.File {
	cmd := cli.CmdCCTXInboundVoter()
	inboundVoterArgs := []string{
		"0x96B05C238b99768F349135de0653b687f9c13fEE",
		strconv.FormatInt(common.GoerliChain().ChainId, 10),
		"0x3b9Fe88DE29efD13240829A0c18E9EC7A44C3CA7",
		"0x96B05C238b99768F349135de0653b687f9c13fEE",
		strconv.FormatInt(common.GoerliChain().ChainId, 10),
		"10000000000000000000",
		message,
		"0x19398991572a825894b34b904ac1e3692720895351466b5c9e6bb7ae1e21d680",
		"100",
		"Gas",
		"",
	}
	txArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(10))).String()),
	}
	args := append(inboundVoterArgs, txArgs...)
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, args)
	require.NoError(t, err)
	unsignerdTx := WriteToNewTempFile(t, out.String())
	res, err := TxSignExec(val.ClientCtx, val.Address, unsignerdTx.Name(),
		"--offline", "--account-number", strconv.FormatUint(account.GetAccountNumber(), 10), "--sequence", strconv.FormatUint(account.GetSequence(), 10))
	require.NoError(t, err)
	return WriteToNewTempFile(t, res.String())
}

func GetBallotIdentifier(message string) string {
	msg := types.NewMsgVoteOnObservedInboundTx(
		"",
		"0x96B05C238b99768F349135de0653b687f9c13fEE",
		common.GoerliChain().ChainId,
		"0x3b9Fe88DE29efD13240829A0c18E9EC7A44C3CA7",
		"0x96B05C238b99768F349135de0653b687f9c13fEE",
		common.GoerliChain().ChainId,
		sdk.NewUint(10000000000000000000),
		message,
		"0x19398991572a825894b34b904ac1e3692720895351466b5c9e6bb7ae1e21d680",
		100,
		250_000,
		common.CoinType_Gas,
		"",
	)
	return msg.Digest()
}

func GetBallotIdentifierOutBound(cctxindex, outtxHash, zetaminted string) string {
	math.NewUintFromString(zetaminted)

	msg := types.NewMsgVoteOnObservedOutboundTx(
		"",
		cctxindex,
		outtxHash,
		1,
		0,
		math.NewUintFromString(zetaminted),
		0,
		common.GoerliChain().ChainId,
		1,
		common.CoinType_Gas,
	)
	return msg.Digest()
}
