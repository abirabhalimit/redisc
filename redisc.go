package redisc

import "time"

func NewSettings() *Settings {
	s := &Settings{}
	s.
		SetRetry(newRetrySettings()).
		SetTimeout(newTimeoutSettings()).
		SetPool(newPoolSettings()).
		SetConn(newConnSettings())
	return s
}

func newConnSettings() *connectionSettings {
	c := &connectionSettings{
		network:  "tcp", // Use TCP for most connections. Use "unix" if you prefer a Unix domain socket.
		database: 0,     // Connects to the first logical database. Change if your application needs a different DB.
	}
	return c
}

func newRetrySettings() *retrySettings {
	r := &retrySettings{
		maxRetries:      3,                      // Three retry attempts for commands, balancing resilience and responsiveness.
		minRetryBackoff: 8 * time.Millisecond,   // Provides a short initial delay between retries.
		maxRetryBackoff: 512 * time.Millisecond, // Caps the maximum delay between retries to prevent long waits.
	}
	return r
}

func newTimeoutSettings() *timeoutSettings {
	t := &timeoutSettings{
		connTimeout:  5 * time.Second, // Allows a moderate wait time when establishing a connection.
		readTimeout:  3 * time.Second, // Sufficient for most environments to avoid long hangs during operations.
		writeTimeout: 3 * time.Second, // Sufficient for most environments to avoid long hangs during operations.
	}
	return t
}

func newPoolSettings() *poolSettings {
	p := &poolSettings{
		poolSize:           10,              // Supports moderate concurrency. Increase if your application has a high number of simultaneous requests.
		minIdleConn:        2,               // Keeps a couple of connections always ready to reduce latency.
		maxConnAge:         0,               // Connections are recycled indefinitely. Set a non-zero value to force periodic connection renewal.
		poolTimeout:        4 * time.Second, // Wait up to 4 seconds for a free connection from the pool.
		idleTimeout:        5 * time.Minute, // Closes idle connections after 5 minutes, freeing up resources.
		idleCheckFrequency: 1 * time.Minute, // Checks every minute to clear out idle connections.
	}
	return p
}
