package models

import (
	"encoding/json"
	"time"
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
	EndConnections float64 `json:"end"`
	Duration       int     `json:"duration"`
}

//Test type is used for sending messages
type Test struct {
	Type     string          `json:"type"`
	Duration int             `json:"duration,omitempty"`
	SendJSON json.RawMessage `json:"send,omitempty"`
}

//ConnectionBucket this has a per second count and incremented by previous second
type ConnectionBucket struct {
	Count       int     `json:"count"`
	IncrementBy float64 `json:"incrBy"`
}

//SocketStats used to measure timing stats
type SocketStats struct {
	ConnectTime       time.Duration `json:"connecttime"`
	Success           bool          `json:"success"`
	DNSResolutionTime time.Duration `json:"dnstime"`
}
