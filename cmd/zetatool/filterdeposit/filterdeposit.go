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

type Deposit struct {
	TxID   string
	Amount uint64
}

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

		if _, ok := cctx["code"]; ok {
			missedList = append(missedList, entry)
			//fmt.Println("appending to missed list: ", entry)
		}
	}

	for _, entry := range missedList {
		fmt.Printf("%s, amount: %d\n", entry.TxID, entry.Amount)
	}
}
