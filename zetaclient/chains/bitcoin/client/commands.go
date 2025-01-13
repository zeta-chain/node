package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"

	types "github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"github.com/tendermint/btcd/btcjson"
)

const (
	// github.com/btcsuite/btcd@v0.24.2/rpcclient/rawtransactions.go:22
	defaultMaxFeeRate types.BTCPerkvB = 0.1
)

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.GetBlockCount(ctx)
	return errors.Wrap(err, "ping failed")
}

func (c *Client) GetNetworkInfo(ctx context.Context) (*types.GetNetworkInfoResult, error) {
	out, err := c.sendCommand(ctx, types.NewGetNetworkInfoCmd())
	if err != nil {
		return nil, errors.Wrap(err, "unable to get network info")
	}

	return unmarshalPtr[types.GetNetworkInfoResult](out)
}

func (c *Client) GetBlockCount(ctx context.Context) (int64, error) {
	out, err := c.sendCommand(ctx, types.NewGetBlockCountCmd())
	if err != nil {
		return 0, errors.Wrap(err, "unable to get block count")
	}

	return unmarshal[int64](out)
}

func (c *Client) GetBlockHash(ctx context.Context, blockHeight int64) (*chainhash.Hash, error) {
	out, err := c.sendCommand(ctx, types.NewGetBlockHashCmd(blockHeight))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get block hash for %d", blockHeight)
	}

	return unmarshalPtr[chainhash.Hash](out)
}

func (c *Client) GetBlockHeader(ctx context.Context, hash *chainhash.Hash) (*wire.BlockHeader, error) {
	cmd := types.NewGetBlockHeaderCmd(hash.String(), btcjson.Bool(false))

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get block header for %s", hash.String())
	}

	var bhHex string
	err = json.Unmarshal(out, &bhHex)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal hex")
	}

	serializedBH, err := hex.DecodeString(bhHex)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode hex")
	}

	var bh wire.BlockHeader
	if err = bh.Deserialize(bytes.NewReader(serializedBH)); err != nil {
		return nil, errors.Wrap(err, "unable to deserialize block header")
	}

	return &bh, nil
}

func (c *Client) GetBlockVerbose(ctx context.Context, hash *chainhash.Hash) (*types.GetBlockVerboseTxResult, error) {
	cmd := types.NewGetBlockCmd(hash.String(), btcjson.Int(2))

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get block hash verbose")
	}

	return unmarshalPtr[types.GetBlockVerboseTxResult](out)
}

func (c *Client) GetTransaction(ctx context.Context, hash *chainhash.Hash) (*types.GetTransactionResult, error) {
	out, err := c.sendCommand(ctx, types.NewGetTransactionCmd(hash.String(), nil))
	if err != nil {
		return nil, errors.Wrap(err, "unable to get transaction")
	}

	return unmarshalPtr[types.GetTransactionResult](out)
}

func (c *Client) GetRawTransaction(ctx context.Context, hash *chainhash.Hash) (*btcutil.Tx, error) {
	cmd := types.NewGetRawTransactionCmd(hash.String(), btcjson.Int(0))

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get raw tx")
	}

	var txHex string
	err = json.Unmarshal(out, &txHex)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal raw tx")
	}

	// Decode the serialized transaction hex to raw bytes.
	serializedTx, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode raw tx")
	}

	// Deserialize the transaction and return it.
	var msgTx wire.MsgTx
	if err = msgTx.Deserialize(bytes.NewReader(serializedTx)); err != nil {
		return nil, errors.Wrap(err, "unable to deserialize raw tx")
	}

	return btcutil.NewTx(&msgTx), nil
}

func (c *Client) GetRawTransactionVerbose(ctx context.Context, hash *chainhash.Hash) (*types.TxRawResult, error) {
	cmd := types.NewGetRawTransactionCmd(hash.String(), btcjson.Int(1))

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get raw tx verbose")
	}

	return unmarshalPtr[types.TxRawResult](out)
}

// SendRawTransaction github.com/btcsuite/btcd@v0.24.2/rpcclient/rawtransactions.go
func (c *Client) SendRawTransaction(ctx context.Context, tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error) {
	if tx == nil {
		return nil, errors.New("tx is nil")
	}

	// Serialize the transaction and convert to hex string.
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return nil, errors.Wrap(err, "unable to serialize tx")
	}

	txHex := hex.EncodeToString(buf.Bytes())

	// Using a 0 MaxFeeRate is interpreted as a maximum fee rate not
	// being enforced by bitcoind.
	var maxFeeRate types.BTCPerkvB
	if !allowHighFees {
		maxFeeRate = defaultMaxFeeRate
	}

	cmd := types.NewBitcoindSendRawTransactionCmd(txHex, maxFeeRate)

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to send raw tx")
	}

	var txHashStr string
	err = json.Unmarshal(out, &txHashStr)
	if err != nil {
		return nil, err
	}

	return chainhash.NewHashFromStr(txHashStr)
}

func (c *Client) ListUnspent(ctx context.Context) ([]btcjson.ListUnspentResult, error) {
	cmd := btcjson.NewListUnspentCmd(nil, nil, nil)

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list unspent")
	}

	return unmarshal[[]btcjson.ListUnspentResult](out)
}

func (c *Client) ListUnspentMinMaxAddresses(
	ctx context.Context,
	minConf, maxConf int,
	addresses []btcutil.Address,
) ([]btcjson.ListUnspentResult, error) {
	stringAddresses := make([]string, 0, len(addresses))
	for _, a := range addresses {
		stringAddresses = append(stringAddresses, a.EncodeAddress())
	}

	cmd := btcjson.NewListUnspentCmd(&minConf, &maxConf, &stringAddresses)

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list unspent")
	}

	return unmarshal[[]btcjson.ListUnspentResult](out)
}

func (c *Client) EstimateSmartFee(
	ctx context.Context,
	confTarget int64,
	mode *types.EstimateSmartFeeMode,
) (*types.EstimateSmartFeeResult, error) {
	cmd := types.NewEstimateSmartFeeCmd(confTarget, mode)

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to estimate smart fee")
	}

	return unmarshalPtr[types.EstimateSmartFeeResult](out)
}

func (c *Client) CreateWallet(
	ctx context.Context,
	name string,
	opts ...rpcclient.CreateWalletOpt,
) (*types.CreateWalletResult, error) {
	cmd := types.NewCreateWalletCmd(name, nil, nil, nil, nil)
	for _, opt := range opts {
		opt(cmd)
	}

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create wallet")
	}

	return unmarshalPtr[types.CreateWalletResult](out)
}

func (c *Client) GetBalance(ctx context.Context, account string) (btcutil.Amount, error) {
	cmd := types.NewGetBalanceCmd(&account, nil)

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return 0, errors.Wrap(err, "unable to get balance")
	}

	balanceRaw, err := unmarshal[float64](out)
	if err != nil {
		return 0, errors.Wrap(err, "unable to unmarshal balance")
	}

	return btcutil.NewAmount(balanceRaw)
}

func (c *Client) GetNewAddress(ctx context.Context, account string) (btcutil.Address, error) {
	cmd := types.NewGetNewAddressCmd(&account, nil)

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get new address")
	}

	var addr string
	err = json.Unmarshal(out, &addr)
	if err != nil {
		return nil, err
	}

	return btcutil.DecodeAddress(addr, &c.params)
}

func (c *Client) GenerateToAddress(
	ctx context.Context,
	numBlocks int64,
	address btcutil.Address,
	maxTries *int64,
) ([]*chainhash.Hash, error) {
	cmd := types.NewGenerateToAddressCmd(numBlocks, address.EncodeAddress(), maxTries)

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate to address")
	}

	result, err := unmarshal[[]string](out)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal to strings")
	}

	convertedResult := make([]*chainhash.Hash, len(result))
	for i, hashString := range result {
		convertedResult[i], err = chainhash.NewHashFromStr(hashString)
		if err != nil {
			return nil, err
		}
	}

	return convertedResult, nil
}
