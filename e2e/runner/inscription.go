package runner

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type commitResponse struct {
	Address string `json:"address"`
}

type revealResponse struct {
	RawHex string `json:"rawHex"`
}

type revealRequest struct {
	Txn     string `json:"txn"`
	Idx     int    `json:"idx"`
	Amount  int    `json:"amount"`
	FeeRate int    `json:"feeRate"`
}

// InscriptionBuilder is a util struct that help create inscription commit and reveal transactions
type InscriptionBuilder struct {
	sidecarUrl string
	client     http.Client
}

func (r *InscriptionBuilder) GenerateCommitAddress(memo []byte) (btcutil.Address, error) {
	// Create the payload
	postData := map[string]string{
		"memo": hex.EncodeToString(memo),
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return nil, err
	}

	postUrl := r.sidecarUrl + "/commit"
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrap(err, "cannot create commit request")
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "cannot send to sidecar")
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read commit response body")
	}

	// Parse the JSON response
	var response commitResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.Wrap(err, "cannot parse commit response body")
	}

	return btcutil.DecodeAddress(response.Address, &chaincfg.RegressionNetParams)
}

func (r *InscriptionBuilder) GenerateRevealTxn(txnHash string, idx int, amount float64) (string, error) {
	postData := revealRequest{
		Txn:     txnHash,
		Idx:     idx,
		Amount:  int(amount * 100000000),
		FeeRate: 10,
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return "", err
	}

	postUrl := r.sidecarUrl + "/reveal"
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", errors.Wrap(err, "cannot create reveal request")
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := r.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "cannot send reveal to sidecar")
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "cannot read reveal response body")
	}

	// Parse the JSON response
	var response revealResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", errors.Wrap(err, "cannot parse reveal response body")
	}

	// Access the "address" field
	return response.RawHex, nil
}
