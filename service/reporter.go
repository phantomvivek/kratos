package service

import (
	"fmt"
	"time"

	"github.com/influxdata/tdigest"

	"github.com/kratos/models"
)

//StatsReporter is the struct that holds all test stats
type StatsReporter struct {
	TotalConnections float64
	TotalDuration    time.Duration
	ConnectSuccess   float64
	ConnectFailure   float64
	ConnectTimeout   float64
	Latencies        map[string]*tdigest.TDigest
	ReportChan       chan models.SocketStats
}

//Reporter singleton object
var Reporter StatsReporter

func init() {
	Reporter = StatsReporter{
		TotalConnections: 0.0,
		ConnectSuccess:   0.0,
		ConnectFailure:   0.0,
		ConnectTimeout:   0.0,
		Latencies:        make(map[string]*tdigest.TDigest),
		ReportChan:       make(chan models.SocketStats),
	}

	Reporter.Latencies["p50"] = tdigest.NewWithCompression(100)
	Reporter.Latencies["p95"] = tdigest.NewWithCompression(100)
	Reporter.Latencies["p99"] = tdigest.NewWithCompression(100)
}

//Start starts the reporter to listen to any metric data coming on channel
func (r *StatsReporter) Start() {
	for {

		select {

		case metric := <-r.ReportChan:
			fmt.Println("Stats are", metric.ConnectTime, metric.Success, metric.DNSResolutionTime)
		}
	}
}
