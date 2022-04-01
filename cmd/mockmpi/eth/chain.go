package eth

import (
	"context"
	"fmt"
	common2 "github.com/zeta-chain/zetacore/cmd/mockmpi/common"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient"
)

// What is this?
const MAGIC_HASH = "0xcccb58610b0b65d5b1d8e5f16435254a787e324209d9b3877a8ece68859a0f55"

type ChainETHish struct {
	// TODO: Could these 3 be refactored out?
	tss      zetaclient.TSSSigner
	mpi_abi  abi.ABI
	context  context.Context
	Chain_id uint16

	MPI_CONTRACT                 string
	DEFAULT_DESTINATION_CONTRACT string
	client                       *ethclient.Client
	name                         string
	id                           *big.Int
	topics                       [][]ethcommon.Hash
	channel                      chan types.Log
	subscription                 ethereum.Subscription
}

func RegisterChains() {
	common2.ALL_CHAINS = append(common2.ALL_CHAINS, &ChainETHish{
		name:                         "ETH",
		MPI_CONTRACT:                 "0x132b042bD5198a48E4D273f46b979E5f13Bd9239",
		DEFAULT_DESTINATION_CONTRACT: "0xFf6B270ac3790589A1Fe90d0303e9D4d9A54FD1A",
		Chain_id:                     5,
	})
	common2.ALL_CHAINS = append(common2.ALL_CHAINS, &ChainETHish{
		name:                         "BSC",
		MPI_CONTRACT:                 "0x96cE47e42A73649CFe33d93D93ACFbEc6FD5ee14",
		DEFAULT_DESTINATION_CONTRACT: "0xF47bd84B86d1667e7621c38c72C6905Ca8710b0d",
		Chain_id:                     97,
	})
	//{
	//	name:                         common.Chain("POLYGON"),
	//	MPI_CONTRACT:                 "0x692E8A48634B530b4BFF1e621FC18C82F471892c",
	//	DEFAULT_DESTINATION_CONTRACT: "0x22696Bef41E49FEf5beac1D4765a5b7B1E0Dcb01",
	//},
}

func (cl *ChainETHish) ID() uint16 {
	return cl.Chain_id
}

func (cl *ChainETHish) Name() string {
	return cl.name
}

func (cl *ChainETHish) Init() {
	cl.tss = GetZetaTestSignature()

	_abi, err := abi.JSON(strings.NewReader(ABI_MPI))
	if err != nil {
		log.Err(err).Msg("abi.JSON")
		os.Exit(1)
	}
	cl.mpi_abi = _abi

	cl.context = context.TODO()

	chain, err := zetaclient.NewChainObserver(common.Chain(cl.name), nil, cl.tss, "")
	if err != nil {
		log.Error().Err(err).Msg("NewChainObserver")
	}
	cl.client = chain.Client

	_id, _ := cl.client.ChainID(cl.context)
	log.Debug().Msg(fmt.Sprintf("%s chain id %d", cl.name, _id))
	cl.id, err = cl.client.ChainID(context.TODO())
	if err != nil {
		fmt.Printf("Chain.id error %s\n", err)
		os.Exit(1)
	}

	cl.topics = make([][]ethcommon.Hash, 1)
	cl.topics[0] = []ethcommon.Hash{ethcommon.HexToHash(MAGIC_HASH)}
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(cl.MPI_CONTRACT)},
		Topics:    cl.topics,
	}

	cl.channel = make(chan types.Log)

	_subscription, err := cl.client.SubscribeFilterLogs(cl.context, query, cl.channel)
	if err != nil {
		log.Fatal().Err(err).Msg("SubscribeFilterLogs")
		os.Exit(1)
	}
	cl.subscription = _subscription
}

func (cl *ChainETHish) Start() {
	cl.Init()
	cl.Listen()
}
