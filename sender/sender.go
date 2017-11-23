package sender

import (
	"encoding/json"
	"git-pd.megvii-inc.com/liuwei02/Edward/taurusrpc"
	"log"
	"sync"
	"time"
)

type Record struct {
	// file context, one row
	c []byte
}

type result struct {
	err           error
	statusCode    int
	duration      time.Duration
	reqDuration   time.Duration // request "write" duration
	resDuration   time.Duration // response "read" duration
	contentLength int64
	response      string
}

var DEFAULT_COUNT = 1000

type Edward struct {
	Contents     []string
	Qps          float64
	start        time.Time
	stopCh       chan struct{}
	results      chan *result
	RequestCount int
	Address      string
	ThreadCount  int
}

func (e *Edward) Run() {
	log.SetFlags(0)
	e.results = make(chan *result, int(e.RequestCount))

	e.stopCh = make(chan struct{}, DEFAULT_COUNT)
	e.start = time.Now()

	e.runWorkers()
	e.Finish()
}

func (e *Edward) runWorkers() {
	var wg sync.WaitGroup
	wg.Add(e.ThreadCount)

	// Ignore the case where e.RequestCount % e.ThreadCount != 0.
	for i := 0; i < e.ThreadCount; i++ {
		go func() {
			e.runWorker(e.RequestCount/e.ThreadCount, float64(e.Qps)/float64(e.ThreadCount))
			wg.Done()
		}()
	}
	wg.Wait()
}

func (e *Edward) runWorker(n int, qps float64) {
	var throttle <-chan time.Time
	throttle = time.Tick(time.Duration(1e6/(qps)) * time.Microsecond)

	for i := 0; i < n; i++ {
		<-throttle
		content := e.getRequestBody(i)
		value, ok := content.(string)
		if !ok {
			continue
		}
		var info taurusrpc.SearchXIDInfo
		err := json.Unmarshal([]byte(value), &info)
		if err != nil {
			log.Fatal(err)
			continue
		}
		resStart := time.Now()
		resp := taurusrpc.CSearchXID(e.Address, &info)
		t := time.Now()
		duration := t.Sub(resStart)
		response, _ := json.Marshal(resp)
		e.results <- &result{
			duration: duration,
			response: string(response),
		}
	}
}

func (e *Edward) getRequestBody(idx int) interface{} {
	length := len(e.Contents)
	return e.Contents[idx%length]

}

func (e *Edward) Finish() {
	e.stopCh <- struct{}{}
	close(e.results)
	newReport(e.RequestCount, e.results, time.Now().Sub(e.start)).finalize()
}
