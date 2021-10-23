package metaclientd

import (
	"testing"
)

func TestQueryRouter(t *testing.T) {
	mo := &MetaObserver{}

	mo.WatchRouter("BSC")
}
