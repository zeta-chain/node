package filterdeposit

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zeta-chain/zetacore/cmd/zetatool/config"
	"github.com/zeta-chain/zetacore/pkg/constant"
)

func NewBtcCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "btc",
		Short: "Filter inbound btc deposits",
		RunE:  FilterBTCTransactions,
	}
}

// FilterBTCTransactions is a command that queries the bitcoin explorer for inbound transactions that qualify for
// cross chain transactions.
func FilterBTCTransactions(cmd *cobra.Command, _ []string) error {
	configFile, err := cmd.Flags().GetString(config.FlagConfig)
	fmt.Println("config file name: ", configFile)
	if err != nil {
		return err
	}
	btcChainID, err := cmd.Flags().GetString(BTCChainIDFlag)
	if err != nil {
		return err
	}
	cfg, err := config.GetConfig(configFile)
	if err != nil {
		return err
	}
	fmt.Println("getting tss Address")
	res, err := GetTssAddress(cfg, btcChainID)
	if err != nil {
		return err
	}
	fmt.Println("got tss Address")
	list, err := getHashList(cfg, res.Btc)
	if err != nil {
		return err
	}

	_, err = CheckForCCTX(list, cfg)
	return err
}

// getHashList is called by FilterBTCTransactions to help query and filter inbound transactions on btc
func getHashList(cfg *config.Config, tssAddress string) ([]Deposit, error) {
	var list []Deposit
	lastHash := ""

	// Setup URL for query
	btcURL, err := url.JoinPath(cfg.BtcExplorerURL, "address", tssAddress, "txs")
	if err != nil {
		return list, err
	}

	// This loop will query the bitcoin explorer for transactions associated with the TSS address. Since the api only
	// allows a response of 25 transactions per request, several requests will be required in order to retrieve a
	// complete list.
	for {
		// The Next Query is determined by the last transaction hash provided by the previous response.
		nextQuery := btcURL
		if lastHash != "" {
			nextQuery, err = url.JoinPath(btcURL, "chain", lastHash)
			if err != nil {
				return list, err
			}
		}
		// #nosec G107 url must be variable
		res, getErr := http.Get(nextQuery)
		if getErr != nil {
			return list, getErr
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			return list, readErr
		}
		closeErr := res.Body.Close()
		if closeErr != nil {
			return list, closeErr
		}

		// NOTE: decoding json from request dynamically is not ideal, however there isn't a detailed, defined data structure
		// provided by blockstream. Will need to create one in the future using following definition:
		// https://github.com/Blockstream/esplora/blob/master/API.md#transaction-format
		var txns []map[string]interface{}
		err := json.Unmarshal(body, &txns)
		if err != nil {
			return list, err
		}

		if len(txns) == 0 {
			break
		}

		fmt.Println("Length of txns: ", len(txns))

		// The "/address" blockstream api provides a maximum of 25 transactions associated with a given address. This
		// loop will iterate over that list of transactions to determine whether each transaction can be considered
		// a deposit to ZetaChain.
		for _, txn := range txns {
			// Get tx hash of the current transaction
			hash := txn["txid"].(string)

			// Read the first output of the transaction and parse the destination address.
			// This address should be the TSS address.
			vout := txn["vout"].([]interface{})
			vout0 := vout[0].(map[string]interface{})
			var vout1 map[string]interface{}
			if len(vout) > 1 {
				vout1 = vout[1].(map[string]interface{})
			} else {
				continue
			}
			_, found := vout0["scriptpubkey"]
			scriptpubkey := ""
			if found {
				scriptpubkey = vout0["scriptpubkey"].(string)
			}
			_, found = vout0["scriptpubkey_address"]
			targetAddr := ""
			if found {
				targetAddr = vout0["scriptpubkey_address"].(string)
			}

			//Check if txn is confirmed
			status := txn["status"].(map[string]interface{})
			confirmed := status["confirmed"].(bool)
			if !confirmed {
				continue
			}

			//Filter out deposits less than min base fee
			if vout0["value"].(float64) < 1360 {
				continue
			}

			//Check if Deposit is a donation
			scriptpubkey1 := vout1["scriptpubkey"].(string)
			if len(scriptpubkey1) >= 4 && scriptpubkey1[:2] == "6a" {
				memoSize, err := strconv.ParseInt(scriptpubkey1[2:4], 16, 32)
				if err != nil {
					continue
				}
				if int(memoSize) != (len(scriptpubkey1)-4)/2 {
					continue
				}
				memoBytes, err := hex.DecodeString(scriptpubkey1[4:])
				if err != nil {
					continue
				}
				if bytes.Equal(memoBytes, []byte(constant.DonationMessage)) {
					continue
				}
			} else {
				continue
			}

			//Make sure Deposit is sent to correct tss address
			if strings.Compare("0014", scriptpubkey[:4]) == 0 && targetAddr == tssAddress {
				entry := Deposit{
					hash,
					// #nosec G115 parsing json requires float64 type from blockstream
					uint64(vout0["value"].(float64)),
				}
				list = append(list, entry)
			}
		}

		lastTxn := txns[len(txns)-1]
		lastHash = lastTxn["txid"].(string)
	}

	return list, nil
}
