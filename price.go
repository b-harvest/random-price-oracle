package main

import (
	"math"
	"math/rand"
)

func RandomPriceChangeRatio(r *rand.Rand) float64 {
	const (
		vol                = 1.0
		updateInterval     = 5.0
		priceJumpPerDay    = 5.0
		priceJumpMagnitude = 0.01
	)
	ts := updateInterval / (365 * 24 * 60 * 60)
	if r.Float64()*2 < ts*365*priceJumpPerDay {
		if r.Float64() > 0.5 {
			return 1 + priceJumpMagnitude
		}
		return 1 - priceJumpMagnitude
	}
	return 1 + r.NormFloat64()*vol*math.Sqrt(ts)
}
