package infra

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/buger/jsonparser"
	"github.com/zeta-chain/zetacore/zetaclient/btc"
	"github.com/zeta-chain/zetacore/zetaclient/btc/model"
)

var _ btc.Client = (*JSONRpcClient)(nil)

const verbosity = 2

type JSONRpcClient struct {
	client        *http.Client
	endpoint      string
	targetAddress string
}

func NewJSONRpcClient(endpoint, targetAddress string) *JSONRpcClient {
	return &JSONRpcClient{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		endpoint:      endpoint,
		targetAddress: targetAddress,
	}
}

func (cli *JSONRpcClient) GetBlockHeight() (int64, error) {
	reqBody := []byte(`{"jsonrpc": "1.0", "id": "zeta", "method": "getblockcount", "params": []}`)
	req, err := http.NewRequest(http.MethodPost, cli.endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.client.Do(req)
	if err != nil {
		return 0, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	return jsonparser.GetInt(data, "result")
}

func (cli *JSONRpcClient) GetBlockHash(block int64) (string, error) {
	var hash string
	reqBody := []byte(fmt.Sprintf(`{"jsonrpc": "1.0", "id": "zeta", "method": "getblockhash", "params": [%d]}`, block))
	req, err := http.NewRequest(http.MethodPost, cli.endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return hash, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.client.Do(req)
	if err != nil {
		return hash, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return hash, err
	}
	return jsonparser.GetString(data, "result")
}

func (cli *JSONRpcClient) GetEventsByHash(hash string) ([]*model.RawEvent, error) {
	reqBody := []byte(fmt.Sprintf(`{"jsonrpc": "1.0", "id": "zeta", "method": "getblock", "params": ["%s",%d]}`, hash, verbosity))
	req, err := http.NewRequest(http.MethodPost, cli.endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.client.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var events []*model.RawEvent
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		filter := false
		jsonparser.ArrayEach(value, func(value2 []byte, dataType jsonparser.ValueType, offset int, err error) {
			asm, _ := jsonparser.GetString(value2, "scriptPubKey", "asm")
			amount, _ := jsonparser.GetFloat(value2, "value")
			var addresses []string
			var amounts []float64
			jsonparser.ArrayEach(value2, func(addressBytes []byte, dataType jsonparser.ValueType, offset int, err error) {
				address := string(addressBytes)
				if address == cli.targetAddress {
					filter = true
				} else {
					filter = false
				}
				addresses = append(addresses, address)
				amounts = append(amounts, amount)
			}, "scriptPubKey", "addresses")
			if filter {
				events = append(events, &model.RawEvent{
					ASM:       asm,
					Addresses: addresses,
					Values:    amounts,
				})
			}
		}, "vout")
	}, "result", "tx")
	return events, nil
}
