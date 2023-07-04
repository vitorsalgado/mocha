package dzstd

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLatency(t *testing.T) {
	d := 5 * time.Second
	r := Latency(d)()

	require.Equal(t, d, r)
}

func TestUniformDistribution(t *testing.T) {
	min := int64(100)
	max := int64(300)

	duration := UniformDistribution(min, max)()

	require.LessOrEqual(t, duration.Milliseconds(), max)
	require.GreaterOrEqual(t, duration.Milliseconds(), min)
}

func TestJitter(t *testing.T) {
	deviation := float64(15)
	jitter := Jitter(300*time.Millisecond, deviation)

	deviations := make([]int64, 10)
	durations := make([]time.Duration, 10)
	durations[0] = jitter()
	durations[1] = jitter()
	durations[2] = jitter()
	durations[3] = jitter()
	durations[4] = jitter()
	durations[5] = jitter()
	durations[6] = jitter()
	durations[7] = jitter()
	durations[8] = jitter()
	durations[9] = jitter()

	total := time.Duration(0)
	for _, d := range durations {
		total += d
	}

	mean := total / 10
	totalDeviation := int64(0)

	for i, d := range durations {
		dev := math.Abs(float64(d - mean))
		deviations[i] = int64(dev)
		totalDeviation += int64(dev)
	}

	avgDeviation := totalDeviation / 10
	deviationPercentage := (float64(avgDeviation) / float64(mean)) * 100

	require.LessOrEqual(t, deviationPercentage, deviation)
}
