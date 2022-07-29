package cmd

import (
	"fmt"
	"math"
	"testing"
)

func Test_syncPools(t *testing.T) {
	from := 1
	to := 260
	limit := 50

	batches := int(math.Ceil(float64(to) / float64(limit)))

	fromBlock := from
	toBlock := fromBlock + limit

	for i := 1; i <= batches; i++ {
		if fromBlock > toBlock {
			continue
		}
		fmt.Printf("Scanning blocks from %d to %d, batch %d/%d\n", fromBlock, toBlock, i, batches)

		fromBlock = toBlock + 1
		toBlock = fromBlock + limit
		if toBlock > to {
			toBlock = to
		}
	}
}
