package filterdeposit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

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
func CheckForCCTX(list []Deposit, cfg *config.Config) {
	var missedList []Deposit

	fmt.Println("Going through list, num of transactions: ", len(list))
	for _, entry := range list {
		zetaURL, err := url.JoinPath(cfg.ZetaURL, "zeta-chain", "crosschain", "in_tx_hash_to_cctx_data", entry.TxID)
		if err != nil {
			log.Fatal(err)
		}
		// #nosec G107 url must be variable
		res, getErr := http.Get(zetaURL)
		if getErr != nil {
			log.Fatal(getErr)
		}

		data, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
		closeErr := res.Body.Close()
		if closeErr != nil {
			log.Fatal(closeErr)
		}

		var cctx map[string]interface{}
		err = json.Unmarshal(data, &cctx)
		if err != nil {
			fmt.Println("error unmarshalling: ", err.Error())
		}

		// successful query of the given cctx will not contain a "code" field, therefore if it exists then the cctx
		// was not found and is added to the missing list.
		if _, ok := cctx["code"]; ok {
			missedList = append(missedList, entry)
		}
	}

	for _, entry := range missedList {
		fmt.Printf("%s, amount: %d\n", entry.TxID, entry.Amount)
	}
}
