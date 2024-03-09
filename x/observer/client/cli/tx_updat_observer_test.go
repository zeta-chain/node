package cli_test

import (
	"fmt"
	"testing"

	"github.com/zeta-chain/zetacore/x/observer/client/cli"
)

func Test_UpdateObserver(t *testing.T) {
	r, err := cli.ParseUpdateReason(0)
	fmt.Println(r, err)
}
