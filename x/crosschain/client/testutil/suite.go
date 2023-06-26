//go:build TESTNET
// +build TESTNET

package testutil

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethcfg "github.com/evmos/ethermint/cmd/config"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/zeta-chain/zetacore/app"
	cmdcfg "github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/network"
	"github.com/zeta-chain/zetacore/x/crosschain/client/cli"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerCli "github.com/zeta-chain/zetacore/x/observer/client/cli"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"os"
	"strconv"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) Setconfig() {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	ethcfg.SetBip44CoinType(config)
	// Make sure address is compatible with ethereum
	config.SetAddressVerifier(app.VerifyAddressFormat)
	config.Seal()
}
func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")
	s.Setconfig()
	minOBsDel, ok := sdk.NewIntFromString("100000000000000000000")
	s.Require().True(ok)
	s.cfg.StakingTokens = minOBsDel.Mul(sdk.NewInt(int64(10)))
	s.cfg.BondedTokens = minOBsDel
	genesisState := s.cfg.GenesisState

	// Cross-chain genesis state
	var crossChainGenesis types.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[types.ModuleName], &crossChainGenesis))
	crossChainGenesis.Params.Enabled = true
	crossChainGenesisBz, err := s.cfg.Codec.MarshalJSON(&crossChainGenesis)
	s.Require().NoError(err)

	// Staking genesis state
	var stakingGenesisState stakingtypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenesisState))
	stakingGenesisState.Params.BondDenom = cmdcfg.BaseDenom
	stakingGenesisStateBz, err := s.cfg.Codec.MarshalJSON(&stakingGenesisState)
	s.Require().NoError(err)

	// Observer genesis state
	var observerGenesis observerTypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[observerTypes.ModuleName], &observerGenesis))
	var observerMapper []*observerTypes.ObserverMapper
	observerList := []string{"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax", "zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2"}
	for _, chain := range common.DefaultChainsList() {
		observerMapper = append(observerMapper, &observerTypes.ObserverMapper{
			ObserverChain: chain,
			ObserverList:  observerList,
		})
	}
	observerGenesis.Observers = observerMapper
	observerGenesisBz, err := s.cfg.Codec.MarshalJSON(&observerGenesis)
	s.Require().NoError(err)

	genesisState[types.ModuleName] = crossChainGenesisBz
	genesisState[stakingtypes.ModuleName] = stakingGenesisStateBz
	genesisState[observerTypes.ModuleName] = observerGenesisBz
	s.cfg.GenesisState = genesisState

	s.network, err = network.New(s.T(), app.NodeDir, s.cfg)
	s.Assert().NoError(err)
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestCCTXInboundVoter() {

	val1 := s.network.Validators[1]
	val := s.network.Validators[0]

	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, authcli.GetAccountCmd(), []string{val1.Address.String()})
	fmt.Println("out", out.String())

	fmt.Println("val.Address", val1.Address.String())

	cmd := cli.CmdCCTXInboundVoter()
	args := []string{
		"0x96B05C238b99768F349135de0653b687f9c13fEE",
		strconv.FormatInt(common.GoerliChain().ChainId, 10),
		"0x3b9Fe88DE29efD13240829A0c18E9EC7A44C3CA7",
		"0x96B05C238b99768F349135de0653b687f9c13fEE",
		strconv.FormatInt(common.GoerliChain().ChainId, 10),
		"10000000000000000000",
		"",
		"0x19398991572a825894b34b904ac1e3692720895351466b5c9e6bb7ae1e21d680",
		"100",
		"Gas",
		"",
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val1.Address),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=true", flags.FlagGenerateOnly),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	out, err = clitestutil.ExecTestCLICmd(val1.ClientCtx, cmd, args)
	s.Require().NoError(err)
	fmt.Println("unsigned TX :", out.String())
	unsignerdTx := WriteToNewTempFile(s.T(), out.String())
	res, err := TxSignExec(val1.ClientCtx, val1.Address, unsignerdTx.Name(), "--offline", "--account-number", "1", "--sequence", "1")
	s.Require().NoError(err)
	fmt.Println("signed TX :", out.String())
	signerdTx := WriteToNewTempFile(s.T(), res.String())
	out, err = clitestutil.ExecTestCLICmd(s.network.Validators[0].ClientCtx, authcli.GetBroadcastCommand(), []string{signerdTx.Name(), "--broadcast-mode", "block"})

	s.Require().NoError(err)
	fmt.Println(out.String())
	out, err = clitestutil.ExecTestCLICmd(val.ClientCtx, observerCli.CmdBallotByIdentifier(), []string{"0x583634d8d39952de71b9564a186d1aaa9576da0e0980174ea9c556109b098ddd"})
	s.Require().NoError(err)
	fmt.Println(out.String())

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
