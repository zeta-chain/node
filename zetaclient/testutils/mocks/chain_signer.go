package mocks

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/zetacore/pkg/chains"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/outboundprocessor"
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
