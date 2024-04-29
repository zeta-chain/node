package metrics

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
)

// BurnRate calculates the average burn rate for a range of blocks.
type BurnRate struct {
	blockLow   int64
	blockHigh  int64
	windowSize int64
	total      int64
	queue      []int64
}

// NewBurnRate creates a new BurnRate instance with a window size.
func NewBurnRate(windowSize int64) *BurnRate {
	return &BurnRate{
		blockLow:   1,
		blockHigh:  1,
		windowSize: windowSize,
		total:      0,
		queue:      make([]int64, 1),
	}
}

// AddFee adds fee amount spent on a tx for a particular block. It is added to a queue which is used to calculate
// the average burn rate for a range of blocks determined by the window size.
func (br *BurnRate) AddFee(amount int64, block int64) error {
	// Check if block is in range of the window
	if block < br.blockLow {
		return fmt.Errorf("block out of range %d", block)
	}

	// If block is greater than blockHigh, shift the window up
	if block > br.blockHigh {
		err := br.enqueueEntry(block, amount)
		if err != nil {
			return err
		}
		br.blockHigh = block

		if br.blockHigh-br.blockLow >= br.windowSize {
			// Remove oldest block(s) from queue
			err = br.dequeueOldEntries()
			if err != nil {
				return err
			}
		}
	} else {
		// Add amount to existing entry in queue
		index := block - br.blockLow
		br.queue[index] += amount
		br.total += amount
	}

	return nil
}

// enqueueEntry adds fee entry into queue if is in range of the window. A padding is added if the block height is
// more than one block greater than the highest range.
func (br *BurnRate) enqueueEntry(block int64, amount int64) error {
	diff := block - br.blockHigh
	if diff < 1 {
		return fmt.Errorf("enqueueEntry: block difference is too low: %d", diff)
	}
	// if difference in block num is greater than 1, need to pad the queue
	if diff > 1 {
		for i := int64(0); i < diff-1; i++ {
			br.queue = append(br.queue, 0)
		}
	}

	// Adjust total with new entry
	br.total += amount

	// enqueue latest entry
	br.queue = append(br.queue, amount)

	return nil
}

// dequeueOldEntries dequeues old entries
// when the window slides forward, older entries in the queue need to be cleared.
func (br *BurnRate) dequeueOldEntries() error {
	diff := br.blockHigh - br.blockLow
	if diff < br.windowSize {
		return fmt.Errorf("dequeueOldEntries: queue is less than or equal to window size, no need to dequeue")
	}
	dequeueLen := diff - (br.windowSize - 1)

	// Adjust total with dequeued elements
	for i := int64(0); i < dequeueLen; i++ {
		br.total -= br.queue[i]
	}

	// dequeue old entries
	br.queue = br.queue[dequeueLen:]
	br.blockLow += dequeueLen

	return nil
}

// GetBurnRate calculates current burn rate and return the value.
func (br *BurnRate) GetBurnRate() sdkmath.Int {
	if br.blockHigh < br.windowSize {
		return sdkmath.NewInt(br.total).QuoRaw(br.blockHigh)
	}
	return sdkmath.NewInt(br.total).QuoRaw(br.windowSize)
}
