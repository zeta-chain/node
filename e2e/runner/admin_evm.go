package runner

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
)

func (r *E2ERunner) UpdateTssAddressForConnector() {
	require.NoError(r, r.SetTSSAddresses())

	tx, err := r.ConnectorEth.UpdateTssAddress(r.EVMAuth, r.TSSAddress)
	require.NoError(r, err)
	r.Logger.Info(fmt.Sprintf("TSS Address Update Tx: %s", tx.Hash().String()))
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tssAddressOnConnector, err := r.ConnectorEth.TssAddress(&bind.CallOpts{Context: r.Ctx})
	require.NoError(r, err)
	require.Equal(r, r.TSSAddress, tssAddressOnConnector)
}

func (r *E2ERunner) UpdateTssAddressForErc20custody() {
	require.NoError(r, r.SetTSSAddresses())

	tx, err := r.ERC20Custody.UpdateTSSAddress(r.EVMAuth, r.TSSAddress)
	require.NoError(r, err)
	r.Logger.Info(fmt.Sprintf("TSS ERC20 Address Update Tx: %s", tx.Hash().String()))
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tssAddressOnCustody, err := r.ERC20Custody.TSSAddress(&bind.CallOpts{Context: r.Ctx})
	require.NoError(r, err)
	require.Equal(r, r.TSSAddress, tssAddressOnCustody)
}
