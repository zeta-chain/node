package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	etherminttypes "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// ZetaTxServer is a ZetaChain tx server for smoke test
type ZetaTxServer struct {
	clientCtx client.Context
	txFactory tx.Factory
}

// NewZetaTxServer returns a new TxServer with provided account
func NewZetaTxServer(rpcAddr string, names []string, mnemonics []string) (ZetaTxServer, error) {
	ctx := context.Background()

	if len(names) != len(mnemonics) {
		return ZetaTxServer{}, errors.New("invalid names and mnemonics")
	}

	// initialize rpc and check status
	rpc, err := rpchttp.New(rpcAddr, "/websocket")
	if err != nil {
		return ZetaTxServer{}, fmt.Errorf("failed to initialize rpc: %s", err.Error())
	}
	if _, err = rpc.Status(ctx); err != nil {
		return ZetaTxServer{}, fmt.Errorf("failed to query rpc: %s", err.Error())
	}

	// initialize codec
	cdc, reg := newCodec()

	// initialize keyring
	kr := keyring.NewInMemory(cdc)

	// create accounts
	for i := range names {
		r, err := kr.NewAccount(names[i], mnemonics[i], "", sdktypes.FullFundraiserPath, hd.Secp256k1)
		if err != nil {
			return ZetaTxServer{}, fmt.Errorf("failed to create account: %s", err.Error())
		}
		_, err = r.GetAddress()
		if err != nil {
			return ZetaTxServer{}, fmt.Errorf("failed to get account address: %s", err.Error())
		}
	}

	clientCtx := newContext(rpc, cdc, reg, kr)
	txf := newFactory(clientCtx)

	return ZetaTxServer{
		clientCtx: clientCtx,
		txFactory: txf,
	}, nil
}

// BroadcastTx broadcasts a tx to ZetaChain with the provided msg from the account
func (zts ZetaTxServer) BroadcastTx(account string, msg sdktypes.Msg) (*sdktypes.TxResponse, error) {
	// Find number and sequence and set it
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return nil, err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return nil, err
	}
	accountNumber, accountSeq, err := zts.clientCtx.AccountRetriever.GetAccountNumberSequence(zts.clientCtx, addr)
	if err != nil {
		return nil, err
	}
	zts.txFactory = zts.txFactory.WithAccountNumber(accountNumber).WithSequence(accountSeq)

	txBuilder, err := zts.txFactory.BuildUnsignedTx(msg)
	if err != nil {
		return nil, err
	}

	// Sign tx
	err = tx.Sign(zts.txFactory, account, txBuilder, true)
	if err != nil {
		return nil, err
	}
	txBytes, err := zts.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	// Broadcast tx
	return zts.clientCtx.BroadcastTx(txBytes)
}

// newCodec returns the codec for msg server
func newCodec() (*codec.ProtoCodec, codectypes.InterfaceRegistry) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	sdktypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	authz.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	stakingtypes.RegisterInterfaces(interfaceRegistry)
	slashingtypes.RegisterInterfaces(interfaceRegistry)
	upgradetypes.RegisterInterfaces(interfaceRegistry)
	distrtypes.RegisterInterfaces(interfaceRegistry)
	evidencetypes.RegisterInterfaces(interfaceRegistry)
	crisistypes.RegisterInterfaces(interfaceRegistry)
	evmtypes.RegisterInterfaces(interfaceRegistry)
	etherminttypes.RegisterInterfaces(interfaceRegistry)
	crosschaintypes.RegisterInterfaces(interfaceRegistry)
	emissionstypes.RegisterInterfaces(interfaceRegistry)
	fungibletypes.RegisterInterfaces(interfaceRegistry)
	observertypes.RegisterInterfaces(interfaceRegistry)

	return cdc, interfaceRegistry
}

// newContext returns the client context for msg server
func newContext(rpc *rpchttp.HTTP, cdc *codec.ProtoCodec, reg codectypes.InterfaceRegistry, kr keyring.Keyring) client.Context {
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	return client.Context{}.
		WithChainID(ZetaChainID).
		WithInterfaceRegistry(reg).
		WithCodec(cdc).
		WithTxConfig(txConfig).
		WithLegacyAmino(codec.NewLegacyAmino()).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithBroadcastMode(flags.BroadcastBlock).
		WithClient(rpc).
		WithSkipConfirmation(true).
		WithFromName("creator").
		WithFromAddress(sdktypes.AccAddress{}).
		WithKeyring(kr).
		WithAccountRetriever(authtypes.AccountRetriever{})
}

// newFactory returns the tx factory for msg server
func newFactory(clientCtx client.Context) tx.Factory {
	return tx.Factory{}.
		WithChainID(clientCtx.ChainID).
		WithKeybase(clientCtx.Keyring).
		WithGas(10000000).
		WithGasAdjustment(1).
		WithSignMode(signing.SignMode_SIGN_MODE_UNSPECIFIED).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithTxConfig(clientCtx.TxConfig).
		WithFees("1000azeta")
}

type messageLog struct {
	Events []event `json:"events"`
}

type event struct {
	Type       string      `json:"type"`
	Attributes []attribute `json:"attributes"`
}

type attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// fetchAttribute fetches the attribute from the tx response
func fetchAttribute(res *sdktypes.TxResponse, key string) (string, error) {
	var logs []messageLog
	err := json.Unmarshal([]byte(res.RawLog), &logs)
	if err != nil {
		return "", err
	}

	for _, log := range logs {
		for _, event := range log.Events {
			for _, attr := range event.Attributes {
				if attr.Key == key {
					return attr.Value, nil
				}
			}
		}
	}

	return "", errors.New("attribute not found")
}
