// Package rpc implements a client for HTTP-RPC using toncenter API V2 spec
// See: https://toncenter.com/api/v2
// See: https://github.com/toncenter/ton-http-api
// See: https://docs.ton.org/v3/guidelines/dapps/apis-sdks/ton-http-apis
package rpc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

type Client struct {
	client   *http.Client
	endpoint string
}

const pageSize = 100

type Opt func(c *Client)

func WithHTTPClient(client *http.Client) Opt {
	return func(c *Client) { c.client = client }
}

var ErrNotFound = errors.New("not found")

// New Client constructor
func New(endpoint string, opts ...Opt) *Client {
	const defaultTimeout = 10 * time.Second

	// todo metrics

	// Most API providers expose a url with api in in the path
	// - https://ton-testnet.core.chainstack.com/$key/api/v2
	// - https://$node.ton-mainnet.quiknode.pro/$key/
	//
	// And we need to add /jsonRPC to the end of the url
	endpoint = strings.TrimRight(endpoint, "/") + "/jsonRPC"

	c := &Client{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) GetMasterchainInfo(ctx context.Context) (MasterchainInfo, error) {
	var info MasterchainInfo

	err := c.callAndUnmarshal(ctx, "getMasterchainInfo", nil, &info)

	return info, err
}

func (c *Client) GetBlockHeader(ctx context.Context, blockID BlockIDExt) (BlockHeader, error) {
	// todo should we have cache?

	params := map[string]any{
		"workchain": blockID.Workchain,
		"shard":     blockID.Shard,
		"seqno":     blockID.Seqno,
	}

	var header BlockHeader

	err := c.callAndUnmarshal(ctx, "getBlockHeader", params, &header)

	return header, err
}

func (c *Client) HealthCheck(ctx context.Context) (time.Time, error) {
	info, err := c.GetMasterchainInfo(ctx)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to get masterchain info")
	}

	blockHeader, err := c.GetBlockHeader(ctx, info.Last)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to get block header")
	}

	blockTime := time.Unix(int64(blockHeader.GenUtime), 0).UTC()

	return blockTime, nil
}

func (c *Client) GetAccountState(ctx context.Context, acc ton.AccountID) (Account, error) {
	params := map[string]any{
		"address": acc.ToRaw(),
	}

	var account Account

	err := c.callAndUnmarshal(ctx, "getExtendedAddressInformation", params, &account)

	return account, err
}

// getLastTransactionHash returns logical time and hash of the last transaction
func (c *Client) getLastTransactionHash(ctx context.Context, acc ton.AccountID) (uint64, tlb.Bits256, error) {
	state, err := c.GetAccountState(ctx, acc)
	if err != nil {
		return 0, tlb.Bits256{}, errors.Wrap(err, "unable to get account state")
	}

	if state.Status != tlb.AccountActive {
		return 0, tlb.Bits256{}, errors.New("account is not active")
	}

	return state.LastTxLT, state.LastTxHash, nil
}

func (c *Client) GetConfigParam(ctx context.Context, index uint32) (*boc.Cell, error) {
	params := map[string]any{
		"config_id": index,
	}

	response, err := c.call(ctx, "getConfigParam", params)
	if err != nil {
		return nil, err
	}

	rawBase64 := gjson.GetBytes(response, "config.bytes").String()
	if rawBase64 == "" {
		return nil, errors.Errorf("config.bytes is empty (%s)", response)
	}

	cells, err := boc.DeserializeBocBase64(rawBase64)

	switch {
	case err != nil:
		return nil, errors.Wrapf(err, "unable to deserialize boc from %q", rawBase64)
	case len(cells) == 0:
		return nil, errors.Errorf("expected at least one cell, got 0")
	default:
		return cells[0], nil
	}
}

func (c *Client) GetTransactions(
	ctx context.Context,
	count uint32,
	accountID ton.AccountID,
	lt uint64,
	hash ton.Bits256,
) ([]ton.Transaction, error) {
	params := map[string]any{
		"address": accountID.ToRaw(),
		"limit":   count,
	}

	if lt > 0 {
		params["lt"] = lt
	}

	if hash != (ton.Bits256{}) {
		params["hash"] = hash.Base64()
	}

	// todo should we support ARCHIVAL nodes?
	// «By default getTransaction request is processed by any available liteserver.
	// If archival=true ONLY lite-servers with full history are used»

	response, err := c.call(ctx, "getTransactions", params)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get transactions")
	}

	// https://github.com/tidwall/gjson?tab=readme-ov-file#path-syntax
	txsRaw := gjson.GetBytes(response, "#.data").Array()
	if len(txsRaw) == 0 {
		return nil, nil
	}

	txs := make([]ton.Transaction, 0, len(txsRaw))
	for _, txRaw := range txsRaw {
		var tx ton.Transaction

		if err := unmarshalFromBase64(txRaw.String(), &tx); err != nil {
			return nil, errors.Wrapf(err, "unable to unmarshal tx %q", txRaw.String())
		}

		txs = append(txs, tx)
	}

	return txs, nil
}

func (c *Client) GetTransaction(
	ctx context.Context,
	acc ton.AccountID,
	lt uint64,
	hash ton.Bits256,
) (ton.Transaction, error) {
	txs, err := c.GetTransactions(ctx, 1, acc, lt, hash)
	if err != nil {
		return ton.Transaction{}, err
	}

	if len(txs) == 0 {
		return ton.Transaction{}, ErrNotFound
	}

	return txs[0], nil
}

// GetTransactionsSince returns all account transactions since the given logicalTime and hash (exclusive).
// The result is ordered from oldest to newest. Used to detect new txs to observe.
func (c *Client) GetTransactionsSince(
	ctx context.Context,
	acc ton.AccountID,
	oldestLT uint64,
	oldestHash ton.Bits256,
) (txs []ton.Transaction, err error) {
	lt, hash, err := c.getLastTransactionHash(ctx, acc)
	if err != nil {
		return nil, err
	}

	var result []ton.Transaction

	// reverse the result to get the oldest tx first
	defer func() {
		if len(result) > 0 {
			slices.Reverse(result)
		}
	}()

	for {
		hashBits := ton.Bits256(hash)

		// note that ton liteapi works in the reverse order.
		// Here we go from the LATEST txs to the oldest at N txs per page
		txs, err := c.GetTransactions(ctx, pageSize, acc, lt, hashBits)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to get transactions [lt %d, hash %s]", lt, hashBits.Hex())
		}

		if len(txs) == 0 {
			break
		}

		for i := range txs {
			found := txs[i].Lt == oldestLT && txs[i].Hash() == tlb.Bits256(oldestHash)
			if !found {
				continue
			}

			// early exit
			result = append(result, txs[:i]...)

			return result, nil
		}

		// otherwise, append all page results
		result = append(result, txs...)

		// prepare pagination params for the next page
		oldestIndex := len(txs) - 1

		lt, hash = txs[oldestIndex].PrevTransLt, txs[oldestIndex].PrevTransHash
	}

	return result, nil
}

// GetFirstTransaction scrolls through the transactions of the given account to find the first one.
// Note that it might fail w/o using an archival node. Also returns the number of
// scrolled transactions for this account i.e. total transactions
func (c *Client) GetFirstTransaction(ctx context.Context, acc ton.AccountID) (*ton.Transaction, int, error) {
	lt, hash, err := c.getLastTransactionHash(ctx, acc)
	if err != nil {
		return nil, 0, err
	}

	var (
		tx       *ton.Transaction
		scrolled int
	)

	for {
		hashBits := ton.Bits256(hash)

		txs, err := c.GetTransactions(ctx, pageSize, acc, lt, hashBits)
		if err != nil {
			return nil, scrolled, errors.Wrapf(err, "unable to get transactions [lt %d, hash %s]", lt, hashBits.Hex())
		}

		if len(txs) == 0 {
			break
		}

		scrolled += len(txs)

		tx = &txs[len(txs)-1]

		// Not we take the latest item in the list (oldest tx in the page)
		// and set it as the new last tx
		lt, hash = tx.PrevTransLt, tx.PrevTransHash
	}

	if tx == nil {
		return nil, scrolled, errors.Errorf("no transactions found [lt %d, hash %s]", lt, ton.Bits256(hash).Hex())
	}

	return tx, scrolled, nil
}

func (c *Client) SendMessage(ctx context.Context, payload []byte) (uint32, error) {
	const method = "sendBoc"

	params := map[string]any{
		"boc": base64.StdEncoding.EncodeToString(payload),
	}

	req := newRPCRequest(method, params)

	res, err := c.rpcRequest(ctx, req)
	if err != nil {
		return 0, errors.Wrapf(err, "%s: unable to call rpc with params: %v", method, req.Params)
	}

	// todo: future: this should be explored during e2e wiring,
	// todo: probably need to parse code from res.Result
	// #nosec G115 in range
	return uint32(res.Code), nil
}

func (c *Client) callAndUnmarshal(
	ctx context.Context,
	method string,
	params map[string]any,
	value any,
) error {
	resp, err := c.call(ctx, method, params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp, value); err != nil {
		return errors.Wrapf(err, "%s: unable to unmarshal rpc response (%s)", method, resp)
	}

	return nil
}

func (c *Client) call(ctx context.Context, method string, params map[string]any) (json.RawMessage, error) {
	req := newRPCRequest(method, params)

	res, err := c.rpcRequest(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: unable to call rpc with params: %v", method, req.Params)
	}

	if !res.Success {
		return nil, errors.Errorf(
			"%s: rpc call failed: %s (code: %d) with params: %v",
			method,
			res.Error,
			res.Code,
			req.Params,
		)
	}

	return res.Result, nil
}

// rpcRequest perform rpc request using HTTP transport
func (c *Client) rpcRequest(ctx context.Context, req rpcRequest) (rpcResponse, error) {
	httpReqBody, err := req.asBody()
	if err != nil {
		return rpcResponse{}, errors.Wrapf(err, "unable to marshal rpc request")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, httpReqBody)
	if err != nil {
		return rpcResponse{}, errors.Wrapf(err, "unable to create http request")
	}

	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return rpcResponse{}, errors.Wrap(err, "unable to send http request")
	}

	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return rpcResponse{}, errors.Wrap(err, "unable to read http response")
	}

	var rpcResp rpcResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return rpcResponse{}, errors.Wrap(err, "unable to unmarshal rpc response")
	}

	return rpcResp, nil
}
