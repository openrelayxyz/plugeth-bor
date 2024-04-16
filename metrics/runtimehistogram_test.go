package metrics

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"reflect"
	"runtime/metrics"
	"testing"
	"time"
)

var _ Histogram = (*runtimeHistogram)(nil)

type runtimeHistogramTest struct {
	h metrics.Float64Histogram

	Count       int64
	Min         int64
	Max         int64
	Sum         int64
	Mean        float64
	Variance    float64
	StdDev      float64
	Percentiles []float64 // .5 .8 .9 .99 .995
}

// This test checks the results of statistical functions implemented
// by runtimeHistogramSnapshot.
func TestRuntimeHistogramStats(t *testing.T) {
	t.Parallel()

	tests := []runtimeHistogramTest{
		0: {
			h: metrics.Float64Histogram{
				Counts:  []uint64{},
				Buckets: []float64{},
			},
			Count:       0,
			Max:         0,
			Min:         0,
			Sum:         0,
			Mean:        0,
			Variance:    0,
			StdDev:      0,
			Percentiles: []float64{0, 0, 0, 0, 0},
		},
		1: {
			// This checks the case where the highest bucket is +Inf.
			h: metrics.Float64Histogram{
				Counts:  []uint64{0, 1, 2},
				Buckets: []float64{0, 0.5, 1, math.Inf(1)},
			},
			Count:       3,
			Max:         1,
			Min:         0,
			Sum:         3,
			Mean:        0.9166666,
			Percentiles: []float64{1, 1, 1, 1, 1},
			Variance:    0.020833,
			StdDev:      0.144433,
		},
		2: {
			h: metrics.Float64Histogram{
				Counts:  []uint64{8, 6, 3, 1},
				Buckets: []float64{12, 16, 18, 24, 25},
			},
			Count:       18,
			Max:         25,
			Min:         12,
			Sum:         270,
			Mean:        16.75,
			Variance:    10.3015,
			StdDev:      3.2096,
			Percentiles: []float64{16, 18, 18, 24, 24},
		},
	}

	for i, test := range tests {
		i, test := i, test

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			s := RuntimeHistogramFromData(1.0, &test.h).Snapshot()

			if v := s.Count(); v != test.Count {
				t.Errorf("Count() = %v, want %v", v, test.Count)
			}

			if v := s.Min(); v != test.Min {
				t.Errorf("Min() = %v, want %v", v, test.Min)
			}

			if v := s.Max(); v != test.Max {
				t.Errorf("Max() = %v, want %v", v, test.Max)
			}

			if v := s.Sum(); v != test.Sum {
				t.Errorf("Sum() = %v, want %v", v, test.Sum)
			}

			if v := s.Mean(); !approxEqual(v, test.Mean, 0.0001) {
				t.Errorf("Mean() = %v, want %v", v, test.Mean)
			}

			if v := s.Variance(); !approxEqual(v, test.Variance, 0.0001) {
				t.Errorf("Variance() = %v, want %v", v, test.Variance)
			}

			if v := s.StdDev(); !approxEqual(v, test.StdDev, 0.0001) {
				t.Errorf("StdDev() = %v, want %v", v, test.StdDev)
			}

			ps := []float64{.5, .8, .9, .99, .995}
			if v := s.Percentiles(ps); !reflect.DeepEqual(v, test.Percentiles) {
				t.Errorf("Percentiles(%v) = %v, want %v", ps, v, test.Percentiles)
			}
		})
	}
}

func approxEqual(x, y, ε float64) bool {
	if math.IsInf(x, -1) && math.IsInf(y, -1) {
		return true
	}

	if math.IsInf(x, 1) && math.IsInf(y, 1) {
		return true
	}

	if math.IsNaN(x) && math.IsNaN(y) {
		return true
	}

	return math.Abs(x-y) < ε
}

// This test verifies that requesting Percentiles in unsorted order
// returns them in the requested order.
func TestRuntimeHistogramStatsPercentileOrder(t *testing.T) {
	s := RuntimeHistogramFromData(1.0, &metrics.Float64Histogram{
		Counts:  []uint64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		Buckets: []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}).Snapshot()
	result := s.Percentiles([]float64{1, 0.2, 0.5, 0.1, 0.2})
	expected := []float64{10, 2, 5, 1, 2}

	if !reflect.DeepEqual(result, expected) {
		t.Fatal("wrong result:", result)
	}
}

func BenchmarkRuntimeHistogramSnapshotRead(b *testing.B) {
	var sLatency = "7\xff\x81\x03\x01\x01\x10Float64Histogram\x01\xff\x82\x00\x01\x02\x01\x06Counts\x01\xff\x84\x00\x01\aBuckets\x01\xff\x86\x00\x00\x00\x16\xff\x83\x02\x01\x01\b[]uint64\x01\xff\x84\x00\x01\x06\x00\x00\x17\xff\x85\x02\x01\x01\t[]float64\x01\xff\x86\x00\x01\b\x00\x00\xfe\x06T\xff\x82\x01\xff\xa2\x00\xfe\r\xef\x00\x01\x02\x02\x04\x05\x04\b\x15\x17 B?6.L;$!2) \x1a? \x190aH7FY6#\x190\x1d\x14\x10\x1b\r\t\x04\x03\x01\x01\x00\x03\x02\x00\x03\x05\x05\x02\x02\x06\x04\v\x06\n\x15\x18\x13'&.\x12=H/L&\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\xff\xa3\xfe\xf0\xff\x00\xf8\x95\xd6&\xe8\v.q>\xf8\x95\xd6&\xe8\v.\x81>\xf8\xdfA:\xdc\x11ŉ>\xf8\x95\xd6&\xe8\v.\x91>\xf8:\x8c0\xe2\x8ey\x95>\xf8\xdfA:\xdc\x11ř>\xf8\x84\xf7C֔\x10\x9e>\xf8\x95\xd6&\xe8\v.\xa1>\xf8:\x8c0\xe2\x8ey\xa5>\xf8\xdfA:\xdc\x11ũ>\xf8\x84\xf7C֔\x10\xae>\xf8\x95\xd6&\xe8\v.\xb1>\xf8:\x8c0\xe2\x8ey\xb5>\xf8\xdfA:\xdc\x11Ź>\xf8\x84\xf7C֔\x10\xbe>\xf8\x95\xd6&\xe8\v.\xc1>\xf8:\x8c0\xe2\x8ey\xc5>\xf8\xdfA:\xdc\x11\xc5\xc9>\xf8\x84\xf7C֔\x10\xce>\xf8\x95\xd6&\xe8\v.\xd1>\xf8:\x8c0\xe2\x8ey\xd5>\xf8\xdfA:\xdc\x11\xc5\xd9>\xf8\x84\xf7C֔\x10\xde>\xf8\x95\xd6&\xe8\v.\xe1>\xf8:\x8c0\xe2\x8ey\xe5>\xf8\xdfA:\xdc\x11\xc5\xe9>\xf8\x84\xf7C֔\x10\xee>\xf8\x95\xd6&\xe8\v.\xf1>\xf8:\x8c0\xe2\x8ey\xf5>\xf8\xdfA:\xdc\x11\xc5\xf9>\xf8\x84\xf7C֔\x10\xfe>\xf8\x95\xd6&\xe8\v.\x01?\xf8:\x8c0\xe2\x8ey\x05?\xf8\xdfA:\xdc\x11\xc5\t?\xf8\x84\xf7C֔\x10\x0e?\xf8\x95\xd6&\xe8\v.\x11?\xf8:\x8c0\xe2\x8ey\x15?\xf8\xdfA:\xdc\x11\xc5\x19?\xf8\x84\xf7C֔\x10\x1e?\xf8\x95\xd6&\xe8\v.!?\xf8:\x8c0\xe2\x8ey%?\xf8\xdfA:\xdc\x11\xc5)?\xf8\x84\xf7C֔\x10.?\xf8\x95\xd6&\xe8\v.1?\xf8:\x8c0\xe2\x8ey5?\xf8\xdfA:\xdc\x11\xc59?\xf8\x84\xf7C֔\x10>?\xf8\x95\xd6&\xe8\v.A?\xf8:\x8c0\xe2\x8eyE?\xf8\xdfA:\xdc\x11\xc5I?\xf8\x84\xf7C֔\x10N?\xf8\x95\xd6&\xe8\v.Q?\xf8:\x8c0\xe2\x8eyU?\xf8\xdfA:\xdc\x11\xc5Y?\xf8\x84\xf7C֔\x10^?\xf8\x95\xd6&\xe8\v.a?\xf8:\x8c0\xe2\x8eye?\xf8\xdfA:\xdc\x11\xc5i?\xf8\x84\xf7C֔\x10n?\xf8\x95\xd6&\xe8\v.q?\xf8:\x8c0\xe2\x8eyu?\xf8\xdfA:\xdc\x11\xc5y?\xf8\x84\xf7C֔\x10~?\xf8\x95\xd6&\xe8\v.\x81?\xf8:\x8c0\xe2\x8ey\x85?\xf8\xdfA:\xdc\x11ŉ?\xf8\x84\xf7C֔\x10\x8e?\xf8\x95\xd6&\xe8\v.\x91?\xf8:\x8c0\xe2\x8ey\x95?\xf8\xdfA:\xdc\x11ř?\xf8\x84\xf7C֔\x10\x9e?\xf8\x95\xd6&\xe8\v.\xa1?\xf8:\x8c0\xe2\x8ey\xa5?\xf8\xdfA:\xdc\x11ũ?\xf8\x84\xf7C֔\x10\xae?\xf8\x95\xd6&\xe8\v.\xb1?\xf8:\x8c0\xe2\x8ey\xb5?\xf8\xdfA:\xdc\x11Ź?\xf8\x84\xf7C֔\x10\xbe?\xf8\x95\xd6&\xe8\v.\xc1?\xf8:\x8c0\xe2\x8ey\xc5?\xf8\xdfA:\xdc\x11\xc5\xc9?\xf8\x84\xf7C֔\x10\xce?\xf8\x95\xd6&\xe8\v.\xd1?\xf8:\x8c0\xe2\x8ey\xd5?\xf8\xdfA:\xdc\x11\xc5\xd9?\xf8\x84\xf7C֔\x10\xde?\xf8\x95\xd6&\xe8\v.\xe1?\xf8:\x8c0\xe2\x8ey\xe5?\xf8\xdfA:\xdc\x11\xc5\xe9?\xf8\x84\xf7C֔\x10\xee?\xf8\x95\xd6&\xe8\v.\xf1?\xf8:\x8c0\xe2\x8ey\xf5?\xf8\xdfA:\xdc\x11\xc5\xf9?\xf8\x84\xf7C֔\x10\xfe?\xf8\x95\xd6&\xe8\v.\x01@\xf8:\x8c0\xe2\x8ey\x05@\xf8\xdfA:\xdc\x11\xc5\t@\xf8\x84\xf7C֔\x10\x0e@\xf8\x95\xd6&\xe8\v.\x11@\xf8:\x8c0\xe2\x8ey\x15@\xf8\xdfA:\xdc\x11\xc5\x19@\xf8\x84\xf7C֔\x10\x1e@\xf8\x95\xd6&\xe8\v.!@\xf8:\x8c0\xe2\x8ey%@\xf8\xdfA:\xdc\x11\xc5)@\xf8\x84\xf7C֔\x10.@\xf8\x95\xd6&\xe8\v.1@\xf8:\x8c0\xe2\x8ey5@\xf8\xdfA:\xdc\x11\xc59@\xf8\x84\xf7C֔\x10>@\xf8\x95\xd6&\xe8\v.A@\xf8:\x8c0\xe2\x8eyE@\xf8\xdfA:\xdc\x11\xc5I@\xf8\x84\xf7C֔\x10N@\xf8\x95\xd6&\xe8\v.Q@\xf8:\x8c0\xe2\x8eyU@\xf8\xdfA:\xdc\x11\xc5Y@\xf8\x84\xf7C֔\x10^@\xf8\x95\xd6&\xe8\v.a@\xf8:\x8c0\xe2\x8eye@\xf8\xdfA:\xdc\x11\xc5i@\xf8\x84\xf7C֔\x10n@\xf8\x95\xd6&\xe8\v.q@\xf8:\x8c0\xe2\x8eyu@\xf8\xdfA:\xdc\x11\xc5y@\xf8\x84\xf7C֔\x10~@\xf8\x95\xd6&\xe8\v.\x81@\xf8:\x8c0\xe2\x8ey\x85@\xf8\xdfA:\xdc\x11ŉ@\xf8\x84\xf7C֔\x10\x8e@\xf8\x95\xd6&\xe8\v.\x91@\xf8:\x8c0\xe2\x8ey\x95@\xf8\xdfA:\xdc\x11ř@\xf8\x84\xf7C֔\x10\x9e@\xf8\x95\xd6&\xe8\v.\xa1@\xf8:\x8c0\xe2\x8ey\xa5@\xf8\xdfA:\xdc\x11ũ@\xf8\x84\xf7C֔\x10\xae@\xf8\x95\xd6&\xe8\v.\xb1@\xf8:\x8c0\xe2\x8ey\xb5@\xf8\xdfA:\xdc\x11Ź@\xf8\x84\xf7C֔\x10\xbe@\xf8\x95\xd6&\xe8\v.\xc1@\xf8:\x8c0\xe2\x8ey\xc5@\xf8\xdfA:\xdc\x11\xc5\xc9@\xf8\x84\xf7C֔\x10\xce@\xf8\x95\xd6&\xe8\v.\xd1@\xf8:\x8c0\xe2\x8ey\xd5@\xf8\xdfA:\xdc\x11\xc5\xd9@\xf8\x84\xf7C֔\x10\xde@\xf8\x95\xd6&\xe8\v.\xe1@\xf8:\x8c0\xe2\x8ey\xe5@\xf8\xdfA:\xdc\x11\xc5\xe9@\xf8\x84\xf7C֔\x10\xee@\xf8\x95\xd6&\xe8\v.\xf1@\xf8:\x8c0\xe2\x8ey\xf5@\xf8\xdfA:\xdc\x11\xc5\xf9@\xf8\x84\xf7C֔\x10\xfe@\xf8\x95\xd6&\xe8\v.\x01A\xfe\xf0\x7f\x00"

	dserialize := func(data string) *metrics.Float64Histogram {
		var res metrics.Float64Histogram
		if err := gob.NewDecoder(bytes.NewReader([]byte(data))).Decode(&res); err != nil {
			panic(err)
		}
		return &res
	}
	latency := RuntimeHistogramFromData(float64(time.Second), dserialize(sLatency))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		snap := latency.Snapshot()
		// These are the fields that influxdb accesses
		_ = snap.Count()
		_ = snap.Max()
		_ = snap.Mean()
		_ = snap.Min()
		_ = snap.StdDev()
		_ = snap.Variance()
		_ = snap.Percentiles([]float64{0.25, 0.5, 0.75, 0.95, 0.99, 0.999, 0.9999})
	}
}
