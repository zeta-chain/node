package integrationtests

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/network"
	"github.com/zeta-chain/zetacore/x/crosschain/client/cli"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungiblecli "github.com/zeta-chain/zetacore/x/fungible/client/cli"
)

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
	t.Cleanup(func() {
		err := os.RemoveAll(tempdir)
		require.NoError(t, err)
	})
	return tempdir
}

func BuildSignedDeploySystemContract(t testing.TB, val *network.Validator, denom string, account authtypes.AccountI) *os.File {
	cmd := fungiblecli.CmdDeploySystemContracts()
	txArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100))).String()),
		// gas limit
		fmt.Sprintf("--%s=%d", flags.FlagGas, 4000000),
	}
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, txArgs)
	require.NoError(t, err)
	unsignerdTx := WriteToNewTempFile(t, out.String())
	res, err := TxSignExec(val.ClientCtx, val.Address, unsignerdTx.Name(),
		"--offline", "--account-number", strconv.FormatUint(account.GetAccountNumber(), 10), "--sequence", strconv.FormatUint(account.GetSequence(), 10))
	require.NoError(t, err)
	return WriteToNewTempFile(t, res.String())
}

func BuildSignedUpdateSystemContract(
	t testing.TB,
	val *network.Validator,
	denom string,
	account authtypes.AccountI,
	systemContractAddress string,
) *os.File {
	cmd := fungiblecli.CmdUpdateSystemContract()
	txArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100))).String()),
		// gas limit
		fmt.Sprintf("--%s=%d", flags.FlagGas, 4000000),
	}
	args := append([]string{systemContractAddress}, txArgs...)
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, args)
	require.NoError(t, err)
	unsignerdTx := WriteToNewTempFile(t, out.String())
	res, err := TxSignExec(val.ClientCtx, val.Address, unsignerdTx.Name(),
		"--offline", "--account-number", strconv.FormatUint(account.GetAccountNumber(), 10), "--sequence", strconv.FormatUint(account.GetSequence(), 10))
	require.NoError(t, err)
	return WriteToNewTempFile(t, res.String())
}

func BuildSignedDeployETHZRC20(
	t testing.TB,
	val *network.Validator,
	denom string,
	account authtypes.AccountI,
) *os.File {
	cmd := fungiblecli.CmdDeployFungibleCoinZRC4()
	txArgs := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100))).String()),
		// gas limit
		fmt.Sprintf("--%s=%d", flags.FlagGas, 10000000),
	}
	args := append([]string{
		"",
		strconv.FormatInt(common.GoerliLocalnetChain().ChainId, 10),
		"18",
		"ETH",
		"gETH",
		strconv.FormatInt(int64(common.CoinType_Gas), 10),
		"1000000",
	}, txArgs...)
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, args)
	require.NoError(t, err)
	unsignerdTx := WriteToNewTempFile(t, out.String())
	res, err := TxSignExec(val.ClientCtx, val.Address, unsignerdTx.Name(),
		"--offline", "--account-number", strconv.FormatUint(account.GetAccountNumber(), 10), "--sequence", strconv.FormatUint(account.GetSequence(), 10))
	require.NoError(t, err)
	return WriteToNewTempFile(t, res.String())
}

func BuildSignedGasPriceVote(t testing.TB, val *network.Validator, denom string, account authtypes.AccountI) *os.File {
	cmd := cli.CmdGasPriceVoter()
	inboundVoterArgs := []string{
		strconv.FormatInt(common.GoerliLocalnetChain().ChainId, 10),
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
		strconv.FormatInt(common.GoerliLocalnetChain().ChainId, 10),
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

func BuildSignedOutboundVote(
	t testing.TB,
	val *network.Validator,
	denom string,
	account authtypes.AccountI,
	nonce uint64,
	cctxIndex,
	outTxHash,
	valueReceived,
	status string,
) *os.File {
	cmd := cli.CmdCCTXOutboundVoter()
	outboundVoterArgs := []string{
		cctxIndex,
		outTxHash,
		"1",
		"0",
		"0",
		"0",
		valueReceived,
		status,
		strconv.FormatInt(common.GoerliLocalnetChain().ChainId, 10),
		strconv.FormatUint(nonce, 10),
		"Zeta",
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

func BuildSignedInboundVote(t testing.TB, val *network.Validator, denom string, account authtypes.AccountI, message string, eventIndex int) *os.File {
	cmd := cli.CmdCCTXInboundVoter()
	inboundVoterArgs := []string{
		"0x96B05C238b99768F349135de0653b687f9c13fEE",
		strconv.FormatInt(common.GoerliLocalnetChain().ChainId, 10),
		"0x3b9Fe88DE29efD13240829A0c18E9EC7A44C3CA7",
		"0x96B05C238b99768F349135de0653b687f9c13fEE",
		strconv.FormatInt(common.GoerliLocalnetChain().ChainId, 10),
		"10000000000000000000",
		message,
		"0x19398991572a825894b34b904ac1e3692720895351466b5c9e6bb7ae1e21d680",
		"100",
		"Zeta",
		"",
		strconv.Itoa(eventIndex),
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

func GetBallotIdentifier(message string, eventIndex int) string {
	msg := types.NewMsgVoteOnObservedInboundTx(
		"",
		"0x96B05C238b99768F349135de0653b687f9c13fEE",
		common.GoerliLocalnetChain().ChainId,
		"0x3b9Fe88DE29efD13240829A0c18E9EC7A44C3CA7",
		"0x96B05C238b99768F349135de0653b687f9c13fEE",
		common.GoerliLocalnetChain().ChainId,
		sdk.NewUint(10000000000000000000),
		message,
		"0x19398991572a825894b34b904ac1e3692720895351466b5c9e6bb7ae1e21d680",
		100,
		250_000,
		common.CoinType_Zeta,
		"",
		uint(eventIndex),
	)
	return msg.Digest()
}

func GetBallotIdentifierOutBound(nonce uint64, cctxindex, outtxHash, valueReceived string) string {
	msg := types.NewMsgVoteOnObservedOutboundTx(
		"",
		cctxindex,
		outtxHash,
		1,
		0,
		math.ZeroInt(),
		0,
		math.NewUintFromString(valueReceived),
		0,
		common.GoerliLocalnetChain().ChainId,
		nonce,
		common.CoinType_Zeta,
	)
	return msg.Digest()
}
