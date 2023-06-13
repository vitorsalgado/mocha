package lib

import (
	"math/rand"
	"time"
)

type Delay func() time.Duration

func Latency(latency time.Duration) Delay {
	return func() time.Duration { return latency }
}

func UniformDistribution(min, max int64) Delay {
	return func() time.Duration {
		return time.Duration(rand.Int63n(max-min+1)+min) * time.Millisecond
	}
}

func Jitter(avg time.Duration, deviation float64) Delay {
	maxDeviation := time.Duration(float64(avg) * (deviation / 100.0))

	return func() time.Duration {
		jitter := time.Duration(rand.Int63n(int64(maxDeviation)))

		return avg + jitter - (maxDeviation / 2)
	}
}
