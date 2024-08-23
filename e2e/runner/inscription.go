package runner

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
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
	To      string `json:"to"`
}

// InscriptionBuilder is a util struct that help create inscription commit and reveal transactions
type InscriptionBuilder struct {
	sidecarURL string
	client     http.Client
}

// GenerateCommitAddress generates a commit p2tr address that one can send funds to this address
func (r *InscriptionBuilder) GenerateCommitAddress(memo []byte) (string, error) {
	// Create the payload
	postData := map[string]string{
		"memo": hex.EncodeToString(memo),
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return "", err
	}

	postURL := r.sidecarURL + "/commit"
	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", errors.Wrap(err, "cannot create commit request")
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := r.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "cannot send to sidecar")
	}
	defer resp.Body.Close()

	// Read the response body
	var response commitResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	fmt.Print("raw commit response ", response.Address)

	return response.Address, nil
}

// GenerateRevealTxn creates the corresponding reveal txn to the commit txn.
func (r *InscriptionBuilder) GenerateRevealTxn(to string, txnHash string, idx int, amount float64) (string, error) {
	postData := revealRequest{
		Txn:     txnHash,
		Idx:     idx,
		Amount:  int(amount * 100000000),
		FeeRate: 10,
		To:      to,
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return "", err
	}

	postURL := r.sidecarURL + "/reveal"
	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(jsonData))
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
