package querytests

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcfg "github.com/evmos/ethermint/cmd/config"
	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/app"
	cmdcfg "github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/testutil/network"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"strconv"
)

type CliTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
	state   *types.GenesisState
}

func NewCLITestSuite(cfg network.Config) *CliTestSuite {
	return &CliTestSuite{cfg: cfg}
}

func (s *CliTestSuite) Setconfig() {
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	ethcfg.SetBip44CoinType(config)
	// Make sure address is compatible with ethereum
	config.SetAddressVerifier(app.VerifyAddressFormat)
	config.Seal()
}
func (s *CliTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")
	s.Setconfig()
	minOBsDel, ok := sdk.NewIntFromString("100000000000000000000")
	s.Require().True(ok)
	s.cfg.StakingTokens = minOBsDel.Mul(sdk.NewInt(int64(10)))
	s.cfg.BondedTokens = minOBsDel
	observerList := []string{"zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax",
		"zeta1f203dypqg5jh9hqfx0gfkmmnkdfuat3jr45ep2",
	}
	s.cfg.GenesisState = network.SetupZetaGenesisState(s.T(), s.cfg.GenesisState, s.cfg.Codec, observerList)
	s.AddCrossChainData(2)
	net, err := network.New(s.T(), app.NodeDir, s.cfg)
	s.Assert().NoError(err)
	s.network = net
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

}

func (s *CliTestSuite) TearDownSuite() {
	s.T().Log("tearing down genesis test suite")
	s.network.Cleanup()
}

func (s *CliTestSuite) AddCrossChainData(n int) {
	state := types.GenesisState{}
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(s.cfg.GenesisState[types.ModuleName], &state))
	// TODO : Fix add EVM balance to deploy contracts
	for i := 0; i < n; i++ {
		state.CrossChainTxs = append(state.CrossChainTxs, &types.CrossChainTx{
			Creator: "ANY",
			Index:   strconv.Itoa(i),
			CctxStatus: &types.Status{
				Status:              types.CctxStatus_PendingInbound,
				StatusMessage:       "",
				LastUpdateTimestamp: 0,
			},
			InboundTxParams:  &types.InboundTxParams{InboundTxObservedHash: fmt.Sprintf("Hash-%d", i), Amount: math.OneUint()},
			OutboundTxParams: []*types.OutboundTxParams{},
			ZetaFees:         math.OneUint()},
		)
	}
	for i := 0; i < n; i++ {
		state.ChainNoncesList = append(state.ChainNoncesList, &types.ChainNonces{Creator: "ANY", Index: strconv.Itoa(i), Signers: []string{}})
	}
	for i := 0; i < n; i++ {
		state.GasPriceList = append(state.GasPriceList, &types.GasPrice{Creator: "ANY", ChainId: int64(i), Index: strconv.Itoa(i), Prices: []uint64{}, BlockNums: []uint64{}, Signers: []string{}})
	}
	for i := 0; i < n; i++ {
		state.LastBlockHeightList = append(state.LastBlockHeightList, &types.LastBlockHeight{Creator: "ANY", Index: strconv.Itoa(i)})
	}
	state.Keygen = &types.Keygen{BlockNumber: 10, GranteePubkeys: []string{}}
	state.Tss = &types.TSS{
		TssPubkey:           "tssPubkey",
		TssParticipantList:  []string{"tssParticipantList"},
		OperatorAddressList: []string{"operatorAddressList"},
		FinalizedZetaHeight: 1,
		KeyGenZetaHeight:    1,
	}
	for i := 0; i < n; i++ {
		state.NodeAccountList = append(state.NodeAccountList, &types.NodeAccount{Operator: strconv.Itoa(i), GranteeAddress: "signer"})
	}
	for i := 0; i < n; i++ {
		outTxTracker := types.OutTxTracker{
			Index:   fmt.Sprintf("%d-%d", i, i),
			ChainId: int64(i),
			Nonce:   uint64(i),
		}
		nullify.Fill(&outTxTracker)
		state.OutTxTrackerList = append(state.OutTxTrackerList, outTxTracker)
	}

	for i := 0; i < n; i++ {
		inTxHashToCctx := types.InTxHashToCctx{
			InTxHash: strconv.Itoa(i),
		}
		nullify.Fill(&inTxHashToCctx)
		state.InTxHashToCctxList = append(state.InTxHashToCctxList, inTxHashToCctx)
	}
	permissionFlags := &types.PermissionFlags{}
	nullify.Fill(&permissionFlags)
	state.PermissionFlags = permissionFlags
	buf, err := s.cfg.Codec.MarshalJSON(&state)
	s.Require().NoError(err)
	s.cfg.GenesisState[types.ModuleName] = buf
	s.state = &state
}
