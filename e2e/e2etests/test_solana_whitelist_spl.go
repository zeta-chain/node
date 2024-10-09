package e2etests

import (
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSolanaWhitelistSPL(r *runner.E2ERunner, args []string) {
	// Deploy a new SPL
	r.Logger.Info("Deploying new SPL")

	// load deployer private key
	privkey, err := solana.PrivateKeyFromBase58(r.Account.SolanaPrivateKey.String())
	require.NoError(r, err)

	spl := r.DeploySPL(&privkey)

	// whitelist erc20 zrc20
	r.Logger.Info("whitelisting ERC20 on new network")
	res, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, crosschaintypes.NewMsgWhitelistERC20(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		spl.PublicKey().String(),
		chains.SolanaLocalnet.ChainId,
		"NEWSOLERC20",
		"NEWSOLERC20",
		6,
		100000,
	))
	require.NoError(r, err)

	// retrieve zrc20 and cctx from event
	whitelistCCTXIndex, err := txserver.FetchAttributeFromTxResponse(res, "whitelist_cctx_index")
	require.NoError(r, err)

	erc20zrc20Addr, err := txserver.FetchAttributeFromTxResponse(res, "zrc20_address")
	require.NoError(r, err)

	err = r.ZetaTxServer.InitializeLiquidityCap(erc20zrc20Addr)
	require.NoError(r, err)

	// ensure CCTX created
	resCCTX, err := r.CctxClient.Cctx(r.Ctx, &crosschaintypes.QueryGetCctxRequest{Index: whitelistCCTXIndex})
	require.NoError(r, err)

	cctx := resCCTX.CrossChainTx
	r.Logger.CCTX(*cctx, "whitelist_cctx")

	// wait for the whitelist cctx to be mined
	r.WaitForMinedCCTXFromIndex(whitelistCCTXIndex)
}
