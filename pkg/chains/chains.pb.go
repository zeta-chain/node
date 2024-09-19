// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: zetachain/zetacore/pkg/chains/chains.proto

package chains

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// ReceiveStatus represents the status of an outbound
// TODO: Rename and move
// https://github.com/zeta-chain/node/issues/2257
type ReceiveStatus int32

const (
	// Created is used for inbounds
	ReceiveStatus_created ReceiveStatus = 0
	ReceiveStatus_success ReceiveStatus = 1
	ReceiveStatus_failed  ReceiveStatus = 2
)

var ReceiveStatus_name = map[int32]string{
	0: "created",
	1: "success",
	2: "failed",
}

var ReceiveStatus_value = map[string]int32{
	"created": 0,
	"success": 1,
	"failed":  2,
}

func (x ReceiveStatus) String() string {
	return proto.EnumName(ReceiveStatus_name, int32(x))
}

func (ReceiveStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{0}
}

// ChainName represents the name of the chain
// Deprecated(v19): replaced with Chain.Name as string
type ChainName int32

const (
	ChainName_empty              ChainName = 0
	ChainName_eth_mainnet        ChainName = 1
	ChainName_zeta_mainnet       ChainName = 2
	ChainName_btc_mainnet        ChainName = 3
	ChainName_polygon_mainnet    ChainName = 4
	ChainName_bsc_mainnet        ChainName = 5
	ChainName_goerli_testnet     ChainName = 6
	ChainName_mumbai_testnet     ChainName = 7
	ChainName_bsc_testnet        ChainName = 10
	ChainName_zeta_testnet       ChainName = 11
	ChainName_btc_testnet        ChainName = 12
	ChainName_sepolia_testnet    ChainName = 13
	ChainName_goerli_localnet    ChainName = 14
	ChainName_btc_regtest        ChainName = 15
	ChainName_amoy_testnet       ChainName = 16
	ChainName_optimism_mainnet   ChainName = 17
	ChainName_optimism_sepolia   ChainName = 18
	ChainName_base_mainnet       ChainName = 19
	ChainName_base_sepolia       ChainName = 20
	ChainName_solana_mainnet     ChainName = 21
	ChainName_solana_devnet      ChainName = 22
	ChainName_solana_localnet    ChainName = 23
	ChainName_btc_signet_testnet ChainName = 24
)

var ChainName_name = map[int32]string{
	0:  "empty",
	1:  "eth_mainnet",
	2:  "zeta_mainnet",
	3:  "btc_mainnet",
	4:  "polygon_mainnet",
	5:  "bsc_mainnet",
	6:  "goerli_testnet",
	7:  "mumbai_testnet",
	10: "bsc_testnet",
	11: "zeta_testnet",
	12: "btc_testnet",
	13: "sepolia_testnet",
	14: "goerli_localnet",
	15: "btc_regtest",
	16: "amoy_testnet",
	17: "optimism_mainnet",
	18: "optimism_sepolia",
	19: "base_mainnet",
	20: "base_sepolia",
	21: "solana_mainnet",
	22: "solana_devnet",
	23: "solana_localnet",
	24: "btc_signet_testnet",
}

var ChainName_value = map[string]int32{
	"empty":              0,
	"eth_mainnet":        1,
	"zeta_mainnet":       2,
	"btc_mainnet":        3,
	"polygon_mainnet":    4,
	"bsc_mainnet":        5,
	"goerli_testnet":     6,
	"mumbai_testnet":     7,
	"bsc_testnet":        10,
	"zeta_testnet":       11,
	"btc_testnet":        12,
	"sepolia_testnet":    13,
	"goerli_localnet":    14,
	"btc_regtest":        15,
	"amoy_testnet":       16,
	"optimism_mainnet":   17,
	"optimism_sepolia":   18,
	"base_mainnet":       19,
	"base_sepolia":       20,
	"solana_mainnet":     21,
	"solana_devnet":      22,
	"solana_localnet":    23,
	"btc_signet_testnet": 24,
}

func (x ChainName) String() string {
	return proto.EnumName(ChainName_name, int32(x))
}

func (ChainName) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{1}
}

// Network represents the network of the chain
// there is a single instance of the network on mainnet
// then the network can have eventual testnets or devnets
type Network int32

const (
	Network_eth      Network = 0
	Network_zeta     Network = 1
	Network_btc      Network = 2
	Network_polygon  Network = 3
	Network_bsc      Network = 4
	Network_optimism Network = 5
	Network_base     Network = 6
	Network_solana   Network = 7
)

var Network_name = map[int32]string{
	0: "eth",
	1: "zeta",
	2: "btc",
	3: "polygon",
	4: "bsc",
	5: "optimism",
	6: "base",
	7: "solana",
}

var Network_value = map[string]int32{
	"eth":      0,
	"zeta":     1,
	"btc":      2,
	"polygon":  3,
	"bsc":      4,
	"optimism": 5,
	"base":     6,
	"solana":   7,
}

func (x Network) String() string {
	return proto.EnumName(Network_name, int32(x))
}

func (Network) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{2}
}

// NetworkType represents the network type of the chain
// Mainnet, Testnet, Privnet, Devnet
type NetworkType int32

const (
	NetworkType_mainnet NetworkType = 0
	NetworkType_testnet NetworkType = 1
	NetworkType_privnet NetworkType = 2
	NetworkType_devnet  NetworkType = 3
)

var NetworkType_name = map[int32]string{
	0: "mainnet",
	1: "testnet",
	2: "privnet",
	3: "devnet",
}

var NetworkType_value = map[string]int32{
	"mainnet": 0,
	"testnet": 1,
	"privnet": 2,
	"devnet":  3,
}

func (x NetworkType) String() string {
	return proto.EnumName(NetworkType_name, int32(x))
}

func (NetworkType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{3}
}

// Vm represents the virtual machine type of the chain to support smart
// contracts
type Vm int32

const (
	Vm_no_vm Vm = 0
	Vm_evm   Vm = 1
	Vm_svm   Vm = 2
)

var Vm_name = map[int32]string{
	0: "no_vm",
	1: "evm",
	2: "svm",
}

var Vm_value = map[string]int32{
	"no_vm": 0,
	"evm":   1,
	"svm":   2,
}

func (x Vm) String() string {
	return proto.EnumName(Vm_name, int32(x))
}

func (Vm) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{4}
}

// Consensus represents the consensus algorithm used by the chain
// this can represent the consensus of a L1
// this can also represent the solution of a L2
type Consensus int32

const (
	Consensus_ethereum         Consensus = 0
	Consensus_tendermint       Consensus = 1
	Consensus_bitcoin          Consensus = 2
	Consensus_op_stack         Consensus = 3
	Consensus_solana_consensus Consensus = 4
)

var Consensus_name = map[int32]string{
	0: "ethereum",
	1: "tendermint",
	2: "bitcoin",
	3: "op_stack",
	4: "solana_consensus",
}

var Consensus_value = map[string]int32{
	"ethereum":         0,
	"tendermint":       1,
	"bitcoin":          2,
	"op_stack":         3,
	"solana_consensus": 4,
}

func (x Consensus) String() string {
	return proto.EnumName(Consensus_name, int32(x))
}

func (Consensus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{5}
}

// CCTXGateway describes for the chain the gateway used to handle CCTX outbounds
type CCTXGateway int32

const (
	// zevm is the internal CCTX gateway to process outbound on the ZEVM and read
	// inbound events from the ZEVM only used for ZetaChain chains
	CCTXGateway_zevm CCTXGateway = 0
	// observers is the CCTX gateway for chains relying on the observer set to
	// observe inbounds and TSS for outbounds
	CCTXGateway_observers CCTXGateway = 1
)

var CCTXGateway_name = map[int32]string{
	0: "zevm",
	1: "observers",
}

var CCTXGateway_value = map[string]int32{
	"zevm":      0,
	"observers": 1,
}

func (x CCTXGateway) String() string {
	return proto.EnumName(CCTXGateway_name, int32(x))
}

func (CCTXGateway) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{6}
}

// Chain represents static data about a blockchain network
// it is identified by a unique chain ID
type Chain struct {
	// ChainId is the unique identifier of the chain
	ChainId int64 `protobuf:"varint,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// ChainName is the name of the chain
	// Deprecated(v19): replaced with Name
	ChainName ChainName `protobuf:"varint,1,opt,name=chain_name,json=chainName,proto3,enum=zetachain.zetacore.pkg.chains.ChainName" json:"chain_name,omitempty"` // Deprecated: Do not use.
	// Network is the network of the chain
	Network Network `protobuf:"varint,3,opt,name=network,proto3,enum=zetachain.zetacore.pkg.chains.Network" json:"network,omitempty"`
	// NetworkType is the network type of the chain: mainnet, testnet, etc..
	NetworkType NetworkType `protobuf:"varint,4,opt,name=network_type,json=networkType,proto3,enum=zetachain.zetacore.pkg.chains.NetworkType" json:"network_type,omitempty"`
	// Vm is the virtual machine used in the chain
	Vm Vm `protobuf:"varint,5,opt,name=vm,proto3,enum=zetachain.zetacore.pkg.chains.Vm" json:"vm,omitempty"`
	// Consensus is the underlying consensus algorithm used by the chain
	Consensus Consensus `protobuf:"varint,6,opt,name=consensus,proto3,enum=zetachain.zetacore.pkg.chains.Consensus" json:"consensus,omitempty"`
	// IsExternal describe if the chain is ZetaChain or external
	IsExternal bool `protobuf:"varint,7,opt,name=is_external,json=isExternal,proto3" json:"is_external,omitempty"`
	// CCTXGateway is the gateway used to handle CCTX outbounds
	CctxGateway CCTXGateway `protobuf:"varint,8,opt,name=cctx_gateway,json=cctxGateway,proto3,enum=zetachain.zetacore.pkg.chains.CCTXGateway" json:"cctx_gateway,omitempty"`
	// Name is the name of the chain
	Name string `protobuf:"bytes,9,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *Chain) Reset()         { *m = Chain{} }
func (m *Chain) String() string { return proto.CompactTextString(m) }
func (*Chain) ProtoMessage()    {}
func (*Chain) Descriptor() ([]byte, []int) {
	return fileDescriptor_236b85e7bff6130d, []int{0}
}
func (m *Chain) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Chain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Chain.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Chain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Chain.Merge(m, src)
}
func (m *Chain) XXX_Size() int {
	return m.Size()
}
func (m *Chain) XXX_DiscardUnknown() {
	xxx_messageInfo_Chain.DiscardUnknown(m)
}

var xxx_messageInfo_Chain proto.InternalMessageInfo

func (m *Chain) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

// Deprecated: Do not use.
func (m *Chain) GetChainName() ChainName {
	if m != nil {
		return m.ChainName
	}
	return ChainName_empty
}

func (m *Chain) GetNetwork() Network {
	if m != nil {
		return m.Network
	}
	return Network_eth
}

func (m *Chain) GetNetworkType() NetworkType {
	if m != nil {
		return m.NetworkType
	}
	return NetworkType_mainnet
}

func (m *Chain) GetVm() Vm {
	if m != nil {
		return m.Vm
	}
	return Vm_no_vm
}

func (m *Chain) GetConsensus() Consensus {
	if m != nil {
		return m.Consensus
	}
	return Consensus_ethereum
}

func (m *Chain) GetIsExternal() bool {
	if m != nil {
		return m.IsExternal
	}
	return false
}

func (m *Chain) GetCctxGateway() CCTXGateway {
	if m != nil {
		return m.CctxGateway
	}
	return CCTXGateway_zevm
}

func (m *Chain) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.ReceiveStatus", ReceiveStatus_name, ReceiveStatus_value)
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.ChainName", ChainName_name, ChainName_value)
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.Network", Network_name, Network_value)
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.NetworkType", NetworkType_name, NetworkType_value)
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.Vm", Vm_name, Vm_value)
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.Consensus", Consensus_name, Consensus_value)
	proto.RegisterEnum("zetachain.zetacore.pkg.chains.CCTXGateway", CCTXGateway_name, CCTXGateway_value)
	proto.RegisterType((*Chain)(nil), "zetachain.zetacore.pkg.chains.Chain")
}

func init() {
	proto.RegisterFile("zetachain/zetacore/pkg/chains/chains.proto", fileDescriptor_236b85e7bff6130d)
}

var fileDescriptor_236b85e7bff6130d = []byte{
	// 781 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xcd, 0x8e, 0xe4, 0x34,
	0x10, 0xee, 0x24, 0xfd, 0x5b, 0x3d, 0x3f, 0x5e, 0xef, 0xb0, 0x84, 0x95, 0x68, 0x06, 0x24, 0xa0,
	0x35, 0x82, 0x1e, 0x01, 0x47, 0x0e, 0xa0, 0x1d, 0xb1, 0x08, 0x21, 0xf6, 0x10, 0x56, 0x2b, 0xc4,
	0xa5, 0x71, 0xbb, 0x8b, 0xb4, 0xd5, 0xb1, 0x1d, 0xc5, 0xee, 0xec, 0x36, 0x4f, 0xc1, 0x43, 0x70,
	0xe0, 0x3d, 0xb8, 0x70, 0xdc, 0x23, 0x47, 0x34, 0xf3, 0x20, 0x20, 0x3b, 0x4e, 0x7a, 0xb8, 0x30,
	0x73, 0x8a, 0xfd, 0xe5, 0xfb, 0xaa, 0xbe, 0xb2, 0xab, 0x0c, 0x17, 0xbf, 0xa0, 0x65, 0x7c, 0xc3,
	0x84, 0xba, 0xf4, 0x2b, 0x5d, 0xe1, 0x65, 0xb9, 0xcd, 0x2f, 0x3d, 0x64, 0xc2, 0x67, 0x51, 0x56,
	0xda, 0x6a, 0xfa, 0x76, 0xc7, 0x5d, 0xb4, 0xdc, 0x45, 0xb9, 0xcd, 0x17, 0x0d, 0xe9, 0xf1, 0x59,
	0xae, 0x73, 0xed, 0x99, 0x97, 0x6e, 0xd5, 0x88, 0xde, 0xfb, 0x27, 0x81, 0xc1, 0x95, 0x23, 0xd0,
	0xb7, 0x60, 0xec, 0x99, 0x4b, 0xb1, 0x4e, 0xe3, 0xf3, 0x68, 0x9e, 0x64, 0x23, 0xbf, 0xff, 0x66,
	0x4d, 0xbf, 0x05, 0x68, 0x7e, 0x29, 0x26, 0x31, 0x8d, 0xce, 0xa3, 0xf9, 0xc9, 0xa7, 0xf3, 0xc5,
	0xff, 0xa6, 0x5b, 0xf8, 0xa0, 0xcf, 0x98, 0xc4, 0x27, 0x71, 0x1a, 0x65, 0x13, 0xde, 0x6e, 0xe9,
	0x97, 0x30, 0x52, 0x68, 0x5f, 0xea, 0x6a, 0x9b, 0x26, 0x3e, 0xd2, 0x07, 0x77, 0x44, 0x7a, 0xd6,
	0xb0, 0xb3, 0x56, 0x46, 0xbf, 0x83, 0xa3, 0xb0, 0x5c, 0xda, 0x7d, 0x89, 0x69, 0xdf, 0x87, 0xb9,
	0xb8, 0x5f, 0x98, 0xe7, 0xfb, 0x12, 0xb3, 0xa9, 0x3a, 0x6c, 0xe8, 0x27, 0x10, 0xd7, 0x32, 0x1d,
	0xf8, 0x20, 0xef, 0xde, 0x11, 0xe4, 0x85, 0xcc, 0xe2, 0x5a, 0xd2, 0xa7, 0x30, 0xe1, 0x5a, 0x19,
	0x54, 0x66, 0x67, 0xd2, 0xe1, 0xfd, 0xce, 0xa3, 0xe5, 0x67, 0x07, 0x29, 0x7d, 0x07, 0xa6, 0xc2,
	0x2c, 0xf1, 0x95, 0xc5, 0x4a, 0xb1, 0x22, 0x1d, 0x9d, 0x47, 0xf3, 0x71, 0x06, 0xc2, 0x7c, 0x15,
	0x10, 0x57, 0x2a, 0xe7, 0xf6, 0xd5, 0x32, 0x67, 0x16, 0x5f, 0xb2, 0x7d, 0x3a, 0xbe, 0x57, 0xa9,
	0x57, 0x57, 0xcf, 0x7f, 0xf8, 0xba, 0x51, 0x64, 0x53, 0xa7, 0x0f, 0x1b, 0x4a, 0xa1, 0xef, 0xaf,
	0x70, 0x72, 0x1e, 0xcd, 0x27, 0x99, 0x5f, 0x5f, 0x7c, 0x0e, 0xc7, 0x19, 0x72, 0x14, 0x35, 0x7e,
	0x6f, 0x99, 0xdd, 0x19, 0x3a, 0x85, 0x11, 0xaf, 0x90, 0x59, 0x5c, 0x93, 0x9e, 0xdb, 0x98, 0x1d,
	0xe7, 0x68, 0x0c, 0x89, 0x28, 0xc0, 0xf0, 0x67, 0x26, 0x0a, 0x5c, 0x93, 0xf8, 0x71, 0xff, 0xf7,
	0xdf, 0x66, 0xd1, 0xc5, 0x1f, 0x09, 0x4c, 0xba, 0x9b, 0xa6, 0x13, 0x18, 0xa0, 0x2c, 0xed, 0x9e,
	0xf4, 0xe8, 0x29, 0x4c, 0xd1, 0x6e, 0x96, 0x92, 0x09, 0xa5, 0xd0, 0x92, 0x88, 0x12, 0x38, 0x72,
	0x56, 0x3b, 0x24, 0x76, 0x94, 0x95, 0xe5, 0x1d, 0x90, 0xd0, 0x87, 0x70, 0x5a, 0xea, 0x62, 0x9f,
	0x6b, 0xd5, 0x81, 0x7d, 0xcf, 0x32, 0x07, 0xd6, 0x80, 0x52, 0x38, 0xc9, 0x35, 0x56, 0x85, 0x58,
	0x5a, 0x34, 0xd6, 0x61, 0x43, 0x87, 0xc9, 0x9d, 0x5c, 0xb1, 0x03, 0x36, 0x6a, 0x85, 0x2d, 0x00,
	0x9d, 0x83, 0x16, 0x99, 0xb6, 0x0e, 0x5a, 0xe0, 0xc8, 0x39, 0x30, 0x58, 0xea, 0x42, 0x1c, 0x58,
	0xc7, 0x0e, 0x0c, 0x09, 0x0b, 0xcd, 0x59, 0xe1, 0xc0, 0x93, 0x56, 0x5a, 0x61, 0xee, 0x88, 0xe4,
	0xd4, 0x45, 0x67, 0x52, 0xef, 0x3b, 0x1d, 0xa1, 0x67, 0x40, 0x74, 0x69, 0x85, 0x14, 0x46, 0x76,
	0xf6, 0x1f, 0xfc, 0x07, 0x0d, 0xb9, 0x08, 0x75, 0xea, 0x15, 0x33, 0xd8, 0xf1, 0x1e, 0x76, 0x48,
	0xcb, 0x39, 0x73, 0x45, 0x1a, 0x5d, 0x30, 0x75, 0x38, 0xc3, 0x37, 0xe8, 0x03, 0x38, 0x0e, 0xd8,
	0x1a, 0x6b, 0x07, 0x3d, 0xf2, 0x35, 0x34, 0x50, 0x67, 0xf7, 0x4d, 0xfa, 0x08, 0xa8, 0xb3, 0x6b,
	0x44, 0xae, 0xd0, 0x76, 0x1e, 0xd3, 0x70, 0x8b, 0x08, 0xa3, 0x30, 0x1d, 0x74, 0x04, 0x09, 0xda,
	0x0d, 0xe9, 0xd1, 0x31, 0xf4, 0xdd, 0x69, 0x91, 0xc8, 0x41, 0x2b, 0xcb, 0x49, 0xec, 0x7a, 0x21,
	0xdc, 0x0f, 0x49, 0x3c, 0x6a, 0x38, 0xe9, 0xd3, 0x23, 0x18, 0xb7, 0x05, 0x91, 0x81, 0x93, 0x39,
	0xdb, 0x64, 0xe8, 0x9a, 0xa5, 0xf1, 0x41, 0x46, 0x21, 0xcd, 0x53, 0x98, 0xde, 0x1a, 0x42, 0x17,
	0xae, 0x2d, 0xc4, 0xf7, 0x59, 0xeb, 0x2a, 0xf2, 0x89, 0x2a, 0x51, 0x37, 0x6d, 0x02, 0x30, 0x0c,
	0xb5, 0x25, 0x21, 0xce, 0x87, 0x10, 0xbf, 0x90, 0xae, 0xd9, 0x94, 0x5e, 0xd6, 0x92, 0xf4, 0xbc,
	0xe9, 0x5a, 0x36, 0x56, 0x4d, 0x2d, 0xbb, 0xee, 0xfc, 0x09, 0x26, 0xdd, 0xd8, 0x39, 0x9f, 0x68,
	0x37, 0x58, 0xe1, 0xce, 0x49, 0x4e, 0x00, 0x2c, 0xaa, 0x35, 0x56, 0x52, 0xa8, 0x90, 0x72, 0x25,
	0x2c, 0xd7, 0x42, 0x91, 0xb8, 0x29, 0x69, 0x69, 0x2c, 0xe3, 0x5b, 0x92, 0xb8, 0x1b, 0x0b, 0x07,
	0xda, 0x0d, 0x2e, 0xe9, 0x87, 0x0c, 0x1f, 0xc1, 0xf4, 0xd6, 0xb0, 0x35, 0x87, 0xe6, 0x2d, 0x1d,
	0xc3, 0x44, 0xaf, 0x0c, 0x56, 0x35, 0x56, 0x86, 0x44, 0x0d, 0xfb, 0xc9, 0x17, 0x7f, 0x5e, 0xcf,
	0xa2, 0xd7, 0xd7, 0xb3, 0xe8, 0xef, 0xeb, 0x59, 0xf4, 0xeb, 0xcd, 0xac, 0xf7, 0xfa, 0x66, 0xd6,
	0xfb, 0xeb, 0x66, 0xd6, 0xfb, 0xf1, 0xfd, 0x5c, 0xd8, 0xcd, 0x6e, 0xb5, 0xe0, 0x5a, 0xfa, 0x87,
	0xfe, 0xe3, 0xe6, 0xcd, 0x57, 0x7a, 0x7d, 0xfb, 0xbd, 0x5f, 0x0d, 0xfd, 0xa3, 0xfd, 0xd9, 0xbf,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x84, 0x13, 0xff, 0x06, 0x17, 0x06, 0x00, 0x00,
}

func (m *Chain) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Chain) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Chain) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintChains(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x4a
	}
	if m.CctxGateway != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.CctxGateway))
		i--
		dAtA[i] = 0x40
	}
	if m.IsExternal {
		i--
		if m.IsExternal {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x38
	}
	if m.Consensus != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.Consensus))
		i--
		dAtA[i] = 0x30
	}
	if m.Vm != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.Vm))
		i--
		dAtA[i] = 0x28
	}
	if m.NetworkType != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.NetworkType))
		i--
		dAtA[i] = 0x20
	}
	if m.Network != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.Network))
		i--
		dAtA[i] = 0x18
	}
	if m.ChainId != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x10
	}
	if m.ChainName != 0 {
		i = encodeVarintChains(dAtA, i, uint64(m.ChainName))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintChains(dAtA []byte, offset int, v uint64) int {
	offset -= sovChains(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Chain) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ChainName != 0 {
		n += 1 + sovChains(uint64(m.ChainName))
	}
	if m.ChainId != 0 {
		n += 1 + sovChains(uint64(m.ChainId))
	}
	if m.Network != 0 {
		n += 1 + sovChains(uint64(m.Network))
	}
	if m.NetworkType != 0 {
		n += 1 + sovChains(uint64(m.NetworkType))
	}
	if m.Vm != 0 {
		n += 1 + sovChains(uint64(m.Vm))
	}
	if m.Consensus != 0 {
		n += 1 + sovChains(uint64(m.Consensus))
	}
	if m.IsExternal {
		n += 2
	}
	if m.CctxGateway != 0 {
		n += 1 + sovChains(uint64(m.CctxGateway))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovChains(uint64(l))
	}
	return n
}

func sovChains(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozChains(x uint64) (n int) {
	return sovChains(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Chain) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowChains
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Chain: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Chain: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainName", wireType)
			}
			m.ChainName = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainName |= ChainName(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainId |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Network", wireType)
			}
			m.Network = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Network |= Network(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NetworkType", wireType)
			}
			m.NetworkType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NetworkType |= NetworkType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vm", wireType)
			}
			m.Vm = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Vm |= Vm(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Consensus", wireType)
			}
			m.Consensus = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Consensus |= Consensus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsExternal", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsExternal = bool(v != 0)
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CctxGateway", wireType)
			}
			m.CctxGateway = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CctxGateway |= CCTXGateway(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowChains
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthChains
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthChains
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipChains(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthChains
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipChains(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowChains
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowChains
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowChains
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthChains
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupChains
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthChains
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthChains        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowChains          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupChains = fmt.Errorf("proto: unexpected end of group")
)
