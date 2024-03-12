package filterdeposit

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/cmd/zetatool/config"
)

var btcCmd = &cobra.Command{
	Use:   "btc",
	Short: "Filter inbound btc deposits",
	Run:   FilterBTCTransactions,
}

func init() {
	Cmd.AddCommand(btcCmd)
}

func FilterBTCTransactions(cmd *cobra.Command, _ []string) {
	configFile, err := cmd.Flags().GetString(config.Flag)
	fmt.Println("config file name: ", configFile)
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.GetConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	list := getHashList(cfg)
	CheckForCCTX(list, cfg)
}

func getHashList(cfg *config.Config) []Deposit {
	var list []Deposit
	lastHash := ""

	url := cfg.BtcExplorer

	for {
		nextQuery := url
		if lastHash != "" {
			path := fmt.Sprintf("/chain/%s", lastHash)
			nextQuery = url + path
		}
		// #nosec G107 url must be variable
		res, getErr := http.Get(nextQuery)
		if getErr != nil {
			log.Fatal(getErr)
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
		closeErr := res.Body.Close()
		if closeErr != nil {
			log.Fatal(closeErr)
		}

		var txns []map[string]interface{}
		err := json.Unmarshal(body, &txns)
		if err != nil {
			fmt.Println("error unmarshalling: ", err.Error())
		}

		if len(txns) == 0 {
			break
		}

		fmt.Println("Length of txns: ", len(txns))

		for _, txn := range txns {
			hash := txn["txid"].(string)

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
				if bytes.Equal(memoBytes, []byte(DonationMessage)) {
					continue
				}
			} else {
				continue
			}

			//Make sure Deposit is sent to correct tss address
			if strings.Compare("0014", scriptpubkey[:4]) == 0 && targetAddr == cfg.TssAddressBTC {
				entry := Deposit{
					hash,
					vout0["value"].(float64),
				}
				list = append(list, entry)
			}
		}

		lastTxn := txns[len(txns)-1]
		lastHash = lastTxn["txid"].(string)
		//fmt.Println("last hash: ", lastHash)
	}

	return list
}
