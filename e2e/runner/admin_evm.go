package runner

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
)

func (r *E2ERunner) UpdateTSSAddressForConnector() {
	require.NoError(r, r.SetTSSAddresses())

	tx, err := r.ConnectorEth.UpdateTssAddress(r.EVMAuth, r.TSSAddress)
	require.NoError(r, err)
	r.Logger.Info("TSS Address Update Tx: %s", tx.Hash().String())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tssAddressOnConnector, err := r.ConnectorEth.TssAddress(&bind.CallOpts{Context: r.Ctx})
	require.NoError(r, err)
	require.Equal(r, r.TSSAddress, tssAddressOnConnector)
}

func (r *E2ERunner) UpdateTSSAddressForERC20custody() {
	require.NoError(r, r.SetTSSAddresses())

	tx, err := r.ERC20CustodyV2.UpdateTSSAddress(r.EVMAuth, r.TSSAddress)
	require.NoError(r, err)
	r.Logger.Info("TSS ERC20 Address Update Tx: %s", tx.Hash().String())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// we have to reference ERC20CustodyV2 since it's `TssAddress` on v2 and `TSSAddress` on v1
	tssAddressOnCustody, err := r.ERC20CustodyV2.TssAddress(&bind.CallOpts{Context: r.Ctx})
	require.NoError(r, err)
	require.Equal(r, r.TSSAddress, tssAddressOnCustody)
}
