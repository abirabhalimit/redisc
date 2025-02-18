package redisc

import "time"

const (
	// defaultPingInterval defines the frequency at which the connection is pinged.
	defaultPingInterval = 30 * time.Second
	defaultTimeFormat   = "2006-01-02 15:04:05.000000"
)
