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

type deposit struct {
	txId   string
	amount float64
}

func CheckForCCTX(list []deposit, cfg *config.Config) {
	var missedList []deposit

	fmt.Println("Going through list, num of transactions: ", len(list))
	for _, entry := range list {
		zetaUrl, err := url.JoinPath(cfg.ZetaUrl, "zeta-chain", "crosschain", "in_tx_hash_to_cctx_data", entry.txId)
		if err != nil {
			log.Fatal(err)
		}
		res, getErr := http.Get(zetaUrl)
		if getErr != nil {
			log.Fatal(getErr)
		}

		data, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
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
		fmt.Printf("%s, amount: %d\n", entry.txId, int64(entry.amount))
	}
}
