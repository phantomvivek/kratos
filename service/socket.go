package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http/httptrace"
	"time"

	"github.com/gorilla/websocket"
	"github.com/phantomvivek/kratos/models"
)

//Socket = single socket connection to the host
type Socket struct {
	Connection  *websocket.Conn
	Dialer      *websocket.Dialer
	SocketStats *models.SocketStats
	Context     context.Context
}

//ContextKey used for getting ref out of context
type ContextKey string

//CustomDialer ..need to figure how to pass ID & channel for this
func CustomDialer(ctx context.Context, network, addr string) (net.Conn, error) {

	//These variables will always be set
	statsRef, ok := ctx.Value(ContextKey("StatsRef")).(*models.SocketStats)
	timeout, _ := ctx.Value(ContextKey("Timeout")).(int)

	var connectStart time.Time
	var dnsStart time.Time
	var overallTime time.Time
	var connectDiff time.Duration
	var dnsDiff time.Duration
	var overallDiff time.Duration

	overallTime = time.Now()

	ctTrace := &httptrace.ClientTrace{
		// GotConn: func(connInfo httptrace.GotConnInfo) {
		// 	fmt.Println("GOT connection")
		// },
		// GetConn: func(hostPort string) {
		// 	fmt.Println("GET connection")
		// },
		DNSStart: func(info httptrace.DNSStartInfo) {
			//fmt.Println("DNS Start")
			dnsStart = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			//fmt.Println("DNS Done")
			dnsDiff = time.Since(dnsStart)
		},
		// GotFirstResponseByte: func() {
		// 	fmt.Println("First byte!")
		// },
		ConnectStart: func(network, addr string) {
			//fmt.Println("Connect start")
			connectStart = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			//fmt.Println("Connect Done")
			connectDiff = time.Since(connectStart)
		},
	}
	traceCtx := httptrace.WithClientTrace(ctx, ctTrace)
	dialer := net.Dialer{
		Timeout: time.Duration(timeout) * time.Second,
	}

	conn, err := dialer.DialContext(traceCtx, network, addr)

	overallDiff = time.Since(overallTime)

	if ok {
		statsRef.ConnectTime = connectDiff
		statsRef.DNSResolutionTime = dnsDiff
		statsRef.OverallTime = overallDiff
		if err != nil {
			statsRef.Success = false
			statsRef.ErrorString = err.Error()
		} else {
			statsRef.Success = true
		}
	}

	if err != nil {
		return conn, err
	}

	return conn, nil
}

//SocketRun goroutine that makes a socket collection with the host and starts the tests
func SocketRun(hostURL string, timeout int, tests []*models.Test, dataIdx int, doneChan chan bool, errChan chan error, hitIdx int, reporterChan chan *models.SocketStats) {

	socket := Socket{
		Dialer: &websocket.Dialer{
			HandshakeTimeout: 10 * time.Second,
			NetDialContext:   CustomDialer,
		},
		SocketStats: &models.SocketStats{
			HitrateIndex: hitIdx,
		},
	}

	socket.Context = context.WithValue(context.Background(), ContextKey("StatsRef"), socket.SocketStats)
	socket.Context = context.WithValue(socket.Context, ContextKey("Timeout"), timeout)

	err := socket.Connect(hostURL)
	if err != nil {
		errChan <- err
		doneChan <- true

		//In case of errors, we need to send the stats
		reporterChan <- socket.SocketStats

		return
	}

	socket.DoTests(tests, dataIdx)

	reporterChan <- socket.SocketStats

	//Tests would be complete
	doneChan <- true
}

//Connect connect the ws to host
func (s *Socket) Connect(url string) error {

	conn, _, err := s.Dialer.DialContext(s.Context, url, nil)
	if err != nil {
		//return err to the error channel
		//fmt.Println("Error in connection", err)
		return err
	}

	s.Connection = conn
	return nil
}

//DoTests runs through tests for this socket
func (s *Socket) DoTests(tests []*models.Test, dataIdx int) {

	for _, test := range tests {

		if test.Type == "message" {

			var msg json.RawMessage
			if test.ReplaceStr {
				msg = test.Data.DataArray[dataIdx]
			} else {
				msg = test.SendJSON
			}
			//Need to send message to the host
			err := s.Connection.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				//Log error
				fmt.Println("Error occured in sending message to host", err)
			}
		} else if test.Type == "sleep" {

			//Sleep for so many seconds
			localTimer := time.NewTimer(time.Duration(test.Duration) * time.Second)
			<-localTimer.C
		} else if test.Type == "disconnect" {

			//Need to disconnect the socket
			err := s.Connection.Close()
			if err != nil {
				fmt.Println("Error in closing")
			}

		} else {
			fmt.Println("Invalid type found", test.Type)
		}

		delay := time.NewTimer(10 * time.Millisecond)
		<-delay.C
	}
}
