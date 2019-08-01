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
	DataFile string           `json:"data,omitempty"`
}

//ConnectionConfig will contain URL & related parameters
type ConnectionConfig struct {
	URL string `json:"url"`
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
	Type     string          `json:"type"`
	Duration int             `json:"duration,omitempty"`
	SendJSON json.RawMessage `json:"send,omitempty"`
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
	Success           bool          `json:"success"`
	DNSResolutionTime time.Duration `json:"dnstime"`
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
	ErrorSet                map[string]int
}
