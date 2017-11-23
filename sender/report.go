package sender

import (
	"log"
	"sort"
	"time"
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
	}
}
