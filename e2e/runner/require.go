package runner

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zrc20.sol"

	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	"github.com/zeta-chain/node/x/observer/types"
)

// EnsureNoTrackers ensures that there are no trackers left on zetacore
func (r *E2ERunner) EnsureNoTrackers() {
	// get all trackers
	res, err := r.CctxClient.OutTxTrackerAll(
		r.Ctx,
		&crosschaintypes.QueryAllOutboundTrackerRequest{},
	)
	require.NoError(r, err)
	require.Empty(r, res.OutboundTracker, "there should be no trackers at the end of the test")
}

// EnsureZeroBalanceAddressZEVM ensures that the balance of the restricted address is zero in the ZEVM
func (r *E2ERunner) EnsureZeroBalanceOnRestrictedAddressZEVM() {
	restrictedAddress := ethcommon.HexToAddress(sample.RestrictedEVMAddressTest)

	// ensure ZETA balance is zero
	balance, err := r.WZeta.BalanceOf(&bind.CallOpts{}, restrictedAddress)
	require.NoError(r, err)
	require.Zero(r, balance.Cmp(big.NewInt(0)), "the wZETA balance of the address should be zero")

	// ensure ZRC20 ETH balance is zero
	ensureZRC20ZeroBalance(r, r.ETHZRC20, restrictedAddress)

	// ensure ZRC20 ERC20 balance is zero
	ensureZRC20ZeroBalance(r, r.ERC20ZRC20, restrictedAddress)

	// ensure ZRC20 BTC balance is zero
	ensureZRC20ZeroBalance(r, r.BTCZRC20, restrictedAddress)

	// ensure ZRC20 SOL balance is zero
	ensureZRC20ZeroBalance(r, r.SOLZRC20, restrictedAddress)
}

// ensureZRC20ZeroBalance ensures that the balance of the ZRC20 token is zero on given address
func ensureZRC20ZeroBalance(r *E2ERunner, zrc20 *zrc20.ZRC20, address ethcommon.Address) {
	balance, err := zrc20.BalanceOf(&bind.CallOpts{}, address)
	require.NoError(r, err)

	zrc20Name, err := zrc20.Name(&bind.CallOpts{})
	require.NoError(r, err)
	require.Zero(
		r,
		balance.Cmp(big.NewInt(0)),
		fmt.Sprintf("the balance of address %s should be zero on ZRC20: %s", address, zrc20Name),
	)
}

// EnsureNoStaleBallots ensures that there are no stale ballots left on the chain.
func (r *E2ERunner) EnsureNoStaleBallots() {
	ballotsRes, err := r.ObserverClient.Ballots(r.Ctx, &types.QueryBallotsRequest{})
	require.NoError(r, err)
	currentBlockHeight, err := r.Clients.Zetacore.GetBlockHeight(r.Ctx)
	require.NoError(r, err)
	emissionsParams, err := r.EmissionsClient.Params(r.Ctx, &emissionstypes.QueryParamsRequest{})
	require.NoError(r, err)
	staleBlockStart := currentBlockHeight - emissionsParams.Params.BallotMaturityBlocks
	if len(ballotsRes.Ballots) < 1 {
		return
	}

	firstBallotCreationHeight := int64(0)

	for _, ballot := range ballotsRes.Ballots {
		if ballot.IsFinalized() {
			firstBallotCreationHeight = ballot.BallotCreationHeight
			break
		}
	}

	if firstBallotCreationHeight == 0 {
		return
	}
	// Log data for debugging
	if firstBallotCreationHeight < staleBlockStart {
		r.Logger.Error(
			"First finalized ballot creation height: %d is less than stale block start: %d",
			firstBallotCreationHeight,
			staleBlockStart,
		)
		for _, ballot := range ballotsRes.Ballots {
			r.Logger.Error(
				"Ballot: %s Creation Height %d BallotStatus %s\n",
				ballot.BallotIdentifier,
				ballot.BallotCreationHeight,
				ballot.BallotStatus,
			)
		}
		r.Logger.Error("Block Maturity Height: %d", emissionsParams.Params.BallotMaturityBlocks)
	}
	require.GreaterOrEqual(r, firstBallotCreationHeight, staleBlockStart, "there should be no stale ballots")
}
