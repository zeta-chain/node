package filterdeposit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/cmd/zetatool/config"
)

var Cmd = &cobra.Command{
	Use:   "filterdeposit",
	Short: "filter missing inbound deposits",
}

// Deposit is a data structure for keeping track of inbound transactions
type Deposit struct {
	TxID   string
	Amount uint64
}

// CheckForCCTX is querying zeta core for a cctx associated with a confirmed transaction hash. If the cctx is not found,
// then the transaction hash is added to the list of missed inbound transactions.
func CheckForCCTX(list []Deposit, cfg *config.Config) []Deposit {
	var missedList []Deposit

	fmt.Println("Going through list, num of transactions: ", len(list))
	for _, entry := range list {
		zetaURL, err := url.JoinPath(cfg.ZetaURL, "zeta-chain", "crosschain", "in_tx_hash_to_cctx_data", entry.TxID)
		if err != nil {
			log.Fatal(err)
		}

		request, err := http.NewRequest(http.MethodGet, zetaURL, nil)
		if err != nil {
			log.Fatal(err)
		}
		request.Header.Add("Accept", "application/json")
		client := &http.Client{}

		response, getErr := client.Do(request)
		if getErr != nil {
			log.Fatal(getErr)
		}

		data, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
		closeErr := response.Body.Close()
		if closeErr != nil {
			log.Fatal(closeErr)
		}

		var cctx map[string]interface{}
		err = json.Unmarshal(data, &cctx)
		if err != nil {
			fmt.Println("error unmarshalling: ", err.Error())
		}

		// successful query of the given cctx will not contain a "message" field with value "not found", if it was not
		// found then it is added to the missing list.
		if _, ok := cctx["message"]; ok {
			if strings.Compare(cctx["message"].(string), "not found") == 0 {
				missedList = append(missedList, entry)
			}
		}
	}

	fmt.Printf("Found %d missed transactions.\n", len(missedList))
	for _, entry := range missedList {
		fmt.Printf("%s, amount: %d\n", entry.TxID, entry.Amount)
	}
	return missedList
}
