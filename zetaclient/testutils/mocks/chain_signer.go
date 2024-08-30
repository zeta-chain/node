package mocks

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/pkg/chains"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/outboundprocessor"
)

// ----------------------------------------------------------------------------
// EVMSigner
// ----------------------------------------------------------------------------
var _ interfaces.ChainSigner = (*EVMSigner)(nil)

// EVMSigner is a mock of evm chain signer for testing
type EVMSigner struct {
	Chain                chains.Chain
	ZetaConnectorAddress ethcommon.Address
	ERC20CustodyAddress  ethcommon.Address
}

func NewEVMSigner(
	chain chains.Chain,
	zetaConnectorAddress ethcommon.Address,
	erc20CustodyAddress ethcommon.Address,
) *EVMSigner {
	return &EVMSigner{
		Chain:                chain,
		ZetaConnectorAddress: zetaConnectorAddress,
		ERC20CustodyAddress:  erc20CustodyAddress,
	}
}

func (s *EVMSigner) TryProcessOutbound(
	_ context.Context,
	_ *crosschaintypes.CrossChainTx,
	_ *outboundprocessor.Processor,
	_ string,
	_ interfaces.ChainObserver,
	_ interfaces.ZetacoreClient,
	_ uint64,
) {
}

func (s *EVMSigner) SetGatewayAddress(_ string) {
}

func (s *EVMSigner) GetGatewayAddress() string {
	return ""
}

func (s *EVMSigner) SetZetaConnectorAddress(address ethcommon.Address) {
	s.ZetaConnectorAddress = address
}

func (s *EVMSigner) SetERC20CustodyAddress(address ethcommon.Address) {
	s.ERC20CustodyAddress = address
}

func (s *EVMSigner) GetZetaConnectorAddress() ethcommon.Address {
	return s.ZetaConnectorAddress
}

func (s *EVMSigner) GetERC20CustodyAddress() ethcommon.Address {
	return s.ERC20CustodyAddress
}

// ----------------------------------------------------------------------------
// BTCSigner
// ----------------------------------------------------------------------------
var _ interfaces.ChainSigner = (*BTCSigner)(nil)

// BTCSigner is a mock of bitcoin chain signer for testing
type BTCSigner struct {
}

func NewBTCSigner() *BTCSigner {
	return &BTCSigner{}
}

func (s *BTCSigner) TryProcessOutbound(
	_ context.Context,
	_ *crosschaintypes.CrossChainTx,
	_ *outboundprocessor.Processor,
	_ string,
	_ interfaces.ChainObserver,
	_ interfaces.ZetacoreClient,
	_ uint64,
) {
}

func (s *BTCSigner) SetGatewayAddress(_ string) {
}

func (s *BTCSigner) GetGatewayAddress() string {
	return ""
}

func (s *BTCSigner) SetZetaConnectorAddress(_ ethcommon.Address) {
}

func (s *BTCSigner) SetERC20CustodyAddress(_ ethcommon.Address) {
}

func (s *BTCSigner) GetZetaConnectorAddress() ethcommon.Address {
	return ethcommon.Address{}
}

func (s *BTCSigner) GetERC20CustodyAddress() ethcommon.Address {
	return ethcommon.Address{}
}

// ----------------------------------------------------------------------------
// SolanaSigner
// ----------------------------------------------------------------------------
var _ interfaces.ChainSigner = (*SolanaSigner)(nil)

// SolanaSigner is a mock of solana chain signer for testing
type SolanaSigner struct {
	GatewayAddress string
}

func NewSolanaSigner() *SolanaSigner {
	return &SolanaSigner{}
}

func (s *SolanaSigner) TryProcessOutbound(
	_ context.Context,
	_ *crosschaintypes.CrossChainTx,
	_ *outboundprocessor.Processor,
	_ string,
	_ interfaces.ChainObserver,
	_ interfaces.ZetacoreClient,
	_ uint64,
) {
}

func (s *SolanaSigner) SetGatewayAddress(address string) {
	s.GatewayAddress = address
}

func (s *SolanaSigner) GetGatewayAddress() string {
	return s.GatewayAddress
}

func (s *SolanaSigner) SetZetaConnectorAddress(_ ethcommon.Address) {
}

func (s *SolanaSigner) SetERC20CustodyAddress(_ ethcommon.Address) {
}

func (s *SolanaSigner) GetZetaConnectorAddress() ethcommon.Address {
	return ethcommon.Address{}
}

func (s *SolanaSigner) GetERC20CustodyAddress() ethcommon.Address {
	return ethcommon.Address{}
}
