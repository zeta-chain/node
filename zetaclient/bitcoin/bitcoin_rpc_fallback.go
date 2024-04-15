package bitcoin

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
)

var _ interfaces.BTCRPCClient = &RPCClientFallback{}

// RPCClientFallback - a decorator type for adding fallback capability to bitcoin rpc client
type RPCClientFallback struct {
	btcConfig  config.BTCConfig
	rpcClients *common.ClientQueue
	logger     zerolog.Logger
}

// NewRPCClientFallback - Constructor, reads config and connects to each client in the endpoint list
func NewRPCClientFallback(cfg config.BTCConfig, logger zerolog.Logger) (*RPCClientFallback, error) {
	if len(cfg.Endpoints) == 0 {
		return nil, errors.New("invalid endpoints")
	}
	rpcClientFallback := RPCClientFallback{
		btcConfig:  cfg,
		rpcClients: common.NewClientQueue(),
		logger:     logger,
	}
	for _, client := range cfg.Endpoints {
		logger.Info().Msgf("endpoint %s", client.RPCHost)
		connCfg := &rpcclient.ConnConfig{
			Host:         client.RPCHost,
			User:         client.RPCUsername,
			Pass:         client.RPCPassword,
			HTTPPostMode: true,
			DisableTLS:   true,
			Params:       client.RPCParams,
		}

		rpcClient, err := rpcclient.New(connCfg, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating rpc client: %s", err)
		}
		err = rpcClient.Ping()
		if err != nil {
			return nil, fmt.Errorf("error ping the bitcoin server: %s", err)
		}
		rpcClientFallback.rpcClients.Append(rpcClient)
	}
	return &rpcClientFallback, nil
}

// Below is an implementation of the BTCRPCClient interface. The logic is similar for all functions, the first client
// in the queue is used to attempt the rpc call. If this fails then it will attempt to call the next client in the list
// until it has tried them all.

func (R *RPCClientFallback) GetNetworkInfo() (*btcjson.GetNetworkInfoResult, error) {
	var res *btcjson.GetNetworkInfoResult
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).GetNetworkInfo()
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) CreateWallet(name string, opts ...rpcclient.CreateWalletOpt) (*btcjson.CreateWalletResult, error) {
	var res *btcjson.CreateWalletResult
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).CreateWallet(name, opts...)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) GetNewAddress(account string) (btcutil.Address, error) {
	var res btcutil.Address
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).GetNewAddress(account)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) GenerateToAddress(numBlocks int64, address btcutil.Address, maxTries *int64) ([]*chainhash.Hash, error) {
	var res []*chainhash.Hash
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).GenerateToAddress(numBlocks, address, maxTries)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) GetBalance(account string) (btcutil.Amount, error) {
	var res btcutil.Amount
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).GetBalance(account)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error) {
	var res *chainhash.Hash
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).SendRawTransaction(tx, allowHighFees)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) ListUnspent() ([]btcjson.ListUnspentResult, error) {
	var res []btcjson.ListUnspentResult
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).ListUnspent()
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) ListUnspentMinMaxAddresses(minConf int, maxConf int, addrs []btcutil.Address) ([]btcjson.ListUnspentResult, error) {
	var res []btcjson.ListUnspentResult
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).ListUnspentMinMaxAddresses(minConf, maxConf, addrs)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) EstimateSmartFee(confTarget int64, mode *btcjson.EstimateSmartFeeMode) (*btcjson.EstimateSmartFeeResult, error) {
	var res *btcjson.EstimateSmartFeeResult
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).EstimateSmartFee(confTarget, mode)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) GetTransaction(txHash *chainhash.Hash) (*btcjson.GetTransactionResult, error) {
	var res *btcjson.GetTransactionResult
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).GetTransaction(txHash)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) GetRawTransaction(txHash *chainhash.Hash) (*btcutil.Tx, error) {
	var res *btcutil.Tx
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).GetRawTransaction(txHash)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) GetRawTransactionVerbose(txHash *chainhash.Hash) (*btcjson.TxRawResult, error) {
	var res *btcjson.TxRawResult
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).GetRawTransactionVerbose(txHash)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) GetBlockCount() (int64, error) {
	var res int64
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).GetBlockCount()
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) GetBlockHash(blockHeight int64) (*chainhash.Hash, error) {
	var res *chainhash.Hash
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).GetBlockHash(blockHeight)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) GetBlockVerbose(blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	var res *btcjson.GetBlockVerboseResult
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).GetBlockVerbose(blockHash)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) GetBlockVerboseTx(blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error) {
	var res *btcjson.GetBlockVerboseTxResult
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).GetBlockVerboseTx(blockHash)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}

func (R *RPCClientFallback) GetBlockHeader(blockHash *chainhash.Hash) (*wire.BlockHeader, error) {
	var res *wire.BlockHeader
	var err error

	for i := 0; i < R.rpcClients.Length(); i++ {
		if client := R.rpcClients.First(); client != nil {
			res, err = client.(interfaces.BTCRPCClient).GetBlockHeader(blockHash)
		}
		if err != nil {
			R.logger.Debug().Err(err).Msg("client endpoint failed attempting fallback client")
			R.rpcClients.Next()
			continue
		}
		break
	}
	return res, err
}
