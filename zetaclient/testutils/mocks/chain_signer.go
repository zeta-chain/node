package mocks

import (
	"context"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/pkg/chains"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

type DummySigner struct{}

func (s *DummySigner) TryProcessOutbound(
	_ context.Context,
	_ *crosschaintypes.CrossChainTx,
	_ interfaces.ChainObserver,
	_ interfaces.ZetacoreClient,
	_ uint64,
) {
}

func (s *DummySigner) SetGatewayAddress(_ string)                     {}
func (s *DummySigner) GetGatewayAddress() (_ string)                  { return }
func (s *DummySigner) SetZetaConnectorAddress(_ ethcommon.Address)    {}
func (s *DummySigner) SetERC20CustodyAddress(_ ethcommon.Address)     {}
func (s *DummySigner) GetZetaConnectorAddress() (_ ethcommon.Address) { return }
func (s *DummySigner) GetERC20CustodyAddress() (_ ethcommon.Address)  { return }

// ----------------------------------------------------------------------------
// EVMSigner
// ----------------------------------------------------------------------------
var _ interfaces.ChainSigner = (*EVMSigner)(nil)

// EVMSigner is a mock of evm chain signer for testing
type EVMSigner struct {
	DummySigner
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
type BTCSigner = DummySigner

func NewBTCSigner() *BTCSigner {
	return &BTCSigner{}
}

// ----------------------------------------------------------------------------
// SolanaSigner
// ----------------------------------------------------------------------------
var _ interfaces.ChainSigner = (*SolanaSigner)(nil)

// SolanaSigner is a mock of solana chain signer for testing
type SolanaSigner struct {
	DummySigner
	GatewayAddress string
}

func NewSolanaSigner() *SolanaSigner {
	return &SolanaSigner{}
}

func (s *SolanaSigner) SetGatewayAddress(address string) {
	s.GatewayAddress = address
}

func (s *SolanaSigner) GetGatewayAddress() string {
	return s.GatewayAddress
}

// ----------------------------------------------------------------------------
// TONSigner
// ----------------------------------------------------------------------------
var _ interfaces.ChainSigner = (*TONSigner)(nil)

// TONSigner is a mock of TON chain signer for testing
type TONSigner struct {
	DummySigner
	GatewayAddress string
}

func NewTONSigner() *TONSigner {
	return &TONSigner{}
}

func (s *TONSigner) SetGatewayAddress(address string) {
	s.GatewayAddress = address
}

func (s *TONSigner) GetGatewayAddress() string {
	return s.GatewayAddress
}
