package chains

// ISC License
//
// Copyright (c) 2013-2024 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers

// this is a copy of the testnet4 parameters from https://github.com/btcsuite/btcd/pull/2275/
// they are not necessarily fully correct but should be sufficient for observation and signing

import (
	"math/big"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

// TestNet4 represents the test network (version 4).
const (
	TestNet4 wire.BitcoinNet = 0x283f161c
)

var (
	// bigOne is 1 represented as a big.Int. It is defined here to avoid
	// the overhead of creating it multiple times.
	bigOne = big.NewInt(1)
	// testNet3PowLimit is the highest proof of work value a Bitcoin block
	// can have for the test network (version 3). It is the value
	// 2^224 - 1.
	testNet3PowLimit = new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)
)

// testNet4GenesisCoinbaseTx is the coinbase transaction for the genesis block
// for the test network (version 4).
var testNet4GenesisCoinbaseTx = wire.MsgTx{
	Version: 1,
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash:  chainhash.Hash{},
				Index: 0xffffffff,
			},
			SignatureScript: []byte{
				0x04, 0xff, 0xff, 0x00, 0x1d, 0x01, 0x04, 0x4c, // |.......L|
				0x4c, 0x30, 0x33, 0x2f, 0x4d, 0x61, 0x79, 0x2f, // |L03/May/|
				0x32, 0x30, 0x32, 0x34, 0x20, 0x30, 0x30, 0x30, // |2024 000|
				0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, // |00000000|
				0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, // |00000000|
				0x30, 0x31, 0x65, 0x62, 0x64, 0x35, 0x38, 0x63, // |01ebd58c|
				0x32, 0x34, 0x34, 0x39, 0x37, 0x30, 0x62, 0x33, // |244970b3|
				0x61, 0x61, 0x39, 0x64, 0x37, 0x38, 0x33, 0x62, // |aa9d783b|
				0x62, 0x30, 0x30, 0x31, 0x30, 0x31, 0x31, 0x66, // |b001011f|
				0x62, 0x65, 0x38, 0x65, 0x61, 0x38, 0x65, 0x39, // |be8ea8e9|
				0x38, 0x65, 0x30, 0x30, 0x65, // |8e00e|
			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value: 0x12a05f200,
			PkScript: []byte{
				0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // |!.......|
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // |........|
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // |........|
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // |........|
				0x00, 0x00, 0xac, // |...|
			},
		},
	},
	LockTime: 0,
}

// testNet4GenesisHash is the hash of the first block in the block chain for the
// test network (version 4).
var testNet4GenesisHash = chainhash.Hash([chainhash.HashSize]byte{
	0x43, 0xf0, 0x8b, 0xda, 0xb0, 0x50, 0xe3, 0x5b,
	0x56, 0x7c, 0x86, 0x4b, 0x91, 0xf4, 0x7f, 0x50,
	0xae, 0x72, 0x5a, 0xe2, 0xde, 0x53, 0xbc, 0xfb,
	0xba, 0xf2, 0x84, 0xda, 0x00, 0x00, 0x00, 0x00,
})

// testNet4GenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the test network (version 4).
var testNet4GenesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{
	0x4e, 0x7b, 0x2b, 0x91, 0x28, 0xfe, 0x02, 0x91,
	0xdb, 0x06, 0x93, 0xaf, 0x2a, 0xe4, 0x18, 0xb7,
	0x67, 0xe6, 0x57, 0xcd, 0x40, 0x7e, 0x80, 0xcb,
	0x14, 0x34, 0x22, 0x1e, 0xae, 0xa7, 0xa0, 0x7a,
})

// testNet4GenesisBlock defines the genesis block of the block chain which
// serves as the public transaction ledger for the test network (version 4).
var testNet4GenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},          // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: testNet4GenesisMerkleRoot, // 7aa0a7ae1e223414cb807e40cd57e667b718e42aaf9306db9102fe28912b7b4e
		Timestamp:  time.Unix(1714777860, 0),  // 2024-05-03 23:11:00 +0000 UTC
		Bits:       0x1d00ffff,                // 486604799 [00000000ffff0000000000000000000000000000000000000000000000000000]
		Nonce:      0x17780cbb,                // 393743547
	},
	Transactions: []*wire.MsgTx{&testNet4GenesisCoinbaseTx},
}

// TestNet4Params defines the network parameters for the test Bitcoin network
// (version 4). Not to be confused with the regression test network, this
// network is sometimes simply called "testnet4".
var TestNet4Params = chaincfg.Params{
	Name:        "testnet4",
	Net:         TestNet4,
	DefaultPort: "48333",
	DNSSeeds: []chaincfg.DNSSeed{
		{"seed.testnet4.bitcoin.sprovoost.nl", true},
		{"seed.testnet4.wiz.biz", true},
	},

	// Chain parameters
	GenesisBlock:             &testNet4GenesisBlock,
	GenesisHash:              &testNet4GenesisHash,
	PowLimit:                 testNet3PowLimit,
	PowLimitBits:             0x1d00ffff,
	CoinbaseMaturity:         100,
	SubsidyReductionInterval: 210000,
	TargetTimespan:           time.Hour * 24 * 14, // 14 days
	TargetTimePerBlock:       time.Minute * 10,    // 10 minutes
	RetargetAdjustmentFactor: 4,                   // 25% less, 400% more
	ReduceMinDifficulty:      true,
	MinDiffReductionTime:     time.Minute * 20, // TargetTimePerBlock * 2
	GenerateSupported:        false,

	// Checkpoints ordered from oldest to newest.
	Checkpoints: []chaincfg.Checkpoint{
		//{500, newHashFromStr("00000000c674047be3a7b25fefe0b6416f6f4e88ff9b01ddc05471b8e2ea603a")},
		//{1000, newHashFromStr("00000000b747d47c3b38161693ad05e26924b3775a8be669751f969da836311e")},
		//{10000, newHashFromStr("000000000037079ff4c37eed57d00eb9ddfde8737b559ffa4101b11e76c97466")},
		//{25000, newHashFromStr("00000000000000c207c423ebb2d935e7b867b51710aaf72967666e83696f01e2")},
		//{35000, newHashFromStr("0000000047f9360bd7e79d3959bd32366e24b4182caf138a8b10d42add3b7fd7")},
		//{45000, newHashFromStr("0000000019ae521883b2597ed74cd21e2efa43fbf487815300cad96206d76f0e")},
	},

	// Consensus rule change deployments.
	//
	// The miner confirmation window is defined as:
	//   target proof of work timespan / target proof of work spacing
	RuleChangeActivationThreshold: 1512, // 75% of MinerConfirmationWindow
	MinerConfirmationWindow:       2016,
	Deployments: [chaincfg.DefinedDeployments]chaincfg.ConsensusDeployment{
		chaincfg.DeploymentTestDummy: {
			BitNumber: 28,
			DeploymentStarter: chaincfg.NewMedianTimeDeploymentStarter(
				time.Time{}, // Always available for vote
			),
			DeploymentEnder: chaincfg.NewMedianTimeDeploymentEnder(
				time.Time{}, // Never expires
			),
		},
		chaincfg.DeploymentTestDummyMinActivation: {
			BitNumber:                 22,
			CustomActivationThreshold: 1815,    // Only needs 90% hash rate.
			MinActivationHeight:       10_0000, // Can only activate after height 10k.
			DeploymentStarter: chaincfg.NewMedianTimeDeploymentStarter(
				time.Time{}, // Always available for vote
			),
			DeploymentEnder: chaincfg.NewMedianTimeDeploymentEnder(
				time.Time{}, // Never expires
			),
		},
		chaincfg.DeploymentCSV: {
			BitNumber: 0,
			DeploymentStarter: chaincfg.NewMedianTimeDeploymentStarter(
				time.Time{}, // Always available for vote
			),
			DeploymentEnder: chaincfg.NewMedianTimeDeploymentEnder(
				time.Time{}, // Never expires
			),
		},
		chaincfg.DeploymentSegwit: {
			BitNumber: 1,
			DeploymentStarter: chaincfg.NewMedianTimeDeploymentStarter(
				time.Time{}, // Always available for vote
			),
			DeploymentEnder: chaincfg.NewMedianTimeDeploymentEnder(
				time.Time{}, // Never expires
			),
		},
		chaincfg.DeploymentTaproot: {
			BitNumber: 2,
			DeploymentStarter: chaincfg.NewMedianTimeDeploymentStarter(
				time.Time{}, // Always available for vote
			),
			DeploymentEnder: chaincfg.NewMedianTimeDeploymentEnder(
				time.Time{}, // Never expires
			),
			CustomActivationThreshold: 1512, // 75%
		},
	},

	// Mempool parameters
	RelayNonStdTxs: true,

	// Human-readable part for Bech32 encoded segwit addresses, as defined in
	// BIP 173.
	Bech32HRPSegwit: "tb", // always tb for test net

	// Address encoding magics
	PubKeyHashAddrID:        0x6f, // starts with m or n
	ScriptHashAddrID:        0xc4, // starts with 2
	WitnessPubKeyHashAddrID: 0x03, // starts with QW
	WitnessScriptHashAddrID: 0x28, // starts with T7n
	PrivateKeyID:            0xef, // starts with 9 (uncompressed) or c (compressed)

	// BIP32 hierarchical deterministic extended key magics
	HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
	HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

	// BIP44 coin type used in the hierarchical deterministic path for
	// address generation.
	HDCoinType: 1,
}
