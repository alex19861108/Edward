package main

import (
	"flag"
	"fmt"
	"git-pd.megvii-inc.com/liuwei02/Edward/reader"
	"git-pd.megvii-inc.com/liuwei02/Edward/sender"
	"git-pd.megvii-inc.com/liuwei02/Edward/taurusrpc"
	"os"
	"os/signal"
	"strings"
)

var (
	Version string
	Build   string
)

const (
	tcp_port = "8000"
)

var (
	h       = flag.Bool("h", false, "show this help")
	t       = flag.String("t", "dumps", "dumps or loads")
	i       = flag.String("i", "input", "input file")
	n       = flag.Int("n", 1000, "request count")
	qps     = flag.Float64("qps", 10, "qps")
	threads = flag.Int("threads", 1, "threads")
	address = flag.String("address", "localhost:"+tcp_port, "tcp client ")
)

var usage = `Version: ` + Version + `
Build: ` + Build + `

Usage: Edward [-h] [-t type] [-p port] [-i input] [-n request_count] [-qps qps]

Options:
`

func main() {

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *h {
		flag.Usage()
		os.Exit(0)
	}

	if strings.ToLower(*t) == "dumps" {
		taurusrpc.InitRPCServer(tcp_port)
	}

	if strings.ToLower(*t) == "loads" {
		contents := reader.TextReader(*i)
		//infos := transformer.Transform(contents)
		e := &sender.Edward{
			Contents:     contents,
			Address:      *address,
			Qps:          *qps,
			RequestCount: *n,
			ThreadCount:  *threads,
		}

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		go func() {
			<-sig
			e.Finish()
			os.Exit(1)
		}()

		e.Run()
	}
}
