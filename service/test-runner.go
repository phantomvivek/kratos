package service

import (
	"fmt"
	"math"
	"time"

	"github.com/kratos/config"
	"github.com/kratos/models"
)

//Runner Handles all the running of tests & running config
type Runner struct {
	TestDoneChan    chan bool
	SocketDoneChan  chan bool
	SocketDoneCount int
	TotalCount      int
	ErrChan         chan error
	HostURL         string
	ConnectTimeout  int
	MaxDataLength   int
	DataIndex       int
	Tests           []*models.Test
	HitRates        []models.HitRate
	Flows           []models.ConnectionBucket
}

//TestRunner singleton runner that will run tests
var TestRunner Runner

//Initialize initializes the runner with config and channels
func (r *Runner) Initialize() {
	TestRunner = Runner{
		ErrChan:         make(chan error),
		TestDoneChan:    make(chan bool, 1),
		SocketDoneChan:  make(chan bool, 1),
		TotalCount:      0,
		SocketDoneCount: 0,
		MaxDataLength:   0,
		DataIndex:       -1,
		HostURL:         config.Config.Config.URL,
		ConnectTimeout:  config.Config.Config.Timeout,
		HitRates:        config.Config.HitRates,
	}

	TestRunner.Tests = make([]*models.Test, 0)

	for _, test := range config.Config.Tests {
		testRef := test
		TestRunner.Tests = append(TestRunner.Tests, &testRef)
	}

	//Defaults to 10 seconds
	if TestRunner.ConnectTimeout == 0 {
		TestRunner.ConnectTimeout = 10
	}

	//Connect reporter module to connect to any third party reporting tool like a statsd daemon
	if config.Config.Reporter.Type == "statsd" {
		Reporter.Connect(&config.Config.Reporter)
	}
}

//Start starts the tests
func (r *Runner) Start() {

	//Prepare per second buckets to determine how many sockets are to be opened per second
	r.PrepareBuckets()

	//Prepare data that will be sent to sockets in case any test has a message & replace string
	handler := DataHandler{}

	r.MaxDataLength = handler.PrepareTestData(config.Config.DataFile, r.TotalCount, r.Tests)

	//Start the error listener
	go r.ErrorListener()

	//Start listening to test completions
	go r.CompleteNotify()

	//Start the reporter
	go Reporter.Start()

	//Run tests!
	r.RunTests()
}

//CompleteNotify is notified when a socket test completes
func (r *Runner) CompleteNotify() {
	for {
		select {
		case <-r.SocketDoneChan:
			r.SocketDoneCount++
			//fmt.Println("Test done for some socket", r.TotalCount, r.SocketDoneCount)
			if r.SocketDoneCount >= r.TotalCount {

				//Tell reporter that test is completed
				Reporter.TestDoneChan <- true

				//Notify that tests are completed & program can exit
				//r.TestDoneChan <- true
				break
			}
		}
	}
}

//ErrorListener prints out errors
func (r *Runner) ErrorListener() {

	for {
		select {
		case err := <-r.ErrChan:
			//Do something with err!
			fmt.Println("Error encountered", err)
		}
	}
}

//PrepareBuckets prepares buckets per second of number of socket connections to open
func (r *Runner) PrepareBuckets() {

	totalDuration := 0
	for _, rate := range r.HitRates {
		totalDuration += rate.Duration
	}

	r.Flows = make([]models.ConnectionBucket, totalDuration)
	counter := 0
	var currConn float64

	for idx, rate := range r.HitRates {

		//Used to count exactly how many connections will be made in this hitrate
		hitrateCount := 0

		//Start connections will be the current connection count
		rate.StartConnections = currConn

		incrPerSecond := (rate.EndConnections - currConn) / float64(rate.Duration)
		for i := 0; i < rate.Duration; i++ {

			currConn += incrPerSecond
			count := int(math.Round(currConn))

			//Add to the total count
			r.TotalCount += count
			hitrateCount += count

			flow := models.ConnectionBucket{
				Idx:         idx,
				Count:       count,
				IncrementBy: incrPerSecond,
			}

			r.Flows[counter] = flow
			counter++
		}

		//Save the number of connections this hitrate will have
		rate.Connections = hitrateCount
		Reporter.MakeHitRateStat(idx, rate)

		//fmt.Printf("Config for hitrate:\tstart=%v, end=%v, total=%v, duration=%vs\n", rate.StartConnections, rate.EndConnections, rate.Connections, rate.Duration)
	}
}

//RunTests runs the tests according to the flow
func (r *Runner) RunTests() {

	//We divide every 10 milliseconds for opening sockets. This can be made more granular
	for _, flow := range r.Flows {

		/*
			We calculate sockets to be opened per 10ms,
			we shave off the decimal from the whole number, like if perTenMs is 1.12, shave is 0.12 and perTenMs becomes 1.00
			We maintain shaveIncr, starting from 0.00, and keep adding shave to it. Once it reaches 1.00 or above, we open a socket
			and reset shaveIncr to 0.00
		*/
		perTenMs := float64(flow.Count) / float64(100)
		shave := math.Mod(perTenMs, 1.00)
		perTenMs = math.Floor(perTenMs - shave)
		shaveIncr := 0.00

		//Round this to two decimals
		shave = math.Round(shave/0.01) * 0.01

		//Every 10 ms we will have multiple sockets to be opened
		for count := 0; count < 100; count++ {

			shaveIncr += shave
			if shaveIncr > 1.00 {
				//fmt.Println("Opening socket in shave condition!")
				r.OpenSocket(flow.Idx)
				shaveIncr -= 1.00
			}

			for i := 0; i < int(perTenMs); i++ {

				//Start a socket!
				r.OpenSocket(flow.Idx)
				//fmt.Println("Opening socket!", i, perTenMs)
			}

			//Wait for 10ms
			localTimer := time.NewTimer(10 * time.Millisecond)
			<-localTimer.C
		}

		//Since shave incr is greater than 0.5, we need to open a socket. This value is mostly very close to 0.99
		if shaveIncr > 0.5 {
			r.OpenSocket(flow.Idx)
		}
	}
}

//OpenSocket opens a socket.. this was repeated code
func (r *Runner) OpenSocket(hitIdx int) {

	r.DataIndex++
	if r.DataIndex >= r.MaxDataLength {
		r.DataIndex = 0
	}

	//Open a socket
	go SocketRun(r.HostURL, r.ConnectTimeout, r.Tests, r.DataIndex, r.SocketDoneChan, r.ErrChan, hitIdx, Reporter.ReportChan)
}
