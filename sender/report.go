package sender

import (
	"log"
	"sort"
	"strings"
	"time"
)

const (
	barChar = "∎"
)

type report struct {
	avgTotal float64
	fastest  float64
	slowest  float64
	average  float64
	rps      float64

	results chan *result
	total   time.Duration

	errorDist      map[string]int
	statusCodeDist map[int]int
	lats           []float64
	sizeTotal      int64
}

func newReport(size int, results chan *result, total time.Duration) *report {
	return &report{
		results: results,
		total:   total,
	}
}

func (r *report) finalize() {
	for res := range r.results {
		if res.err != nil {
			r.errorDist[res.err.Error()]++
		} else {
			r.lats = append(r.lats, res.duration.Seconds())
			r.avgTotal += res.duration.Seconds()
		}
	}
	r.rps = float64(len(r.lats)) / r.total.Seconds()
	r.average = r.avgTotal / float64(len(r.lats))
	r.print()
}

func (r *report) print() {
	if len(r.lats) > 0 {
		sort.Float64s(r.lats)
		r.fastest = r.lats[0]
		r.slowest = r.lats[len(r.lats)-1]
		log.Printf("\nSummary:\n")
		log.Printf("  Total:\t%4.4f secs\n", r.total.Seconds())
		log.Printf("  Slowest:\t%4.4f secs\n", r.slowest)
		log.Printf("  Fastest:\t%4.4f secs\n", r.fastest)
		log.Printf("  Average:\t%4.4f secs\n", r.average)
		log.Printf("  Requests/sec:\t%4.4f\n", r.rps)

		r.printHistogram()
		r.printLatencies()
	}
}

func (r *report) printHistogram() {
	bc := 10
	buckets := make([]float64, bc+1)
	counts := make([]int, bc+1)
	bs := (r.slowest - r.fastest) / float64(bc)
	for i := 0; i < bc; i++ {
		buckets[i] = r.fastest + bs*float64(i)
	}
	buckets[bc] = r.slowest
	var bi int
	var max int
	for i := 0; i < len(r.lats); {
		if r.lats[i] <= buckets[bi] {
			i++
			counts[bi]++
			if max < counts[bi] {
				max = counts[bi]
			}
		} else if bi < len(buckets)-1 {
			bi++
		}
	}
	log.Printf("\nResponse time histogram:\n")
	for i := 0; i < len(buckets); i++ {
		// Normalize bar lengths.
		var barLen int
		if max > 0 {
			barLen = (counts[i]*40 + max/2) / max
		}
		log.Printf("  %4.3f [%v]\t|%v\n", buckets[i], counts[i], strings.Repeat(barChar, barLen))
	}
}

// printLatencies prints percentile latencies.
func (r *report) printLatencies() {
	pctls := []int{10, 25, 50, 75, 90, 95, 99}
	data := make([]float64, len(pctls))
	j := 0
	for i := 0; i < len(r.lats) && j < len(pctls); i++ {
		current := i * 100 / len(r.lats)
		if current >= pctls[j] {
			data[j] = r.lats[i]
			j++
		}
	}
	log.Printf("\nLatency distribution:\n")
	for i := 0; i < len(pctls); i++ {
		if data[i] > 0 {
			log.Printf("  %v%% in %4.4f secs\n", pctls[i], data[i])
		}
	}
}
