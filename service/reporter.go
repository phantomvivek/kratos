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
	RateStats     map[int]*models.HitRateStats
	ReportChan    chan models.SocketStats
	TestDoneChan  chan bool
	ReportString  string
	HitrateString string
	TabWriter     *tabwriter.Writer
}

//Reporter singleton object
var Reporter StatsReporter

func init() {
	Reporter = StatsReporter{
		RateStats:    make(map[int]*models.HitRateStats),
		ReportChan:   make(chan models.SocketStats),
		TestDoneChan: make(chan bool),
	}

	Reporter.ReportString = "Connections\t[total]\t%v sockets\n" +
		"Connect\t[success, error, timeout]\t%v, %v, %v\n" +
		"Connect Latency\t[min, p50, p95, p99, max]\t%s, %s, %s, %s, %s\n" +
		"DNS Latency\t[min, p50, p95, p99, max]\t%s, %s, %s, %s, %s\n"

	Reporter.HitrateString = "Hitrate Connection Parameters\tstart=%v, end=%v, total=%v, duration=%vs\n"

	Reporter.TabWriter = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', tabwriter.StripEscape)
}

//MakeHitRateStat makes a stat object based on the hitrate id
func (r *StatsReporter) MakeHitRateStat(idx int, hitrate models.HitRate) {

	hrStat := models.HitRateStats{
		HitRateRef:              &hitrate,
		TotalConnections:        0,
		ConnectSuccess:          0.0,
		ConnectFailure:          0.0,
		ConnectTimeout:          0.0,
		ConnectLatencies:        tdigest.NewWithCompression(100),
		DNSResolutionLatencies:  tdigest.NewWithCompression(100),
		ConnectLatencyMin:       0,
		ConnectLatencyMax:       0,
		DNSResolutionLatencyMin: 0,
		DNSResolutionLatencyMax: 0,
		ErrorSet:                make(map[string]int),
	}

	r.RateStats[idx] = &hrStat
}

//Start starts the reporter to listen to any metric data coming on channel
func (r *StatsReporter) Start() {

	for {
		select {
		case metric := <-r.ReportChan:

			//Use this hrStat for measurements
			hrStat := r.RateStats[metric.HitrateIndex]

			hrStat.TotalConnections++
			if metric.Success {
				hrStat.ConnectSuccess++
			} else {
				hrStat.ConnectFailure++
			}

			if hrStat.ConnectLatencyMin > float64(metric.ConnectTime) || hrStat.ConnectLatencyMin == 0 {
				hrStat.ConnectLatencyMin = float64(metric.ConnectTime)
			}

			if hrStat.ConnectLatencyMax < float64(metric.ConnectTime) {
				hrStat.ConnectLatencyMax = float64(metric.ConnectTime)
			}

			if hrStat.DNSResolutionLatencyMin > float64(metric.DNSResolutionTime) || hrStat.DNSResolutionLatencyMin == 0 {
				hrStat.DNSResolutionLatencyMin = float64(metric.DNSResolutionTime)
			}

			if hrStat.DNSResolutionLatencyMax < float64(metric.DNSResolutionTime) {
				hrStat.DNSResolutionLatencyMax = float64(metric.DNSResolutionTime)
			}

			hrStat.ConnectLatencies.Add(float64(metric.ConnectTime), 1)
			hrStat.DNSResolutionLatencies.Add(float64(metric.DNSResolutionTime), 1)

			if metric.ErrorString != "" {
				if _, ok := hrStat.ErrorSet[metric.ErrorString]; ok {
					hrStat.ErrorSet[metric.ErrorString]++
				} else {
					hrStat.ErrorSet[metric.ErrorString] = 1
				}
			}

			if hrStat.TotalConnections >= hrStat.HitRateRef.Connections {
				//Report this hit rate as all connections for this hitrate have finished
				r.LogHitrate(hrStat.HitRateRef)
				r.Report(hrStat)
			}

		case <-r.TestDoneChan:
			//Test done
			fmt.Fprintln(Reporter.TabWriter, "All Tests Complete\tFinal Results Below:")

			//Report all stats from all hitratestats

			//Flush the tabwriter
			Reporter.TabWriter.Flush()

			//Program can exit after above reporting
			TestRunner.TestDoneChan <- true
		}
	}
}

func (r *StatsReporter) durationStr(dur float64) time.Duration {
	return time.Duration(dur)
}

//Report prints out the stats as they currently stand
func (r *StatsReporter) Report(hrStat *models.HitRateStats) {

	//Reporting stats from the tests
	if _, err := fmt.Fprintf(Reporter.TabWriter, r.ReportString,
		hrStat.TotalConnections,
		hrStat.ConnectSuccess, hrStat.ConnectFailure, hrStat.ConnectTimeout,
		time.Duration(hrStat.ConnectLatencyMin), r.durationStr(hrStat.ConnectLatencies.Quantile(0.5)), r.durationStr(hrStat.ConnectLatencies.Quantile(0.95)), r.durationStr(hrStat.ConnectLatencies.Quantile(0.99)), time.Duration(hrStat.ConnectLatencyMax),
		time.Duration(hrStat.DNSResolutionLatencyMin), r.durationStr(hrStat.DNSResolutionLatencies.Quantile(0.5)), r.durationStr(hrStat.DNSResolutionLatencies.Quantile(0.95)), r.durationStr(hrStat.DNSResolutionLatencies.Quantile(0.99)), time.Duration(hrStat.DNSResolutionLatencyMax),
	); err != nil {
		fmt.Println("Reporting error", err)
	}

	if len(hrStat.ErrorSet) == 0 {
		if _, err := fmt.Fprintf(Reporter.TabWriter, "Error Set\t[error, count]\tNo Errors\n\n"); err != nil {
			fmt.Println("Reporting error", err)
		}
	} else {
		for errStr, count := range hrStat.ErrorSet {
			if _, err := fmt.Fprintf(Reporter.TabWriter, "Error Set\t[error, count]\t%s, %v\n\n", errStr, count); err != nil {
				fmt.Println("Reporting error", err)
			}
		}
	}

	//Flush the tabwriter
	Reporter.TabWriter.Flush()
}

//LogHitrate logs the current hitrate
func (r *StatsReporter) LogHitrate(hitrate *models.HitRate) {

	if _, err := fmt.Fprintf(Reporter.TabWriter, r.HitrateString, hitrate.StartConnections, hitrate.EndConnections, hitrate.Connections, hitrate.Duration); err != nil {
		fmt.Println("Reporting error", err)
		return
	}

	//Flush the tabwriter
	Reporter.TabWriter.Flush()
}
