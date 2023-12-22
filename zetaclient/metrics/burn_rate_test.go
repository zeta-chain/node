package metrics

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type BurnRateSuite struct {
	suite.Suite
	burnRate *BurnRate
}

func TestBurnRateSuite(t *testing.T) {
	suite.Run(t, new(BurnRateSuite))
}

func (brs *BurnRateSuite) SetupTest() {
	brs.burnRate = NewBurnRate(10)
}

func (brs *BurnRateSuite) TestBurnRate() {
	// 1. Add number of entries less than window size and verify rate
	for i := 1; i <= 5; i++ {
		err := brs.burnRate.AddFee(int64(i)*10, int64(i))
		brs.NoError(err)
	}

	rate := brs.burnRate.GetBurnRate()
	brs.Equal(int64(30), rate.Int64())

	// 2. Add number of entries equal to window size and verify rate
	for i := 6; i <= 10; i++ {
		err := brs.burnRate.AddFee(int64(i)*10, int64(i))
		brs.NoError(err)
	}

	rate = brs.burnRate.GetBurnRate()
	brs.Equal(int64(55), rate.Int64())

	// 3. Add number of entries greater than window size and verify rate
	for i := 11; i <= 15; i++ {
		err := brs.burnRate.AddFee(int64(i)*10, int64(i))
		brs.NoError(err)
	}

	rate = brs.burnRate.GetBurnRate()
	brs.Equal(int64(105), rate.Int64())

	// 4. Add non-contiguous block entries to queue and verify rate
	entries := []int{16, 18, 22, 23, 24, 29}
	for _, val := range entries {
		err := brs.burnRate.AddFee(int64(val)*10, int64(val))
		brs.NoError(err)
	}

	rate = brs.burnRate.GetBurnRate()
	brs.Equal(int64(98), rate.Int64())
}
