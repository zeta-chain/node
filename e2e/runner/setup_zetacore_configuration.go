package runner

import (
	"math/big"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	e2eutils "github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// UpdateEVMChainParams update the erc20 custody contract and gateway address in the chain params
// TODO: should be used for all protocol contracts including the ZETA connector
// https://github.com/zeta-chain/node/issues/3257
func (r *E2ERunner) UpdateEVMChainParams(testLegacy bool) {
	res, err := r.ObserverClient.GetChainParams(r.Ctx, &observertypes.QueryGetChainParamsRequest{})
	require.NoError(r, err)

	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	// find old chain params
	var (
		chainParams *observertypes.ChainParams
		found       bool
	)
	for _, cp := range res.ChainParams.ChainParams {
		if cp.ChainId == evmChainID.Int64() {
			chainParams = cp
			found = true
			break
		}
	}
	require.True(r, found, "Chain params not found for chain id %d", evmChainID)

	// update with the new ERC20 custody contract address
	chainParams.Erc20CustodyContractAddress = r.ERC20CustodyAddr.Hex()

	// update with the new gateway address
	chainParams.GatewayAddress = r.GatewayEVMAddr.Hex()

	//  update with the new connector address only if not running legacy tests
	// when running legacy tests the connector address is set by the LegacySetupEVM function
	if !testLegacy {
		chainParams.ConnectorContractAddress = r.ConnectorNativeAddr.Hex()
	}

	// update the chain params
	err = r.ZetaTxServer.UpdateChainParams(chainParams)
	require.NoError(r, err)
}

// EmissionsPoolFunding represents the amount of ZETA to fund the emissions pool with
// This is the same value as used originally on mainnet (20M ZETA)
var EmissionsPoolFunding = big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(2e7))

// SetTSSAddresses set TSS addresses from information queried from ZetaChain
func (r *E2ERunner) SetTSSAddresses() error {
	btcChainID, err := chains.GetBTCChainIDFromChainParams(r.BitcoinParams)
	if err != nil {
		return err
	}

	res := &observertypes.QueryGetTssAddressResponse{}
	for i := 0; ; i++ {
		res, err = r.ObserverClient.GetTssAddress(r.Ctx, &observertypes.QueryGetTssAddressRequest{
			BitcoinChainId: btcChainID,
		})
		if err != nil {
			if i%10 == 0 {
				r.Logger.Info("ObserverClient.TSS error %s", err.Error())
				r.Logger.Info("TSS not ready yet, waiting for TSS to be appear in zetacore network...")
			}
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	tssAddress := ethcommon.HexToAddress(res.Eth)

	btcTSSAddress, err := btcutil.DecodeAddress(res.Btc, r.BitcoinParams)
	require.NoError(r, err)

	r.TSSAddress = tssAddress
	r.BTCTSSAddress = btcTSSAddress
	r.SuiTSSAddress = res.Sui

	return nil
}

// EnableHeaderVerification enables the header verification for the given chain IDs
func (r *E2ERunner) EnableHeaderVerification(chainIDList []int64) error {
	r.Logger.Print("⚙️ enabling verification flags for block headers")

	return r.ZetaTxServer.EnableHeaderVerification(e2eutils.AdminPolicyName, chainIDList)
}

// EnableV2ZETAFlows enables the V2 ZETA flows flag
func (r *E2ERunner) EnableV2ZETAFlows() error {
	r.Logger.Print("⚙️ enabling V2 ZETA flows")

	return r.ZetaTxServer.UpdateV2ZETAFlows(e2eutils.OperationalPolicyName, true)
}

// IsV2ZETAEnabled returns true if V2 ZETA flows are enabled
func (r *E2ERunner) IsV2ZETAEnabled() bool {
	res, err := r.ObserverClient.CrosschainFlags(r.Ctx, &observertypes.QueryGetCrosschainFlagsRequest{})
	if err != nil {
		r.Logger.Print("⚠️ failed to query crosschain flags: %v", err)
		return false
	}
	return res.CrosschainFlags.IsV2ZetaEnabled
}
