package ibccrosschain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"

	"github.com/zeta-chain/node/x/ibccrosschain/keeper"
)

var (
	_ porttypes.IBCModule = IBCModule{}
)

// IBCModule implements the ICS26 interface for transfer given the transfer keeper.
type IBCModule struct {
	keeper keeper.Keeper
}

// NewIBCModule creates a new IBCModule given the keeper
func NewIBCModule(k keeper.Keeper) IBCModule {
	return IBCModule{
		keeper: k,
	}
}

// OnChanOpenInit implements the IBCModule interface
// TODO: Implement the function as middleware
// https://github.com/zeta-chain/node/issues/2163
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	// Set variable to blank to remove lint error so we can keep variable name for future use
	_ = ctx
	_ = order
	_ = connectionHops
	_ = portID
	_ = channelID
	_ = counterparty

	return version, nil
}

// OnChanOpenTry implements the IBCModule interface.
// TODO: Implement the function as middleware
// https://github.com/zeta-chain/node/issues/2163
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	// Set variable to blank to remove lint error so we can keep variable name for future us
	_ = ctx
	_ = order
	_ = connectionHops
	_ = portID
	_ = channelID
	_ = counterparty
	_ = counterpartyVersion

	return "", nil
}

// OnChanOpenAck implements the IBCModule interface
// TODO: Implement the function as middleware
// https://github.com/zeta-chain/node/issues/2163
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	// Set variable to blank to remove lint error so we can keep variable name for future use
	_ = ctx
	_ = portID
	_ = channelID
	_ = counterpartyChannelID
	_ = counterpartyVersion

	return nil
}

// OnChanOpenConfirm implements the IBCModule interface
// TODO: Implement the function as middleware
// https://github.com/zeta-chain/node/issues/2163
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Set variable to blank to remove lint error so we can keep variable name for future use
	_ = ctx
	_ = portID
	_ = channelID

	return nil
}

// OnChanCloseInit implements the IBCModule interface
// TODO: Implement the function as middleware
// https://github.com/zeta-chain/node/issues/2163
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Set variable to blank to remove lint error so we can keep variable name for future use
	_ = ctx
	_ = portID
	_ = channelID

	return nil
}

// OnChanCloseConfirm implements the IBCModule interface
// TODO: Implement the function as middleware
// https://github.com/zeta-chain/node/issues/2163
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Set variable to blank to remove lint error so we can keep variable name for future use
	_ = ctx
	_ = portID
	_ = channelID

	return nil
}

// OnRecvPacket implements the IBCModule interface. A successful acknowledgement
// is returned if the packet data is successfully decoded and the receive application
// logic returns without error
// TODO: Implement the function as middleware
// https://github.com/zeta-chain/node/issues/2163
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	channelVersion string,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	// Set variable to blank to remove lint error so we can keep variable name for future use
	_ = ctx
	_ = channelVersion
	_ = packet
	_ = relayer

	return nil
}

// OnAcknowledgementPacket implements the IBCModule interface
// TODO: Implement the function as middleware
// https://github.com/zeta-chain/node/issues/2163
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	channelVersion string,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	// Set variable to blank to remove lint error so we can keep variable name for future use
	_ = ctx
	_ = channelVersion
	_ = packet
	_ = acknowledgement
	_ = relayer

	return nil
}

// OnTimeoutPacket implements the IBCModule interface
// TODO: Implement the function as middleware
// https://github.com/zeta-chain/node/issues/2163
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	channelVersion string,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	// Set variable to blank to remove lint error so we can keep variable name for future use
	_ = ctx
	_ = channelVersion
	_ = packet
	_ = relayer

	return nil
}
