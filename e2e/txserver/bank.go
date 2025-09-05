package txserver

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	e2eutils "github.com/zeta-chain/node/e2e/utils"
)

// TransferZETA transfers ZETA tokens from the admin account to an address
// Note: admin account is only used because it already have a initialized balance
func (zts *ZetaTxServer) TransferZETA(destination sdk.AccAddress, amount int64) error {
	// retrieve account
	accAdmin, err := zts.clientCtx.Keyring.Key(e2eutils.AdminPolicyName)
	if err != nil {
		return err
	}
	addrAdmin, err := accAdmin.GetAddress()
	if err != nil {
		return err
	}

	// add new authorization
	_, err = zts.BroadcastTx(e2eutils.AdminPolicyName, banktypes.NewMsgSend(
		addrAdmin,
		destination,
		sdk.NewCoins(sdk.NewInt64Coin(config.BaseDenom, amount)),
	))
	if err != nil {
		return fmt.Errorf("failed to add authorization: %w", err)
	}

	return nil
}
