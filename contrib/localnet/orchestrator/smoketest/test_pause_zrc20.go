package main

import (
	"fmt"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	"time"
)

func (sm *SmokeTest) TestPauseZRC20() {
	LoudPrintf("Test ZRC20 pause\n")
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()

	// Pause USDT
	msg := fungibletypes.NewMsgUpdateZRC20PausedStatus(
		FungibleAdminAddress,
		[]string{USDTZRC20Addr},
		fungibletypes.UpdatePausedStatusAction_PAUSE,
	)
	res, err := sm.zetaTxServer.BroadcastTx(FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("pause zrc20 tx: %s\n", res.String())
}
