package models

import (
	"encoding/json"
	"time"

	"github.com/influxdata/tdigest"
)

//Configuration for incoming config
type Configuration struct {
	Config   ConnectionConfig `json:"config"`
	HitRates []HitRate        `json:"hitrate"`
	Tests    []Test           `json:"tests"`
	DataFile string           `json:"dataFile,omitempty"`
	Reporter ReporterConfig   `json:"reporter"`
}

//ReporterConfig to read the reporting config
type ReporterConfig struct {
	Type   string `json:"type"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Prefix string `json:"prefix"`
}

//ConnectionConfig will contain URL & related parameters
type ConnectionConfig struct {
	URL     string `json:"url"`
	Timeout int    `json:"timeout,omitempty"`
}

//HitRate defines a rate of hitting the test applicant with connections
type HitRate struct {
	StartConnections float64 `json:"start"`
	EndConnections   float64 `json:"end"`
	Duration         int     `json:"duration"`
	Connections      int     `json:"connCount"`
}

//Test type is used for sending messages
type Test struct {
	Type       string          `json:"type"`
	Duration   int             `json:"duration,omitempty"`
	SendJSON   json.RawMessage `json:"send,omitempty"`
	ReplaceStr bool            `json:"replace,omitempty"`
	Data       *TestData       `json:"testdata,omitempty"`
}

//TestData will hold all the constructed messages after replacing variables from file
type TestData struct {
	Counter   int               `json:"counter"`
	DataArray []json.RawMessage `json:"messageArray"`
}

//TestDataConfig saves a config for replacing a particular column from csv in the message
type TestDataConfig struct {
	ColumnIdx int
	TextBytes []byte
}

//ConnectionBucket this has a per second count and incremented by previous second
type ConnectionBucket struct {
	Idx         int     `json:"index"`
	Count       int     `json:"count"`
	IncrementBy float64 `json:"incrBy"`
}

//SocketStats used to measure timing stats
type SocketStats struct {
	HitrateIndex      int           `json:"hrIdx"`
	ConnectTime       time.Duration `json:"connecttime"`
	DNSResolutionTime time.Duration `json:"dnstime"`
	OverallTime       time.Duration `json:"overalltime"`
	Success           bool          `json:"success"`
	ErrorString       string        `json:"error"`
}

//HitRateStats will store all stats related to this particular hit rate
type HitRateStats struct {
	HitRateRef              *HitRate
	TotalConnections        int
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
	OverallLatencies        *tdigest.TDigest
	OverallLatencyMin       float64
	OverallLatencyMax       float64
	ErrorSet                map[string]int
}
