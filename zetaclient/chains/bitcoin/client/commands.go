package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

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

func (c *Client) CreateRawTransaction(
	ctx context.Context,
	inputs []types.TransactionInput,
	amounts map[btcutil.Address]btcutil.Amount,
	lockTime *int64,
) (*wire.MsgTx, error) {
	convertedAmounts := make(map[string]float64, len(amounts))
	for addr, amount := range amounts {
		convertedAmounts[addr.String()] = amount.ToBTC()
	}

	cmd := types.NewCreateRawTransactionCmd(inputs, convertedAmounts, lockTime)

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create raw tx")
	}

	// Decode the serialized transaction hex to raw bytes.
	serializedTx, err := unmarshalHex(out)
	if err != nil {
		return nil, err
	}

	// Deserialize the transaction and return it.
	var msgTx wire.MsgTx

	if err = msgTx.Deserialize(bytes.NewReader(serializedTx)); err != nil {
		return nil, err
	}

	return &msgTx, nil
}

func (c *Client) SignRawTransactionWithWallet2(
	ctx context.Context,
	tx *wire.MsgTx,
	inputs []types.RawTxWitnessInput,
) (*wire.MsgTx, bool, error) {
	if tx == nil {
		return nil, false, errors.New("tx is nil")
	}

	// Serialize the transaction and convert to hex string.
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return nil, false, errors.Wrap(err, "unable to serialize tx")
	}

	txHex := hex.EncodeToString(buf.Bytes())

	cmd := types.NewSignRawTransactionWithWalletCmd(txHex, &inputs, nil)

	out, err := c.sendCommand(ctx, cmd)
	if err != nil {
		return nil, false, errors.Wrap(err, "unable to sign raw tx")
	}

	result, err := unmarshalPtr[types.SignRawTransactionWithWalletResult](out)
	if err != nil {
		return nil, false, errors.Wrap(err, "unable to unmarshal sign raw tx result")
	}

	// Decode the serialized transaction hex to raw bytes.
	serializedTx, err := hex.DecodeString(result.Hex)
	if err != nil {
		return nil, false, err
	}

	// Deserialize the transaction and return it.
	var msgTx wire.MsgTx
	if err = msgTx.Deserialize(bytes.NewReader(serializedTx)); err != nil {
		return nil, false, err
	}

	return &msgTx, result.Complete, nil
}

func (c *Client) ImportAddress(ctx context.Context, address string) error {
	cmd := types.NewImportAddressCmd(address, "", nil)

	_, err := c.sendCommand(ctx, cmd)
	return err
}

func (c *Client) ImportPrivKeyRescan(ctx context.Context, privKeyWIF *btcutil.WIF, label string, rescan bool) error {
	wif := ""
	if privKeyWIF != nil {
		wif = privKeyWIF.String()
	}

	cmd := types.NewImportPrivKeyCmd(wif, &label, &rescan)

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

	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal params")
	}

	body := fmt.Sprintf(
		`{"jsonrpc":"%s","id":%d,"method":"%s","params":%s}`,
		rpcVersion,
		commandID,
		method,
		paramsBytes,
	)

	req, err := c.newRequest(ctx, []byte(body))
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
