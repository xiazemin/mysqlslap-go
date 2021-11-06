package report

import (
	"fmt"
	"sort"
	"time"
)

// LatencyDistribution holds latency distribution data
type LatencyDistribution struct {
	Percentage int           `json:"percentage"`
	Latency    time.Duration `json:"latency"`
}

func latencies(latencies []float64) []LatencyDistribution {
	pctls := []int{10, 25, 50, 75, 90, 95, 99}
	data := make([]float64, len(pctls))
	lt := float64(len(latencies))
	for i, p := range pctls {
		ip := (float64(p) / 100.0) * lt
		di := int(ip)

		// since we're dealing with 0th based ranks we need to
		// check if ordinal is a whole number that lands on the percentile
		// if so adjust accordingly
		if ip == float64(di) {
			di = di - 1
		}

		if di < 0 {
			di = 0
		}

		data[i] = latencies[di]
	}

	res := make([]LatencyDistribution, len(pctls))
	for i := 0; i < len(pctls); i++ {
		if data[i] > 0 {
			lat := time.Duration(data[i] * float64(time.Second))
			res[i] = LatencyDistribution{Percentage: pctls[i], Latency: lat}
		}
	}
	return res
}

// Bucket holds histogram data
type Bucket struct {
	// The Mark for histogram bucket in seconds
	Mark float64 `json:"mark"`

	// The count in the bucket
	Count int `json:"count"`

	// The frequency of results in the bucket as a decimal percentage
	Frequency float64 `json:"frequency"`
}

func histogram(latencies []float64, slowest, fastest float64) []Bucket {
	bc := 10
	buckets := make([]float64, bc+1)
	counts := make([]int, bc+1)
	bs := (slowest - fastest) / float64(bc)
	for i := 0; i < bc; i++ {
		buckets[i] = fastest + bs*float64(i)
	}
	buckets[bc] = slowest
	var bi int
	var max int
	for i := 0; i < len(latencies); {
		if latencies[i] <= buckets[bi] {
			i++
			counts[bi]++
			if max < counts[bi] {
				max = counts[bi]
			}
		} else if bi < len(buckets)-1 {
			//if there is no latencies in next bucket,the next counts should not increase;
			//so bi should skip the empty bucket
			for bi < len(buckets)-1 && latencies[i] > buckets[bi] {
				bi++
			}
		}
	}
	res := make([]Bucket, len(buckets))
	for i := 0; i < len(buckets); i++ {
		res[i] = Bucket{
			Mark:      buckets[i],
			Count:     counts[i],
			Frequency: float64(counts[i]) / float64(len(latencies)),
		}
	}
	return res
}

func GetReport(okLats []float64) ([]Bucket, []LatencyDistribution) {
	sort.Float64s(okLats)
	fastestNum := okLats[0]
	slowestNum := okLats[len(okLats)-1]
	return histogram(okLats, slowestNum, fastestNum), latencies(okLats)
}

func AverageFloat64(okLats []float64) float64 {
	if len(okLats) < 1 {
		return 0
	}
	var sum float64
	for _, num := range okLats {
		sum += num
	}
	return sum / float64(len(okLats))
}

type Report struct {
	Name string `json:"name,omitempty"`
	// EndReason StopReason `json:"endReason,omitempty"`
	Date time.Time `json:"date"`
	// Options   Options    `json:"options,omitempty"`

	Count   uint64        `json:"count"`
	Total   time.Duration `json:"total"`
	Average time.Duration `json:"average"`
	Fastest time.Duration `json:"fastest"`
	Slowest time.Duration `json:"slowest"`
	Rps     float64       `json:"rps"`

	ErrorDist      map[string]int `json:"errorDistribution"`
	StatusCodeDist map[string]int `json:"statusCodeDistribution"`

	LatencyDistribution []LatencyDistribution `json:"latencyDistribution"`
	Histogram           []Bucket              `json:"histogram"`
	// Details             []ResultDetail        `json:"details"`

	Tags map[string]string `json:"tags,omitempty"`
}

func (rp *Report) Print(s string) error {
	_, err := fmt.Print(s)
	return err
}
