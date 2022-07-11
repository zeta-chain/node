package main

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
	"math/big"
	"strings"
)

const (
	ERC20_ABI_STRING = `[{"inputs":[{"internalType":"uint256","name":"initialSupply","type":"uint256"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]`
)

type Probe struct {
	Client           *ethclient.Client
	ConnectorABI     *abi.ABI
	ConnectorAddress ethcommon.Address
	ERC20ABI         *abi.ABI
	TokenAddress     ethcommon.Address
	Address          ethcommon.Address
	ChainID          *big.Int
	ChainName        string
}

func NewProbe(client *ethclient.Client, connectorABI *abi.ABI, address string, chainID *big.Int, connectorAddress string, tokenAddress string) *Probe {
	ERC20ABI, err := abi.JSON(strings.NewReader(ERC20_ABI_STRING))
	if err != nil {
		log.Fatal().Err(err).Msg("parse erc20 abi error")
		return nil
	}

	return &Probe{
		Client:           client,
		ConnectorABI:     connectorABI,
		ERC20ABI:         &ERC20ABI,
		Address:          ethcommon.HexToAddress(address),
		ChainID:          chainID,
		ConnectorAddress: ethcommon.HexToAddress(connectorAddress),
		TokenAddress:     ethcommon.HexToAddress(tokenAddress),
	}
}

func (probe *Probe) SendTransaction() error {

	return nil
}

// no decimals
func (probe *Probe) GetBalance() (*big.Int, error) {
	bal, err := probe.Client.BalanceAt(context.Background(), probe.Address, nil)
	if err != nil {
		return nil, err
	} else {
		return bal.Div(bal, big.NewInt(1e18)), nil
	}
}

func (probe *Probe) GetZetaBalance() (*big.Int, error) {
	input, err := probe.ERC20ABI.Pack("balanceOf", probe.Address)
	if err != nil {
		return nil, err
	}
	res, err := probe.Client.CallContract(context.Background(), ethereum.CallMsg{
		From: probe.Address,
		To:   &probe.TokenAddress,
		Data: input,
	}, nil)
	if err != nil {
		return nil, err
	}
	output, err := probe.ERC20ABI.Unpack("balanceOf", res)

	bal := *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	return bal.Div(bal, big.NewInt(1e18)), err
}
