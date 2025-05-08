package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"time"

	types "github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
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

	str, err := unmarshal[string](out)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal block hash")
	}

	return chainhash.NewHashFromStr(str)
}

func (c *Client) GetBlockHeader(ctx context.Context, hash *chainhash.Hash) (*wire.BlockHeader, error) {
	cmd := types.NewGetBlockHeaderCmd(hash.String(), types.Bool(false))

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get block header for %s", hash.String())
	}

	serializedBH, err := unmarshalHex(out)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode hex")
	}

	var bh wire.BlockHeader
	if err = bh.Deserialize(bytes.NewReader(serializedBH)); err != nil {
		return nil, errors.Wrap(err, "unable to deserialize block header")
	}

	return &bh, nil
}

// GetRawMempool fetches all mempool transaction hashes.
func (c *Client) GetRawMempool(ctx context.Context) ([]*chainhash.Hash, error) {
	cmd := types.NewGetRawMempoolCmd(types.Bool(false))

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get raw mempool")
	}

	txHashStrs, err := unmarshal[[]string](out)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal to strings")
	}

	txHashes := make([]*chainhash.Hash, len(txHashStrs))
	for i, hashString := range txHashStrs {
		txHashes[i], err = chainhash.NewHashFromStr(hashString)
		if err != nil {
			return nil, err
		}
	}

	return txHashes, nil
}

// GetMempoolEntry fetches the mempool entry for the given transaction hash.
func (c *Client) GetMempoolEntry(ctx context.Context, txHash string) (*types.GetMempoolEntryResult, error) {
	cmd := types.NewGetMempoolEntryCmd(txHash)

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get mempool entry for %s", txHash)
	}

	return unmarshalPtr[types.GetMempoolEntryResult](out)
}

func (c *Client) GetBlockVerbose(ctx context.Context, hash *chainhash.Hash) (*types.GetBlockVerboseTxResult, error) {
	cmd := types.NewGetBlockCmd(hash.String(), types.Int(2))

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
	cmd := types.NewGetRawTransactionCmd(hash.String(), types.Int(0))

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get raw tx")
	}

	// Decode the serialized transaction hex to raw bytes.
	serializedTx, err := unmarshalHex(out)
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
	cmd := types.NewGetRawTransactionCmd(hash.String(), types.Int(1))

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

	txHashStr, err := unmarshal[string](out)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal tx hash")
	}

	return chainhash.NewHashFromStr(txHashStr)
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

// IsTxStuckInMempool checks if the transaction is stuck in the mempool.
//
// A pending tx with 'confirmations == 0' will be considered stuck due to excessive pending time.
func (c *Client) IsTxStuckInMempool(
	ctx context.Context,
	txHash string,
	maxWaitBlocks int64,
) (stuck bool, pendingFor time.Duration, err error) {
	lastBlock, err := c.GetBlockCount(ctx)
	if err != nil {
		return false, 0, errors.Wrap(err, "GetBlockCount failed")
	}

	entry, err := c.GetMempoolEntry(ctx, txHash)
	if err != nil {
		// not a mempool tx, of course not stuck
		if isTxNotInMempoolError(err) {
			return false, 0, nil
		}

		return false, 0, errors.Wrap(err, "GetMempoolEntry failed")
	}

	const blockTimeBTC = 10 * time.Minute

	// is the tx pending for too long?
	pendingFor = time.Since(time.Unix(entry.Time, 0))
	maxPendingFor := blockTimeBTC * time.Duration(maxWaitBlocks)
	pendingDeadline := entry.Height + maxWaitBlocks

	// the block mining is frozen in Regnet for E2E test
	if c.isRegnet {
		maxPendingFor := time.Second * time.Duration(maxWaitBlocks)

		stuck = pendingFor > maxPendingFor && entry.Height == lastBlock

		return stuck, pendingFor, nil
	}

	stuck = pendingFor > maxPendingFor && lastBlock > pendingDeadline

	return stuck, pendingFor, nil
}

func (c *Client) ListUnspent(ctx context.Context) ([]types.ListUnspentResult, error) {
	cmd := types.NewListUnspentCmd(nil, nil, nil)

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list unspent")
	}

	return unmarshal[[]types.ListUnspentResult](out)
}

func (c *Client) ListUnspentMinMaxAddresses(
	ctx context.Context,
	minConf, maxConf int,
	addresses []btcutil.Address,
) ([]types.ListUnspentResult, error) {
	stringAddresses := make([]string, 0, len(addresses))
	for _, a := range addresses {
		stringAddresses = append(stringAddresses, a.EncodeAddress())
	}

	cmd := types.NewListUnspentCmd(&minConf, &maxConf, &stringAddresses)

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list unspent")
	}

	return unmarshal[[]types.ListUnspentResult](out)
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

	addr, err := unmarshal[string](out)
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

func (c *Client) ImportAddress(ctx context.Context, address string) error {
	cmd := types.NewImportAddressCmd(address, "", nil)

	_, err := c.sendCommand(ctx, cmd)
	return err
}

func (c *Client) RawRequest(ctx context.Context, method string, params []json.RawMessage) (json.RawMessage, error) {
	switch {
	case method == "":
		return nil, errors.New("no method")
	case params == nil:
		params = []json.RawMessage{}
	}

	payload := struct {
		Version string            `json:"jsonrpc"`
		ID      uint64            `json:"id"`
		Method  string            `json:"method"`
		Params  []json.RawMessage `json:"params"`
	}{
		Version: string(rpcVersion),
		ID:      commandID,
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal body")
	}

	req, err := c.newRequest(ctx, body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request")
	}

	res, err := c.sendRequest(req, method)
	switch {
	case err != nil:
		return nil, errors.Wrapf(err, "%q failed", method)
	case res.Error != nil:
		return nil, errors.Wrapf(res.Error, "got rpc error for %q", method)
	}

	return res.Result, nil
}
