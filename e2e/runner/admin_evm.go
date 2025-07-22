package runner

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/utils"
)

// UpdateTSSAddressForConnector updates the TSS address for the connector contract
func (r *E2ERunner) UpdateTSSAddressForConnector() {
	require.NoError(r, r.SetTSSAddresses())
	require.NotNil(r, r.ConnectorEth, "Connector is not initialized")

	tx, err := r.ConnectorEth.UpdateTssAddress(r.EVMAuth, r.TSSAddress)
	require.NoError(r, err)
	r.Logger.Info("TSS Address Update Tx: %s", tx.Hash().String())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tssAddressOnConnector, err := r.ConnectorEth.TssAddress(&bind.CallOpts{Context: r.Ctx})
	require.NoError(r, err)
	require.Equal(r, r.TSSAddress, tssAddressOnConnector)
}

// UpdateTSSAddressForConnectorNative updates the TSS address for V2 connector contract
func (r *E2ERunner) UpdateTSSAddressForConnectorNative() {
	require.NoError(r, r.SetTSSAddresses())
	require.NotNil(r, r.ConnectorNative, "ConnectorNative is not initialized")

	tx, err := r.ConnectorNative.UpdateTSSAddress(r.EVMAuth, r.TSSAddress)
	require.NoError(r, err)
	r.Logger.Info("TSS Address Update Tx: %s", tx.Hash().String())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tssAddressOnConnectorNative, err := r.ConnectorNative.TssAddress(&bind.CallOpts{Context: r.Ctx})
	require.NoError(r, err)
	require.Equal(r, r.TSSAddress, tssAddressOnConnectorNative)
}

// UpdateTSSAddressForERC20custody updates the TSS address for the ERC20 custody contract
func (r *E2ERunner) UpdateTSSAddressForERC20custody() {
	require.NoError(r, r.SetTSSAddresses())

	tx, err := r.ERC20Custody.UpdateTSSAddress(r.EVMAuth, r.TSSAddress)
	require.NoError(r, err)
	r.Logger.Info("TSS ERC20 Address Update Tx: %s", tx.Hash().String())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// we have to reference ERC20Custody since it's `TssAddress` on v2 and `TSSAddress` on v1
	tssAddressOnCustody, err := r.ERC20Custody.TssAddress(&bind.CallOpts{Context: r.Ctx})
	require.NoError(r, err)
	require.Equal(r, r.TSSAddress, tssAddressOnCustody)
}

// UpdateTSSAddressForGateway updates the TSS address for the gateway contract
func (r *E2ERunner) UpdateTSSAddressForGateway() {
	require.NoError(r, r.SetTSSAddresses())

	tx, err := r.GatewayEVM.UpdateTSSAddress(r.EVMAuth, r.TSSAddress)
	require.NoError(r, err)
	r.Logger.Info("TSS Gateway Address Update Tx: %s", tx.Hash().String())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tssAddressOnGateway, err := r.GatewayEVM.TssAddress(&bind.CallOpts{Context: r.Ctx})
	require.NoError(r, err)
	require.Equal(r, r.TSSAddress, tssAddressOnGateway)
}
