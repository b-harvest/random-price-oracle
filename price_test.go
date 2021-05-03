package main

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomPriceChangeRatio(t *testing.T) {
	const (
		iterations       = 20000
		startingPrice    = 10.0
		avgDiffThreshold = 0.01
	)
	sum := 0.0
	for seed := int64(0); seed < iterations; seed++ {
		r := rand.New(rand.NewSource(seed))
		currentPrice := startingPrice
		for i := 0; i < 86400/5; i++ {
			changeRatio := RandomPriceChangeRatio(r)
			currentPrice *= changeRatio
		}
		sum += currentPrice
	}
	avgPrice := sum / iterations
	avgDiff := math.Abs(avgPrice/startingPrice - 1)
	require.Truef(t, avgDiff < avgDiffThreshold, "average price diff %f is larger than threshold %f", avgDiff, avgDiffThreshold)
}
