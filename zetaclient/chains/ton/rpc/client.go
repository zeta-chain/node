package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	client   *http.Client
	endpoint string
}

type Opt func(c *Client)

func WithHTTPClient(client *http.Client) Opt {
	return func(c *Client) { c.client = client }
}

// New Client constructor
// https://docs.ton.org/v3/guidelines/dapps/apis-sdks/ton-http-apis
func New(endpoint string, opts ...Opt) *Client {
	const defaultTimeout = 10 * time.Second

	// See: https://toncenter.com/api/v2
	//
	// Most API providers expose a url with api in in the path
	// - https://ton-testnet.core.chainstack.com/$key/api/v2
	// - https://$node.ton-mainnet.quiknode.pro/$key/
	//
	// And we need to add /jsonRPC to the end of the url
	endpoint = strings.TrimRight(endpoint, "/")
	endpoint += "/jsonRPC"

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

func (c *Client) callAndUnmarshal(
	ctx context.Context,
	method string,
	params map[string]any,
	value json.Unmarshaler,
) error {
	resp, err := c.call(ctx, method, params)
	if err != nil {
		return err
	}

	if err := value.UnmarshalJSON(resp); err != nil {
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

type rpcRequest struct {
	Jsonrpc string         `json:"jsonrpc"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
	ID      string         `json:"id"`
}

func newRPCRequest(method string, params map[string]any) rpcRequest {
	if params == nil {
		params = make(map[string]any)
	}

	return rpcRequest{
		Jsonrpc: "2.0",
		ID:      "1",
		Method:  method,
		Params:  params,
	}
}

func (r *rpcRequest) asBody() (io.Reader, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal rpc request")
	}

	return bytes.NewReader(body), nil
}

type rpcResponse struct {
	Success bool            `json:"ok"`
	Result  json.RawMessage `json:"result"`
	Error   string          `json:"error"`
	Code    int             `json:"code"`
}
