package service

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/influxdata/tdigest"

	"github.com/kratos/models"
)

//StatsReporter is the struct that holds all test stats
type StatsReporter struct {
	TotalConnections        float64
	TotalDuration           time.Duration
	ConnectSuccess          float64
	ConnectFailure          float64
	ConnectTimeout          float64
	ConnectLatencies        *tdigest.TDigest
	ConnectLatencyMin       float64
	ConnectLatencyMax       float64
	DNSResolutionLatencies  *tdigest.TDigest
	DNSResolutionLatencyMin float64
	DNSResolutionLatencyMax float64
	ReportChan              chan models.SocketStats
	TestDoneChan            chan bool
	ReportString            string
	ErrorSet                map[string]int
}

//Reporter singleton object
var Reporter StatsReporter

func init() {
	Reporter = StatsReporter{
		TotalConnections:        0.0,
		ConnectSuccess:          0.0,
		ConnectFailure:          0.0,
		ConnectTimeout:          0.0,
		ConnectLatencies:        tdigest.NewWithCompression(100),
		DNSResolutionLatencies:  tdigest.NewWithCompression(100),
		ReportChan:              make(chan models.SocketStats),
		TestDoneChan:            make(chan bool),
		ConnectLatencyMin:       0,
		ConnectLatencyMax:       0,
		DNSResolutionLatencyMin: 0,
		DNSResolutionLatencyMax: 0,
		ErrorSet:                make(map[string]int),
	}

	Reporter.ReportString = "Connections\t[total]\t%v sockets\n" +
		"Connect\t[success, error, timeout]\t%v, %v, %v\n" +
		"Connect Latency\t[min, p50, p95, p99, max]\t%s, %s, %s, %s, %s\n" +
		"DNS Latency\t[min, p50, p95, p99, max]\t%s, %s, %s, %s, %s\n"
}

//Start starts the reporter to listen to any metric data coming on channel
func (r *StatsReporter) Start() {

	for {
		select {
		case metric := <-r.ReportChan:
			r.TotalConnections++
			if metric.Success {
				r.ConnectSuccess++
			} else {
				r.ConnectFailure++
			}

			if r.ConnectLatencyMin > float64(metric.ConnectTime) || r.ConnectLatencyMin == 0 {
				r.ConnectLatencyMin = float64(metric.ConnectTime)
			}

			if r.ConnectLatencyMax < float64(metric.ConnectTime) {
				r.ConnectLatencyMax = float64(metric.ConnectTime)
			}

			if r.DNSResolutionLatencyMin > float64(metric.DNSResolutionTime) || r.DNSResolutionLatencyMin == 0 {
				r.DNSResolutionLatencyMin = float64(metric.DNSResolutionTime)
			}

			if r.DNSResolutionLatencyMax < float64(metric.DNSResolutionTime) {
				r.DNSResolutionLatencyMax = float64(metric.DNSResolutionTime)
			}

			r.ConnectLatencies.Add(float64(metric.ConnectTime), 1)
			r.DNSResolutionLatencies.Add(float64(metric.DNSResolutionTime), 1)

			if metric.ErrorString != "" {
				if _, ok := r.ErrorSet[metric.ErrorString]; ok {
					r.ErrorSet[metric.ErrorString]++
				} else {
					r.ErrorSet[metric.ErrorString] = 1
				}
			}

		case <-r.TestDoneChan:
			//Test done
			fmt.Println("All done!")

			r.Report()

			//Program can exit after above reporting
			TestRunner.TestDoneChan <- true
		}
	}
}

func (r *StatsReporter) durationStr(dur float64) time.Duration {
	return time.Duration(dur)
}

//Report prints out the stats as they currently stand
func (r *StatsReporter) Report() {

	//Reporting stats from the tests
	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', tabwriter.StripEscape)
	if _, err := fmt.Fprintf(tw, r.ReportString,
		r.TotalConnections,
		r.ConnectSuccess, r.ConnectFailure, r.ConnectTimeout,
		time.Duration(r.ConnectLatencyMin), r.durationStr(r.ConnectLatencies.Quantile(0.5)), r.durationStr(r.ConnectLatencies.Quantile(0.95)), r.durationStr(r.ConnectLatencies.Quantile(0.99)), time.Duration(r.ConnectLatencyMax),
		time.Duration(r.DNSResolutionLatencyMin), r.durationStr(r.DNSResolutionLatencies.Quantile(0.5)), r.durationStr(r.DNSResolutionLatencies.Quantile(0.95)), r.durationStr(r.DNSResolutionLatencies.Quantile(0.99)), time.Duration(r.DNSResolutionLatencyMax),
	); err != nil {
		fmt.Println("Reporting error", err)
	}

	if len(r.ErrorSet) == 0 {
		if _, err := fmt.Fprintf(tw, "Error Set\t[error, count]\n"); err != nil {
			fmt.Println("Reporting error", err)
		}
	} else {
		for errStr, count := range r.ErrorSet {
			if _, err := fmt.Fprintf(tw, "Error Set\t[error, count]\t%s, %v\n", errStr, count); err != nil {
				fmt.Println("Reporting error", err)
			}
		}
	}

	tw.Flush()
}
