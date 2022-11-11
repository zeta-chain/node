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

var MagicHash = "0xcccb58610b0b65d5b1d8e5f16435254a787e324209d9b3877a8ece68859a0f55"

type ChainETHish struct {
	// TODO: Could these 3 be refactored out?
	tss     zetaclient.TSSSigner
	mpiAbi  abi.ABI
	context context.Context
	chainID *big.Int

	MpiContract                string
	DefaultDestinationContract string
	client                     *ethclient.Client
	name                       common.Chain
	id                         *big.Int
	topics                     [][]ethcommon.Hash
	channel                    chan types.Log
	subscription               ethereum.Subscription
}

func (cl *ChainETHish) Init() {
	cl.tss = GetZetaTestSignature()

	_abi, err := abi.JSON(strings.NewReader(AbiMpi))
	if err != nil {
		log.Err(err).Msg("abi.JSON")
		os.Exit(1)
	}
	cl.mpiAbi = _abi

	cl.context = context.TODO()

	chain, err := zetaclient.NewEVMChainClient(cl.name, nil, cl.tss, "", nil)
	if err != nil {
		log.Error().Err(err).Msg("NewEVMChainClient")
	}
	cl.client = chain.EvmClient

	_id, _ := cl.client.ChainID(cl.context)
	log.Debug().Msg(fmt.Sprintf("%s chain id %d", cl.name, _id))
	cl.id, err = cl.client.ChainID(context.TODO())
	if err != nil {
		fmt.Printf("Chain.id error %s\n", err)
		os.Exit(1)
	}

	cl.topics = make([][]ethcommon.Hash, 1)
	cl.topics[0] = []ethcommon.Hash{ethcommon.HexToHash(MagicHash)}
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(cl.MpiContract)},
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

func FindChainByID(id *big.Int) (*ChainETHish, error) {
	for _, chain := range AllChains {
		if chain.chainID.Cmp(id) == 0 {
			return chain, nil
		}
	}
	return nil, fmt.Errorf("not listening for chain with ID: %d", id)
}

func FindChainByName(name string) (*ChainETHish, error) {
	for _, chain := range AllChains {
		if chain.name.String() == name {
			return chain, nil
		}
	}
	return nil, fmt.Errorf("couldn't find chain: %s", name)
}
