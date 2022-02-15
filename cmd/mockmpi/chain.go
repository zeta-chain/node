package main

import (
	"context"
	"fmt"
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
const MAGIC_HASH = "0x38f8fa9ce079e7e087c700936fd84330f80123e22a6aea6e125b55e95dcd585a"

type ChainETHish struct {
	tss                          zetaclient.TSSSigner
	mpi_abi                      abi.ABI
	MPI_CONTRACT                 string
	DEFAULT_DESTINATION_CONTRACT string
	context                      context.Context
	client                       *ethclient.Client
	name                         common.Chain
	id                           *big.Int
	topics                       [][]ethcommon.Hash
	channel                      chan types.Log
	subscription                 ethereum.Subscription
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

	chain, err := zetaclient.NewChainObserver(cl.name, nil, cl.tss, "")
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

func FindChainByID(id string) (*ChainETHish, error) {
	for _, chain := range ALL_CHAINS {
		if chain.id.String() == id {
			return chain, nil
		}
	}
	return nil, fmt.Errorf("Not listening for chain with ID: %s", id)
}

func FindChainByName(name string) (*ChainETHish, error) {
	for _, chain := range ALL_CHAINS {
		if chain.name.String() == name {
			return chain, nil
		}
	}
	return nil, fmt.Errorf("Couldn't find chain: %s", name)
}
