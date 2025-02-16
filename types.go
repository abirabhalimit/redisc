package redisc

import "time"

type settings struct {
	enabled   bool
	debugging bool
	conn      *connectionSettings
	retry     *retrySettings
	timeout   *timeoutSettings
	pool      *poolSettings
}

type connectionSettings struct {
	// Specifies the network protocol ("tcp" or "unix").
	// When you need a standard TCP connection or a Unix domain socket.
	// Default network is "tcp". Use "unix" for Unix domain sockets.
	network string

	// The address of your Redis server in "host:port" format.
	// Always required to locate your Redis server.
	// Default address for a local Redis instance.
	connectionStrings string

	// The password for authentication with the Redis server.
	// Required when your Redis instance is protected by a password.
	// Default is no password (set if your Redis requires authentication).
	password string

	// The Redis logical database number to select (default is 0).
	// Change this if your application uses a specific logical database.
	// Default DB index is 0.
	database int
}

type retrySettings struct {
	// Maximum number of retry attempts for a command if errors occur.
	// Increase in environments with intermittent connectivity issues.
	// When you want to automatically retry failed commands (e.g., due to temporary network issues).
	// Try up to n times before failing a command.
	maxRetries int

	// The minimum wait time between retries.
	// Adjust to control the pace of retry attempts.
	minRetryBackoff time.Duration

	// The maximum wait time between retry attempts.
	// Prevents the backoff duration from growing excessively.
	maxRetryBackoff time.Duration
}

type timeoutSettings struct {
	// Maximum duration to wait for a new connection to be established.
	// Helps ensure your application fails fast if the Redis server is unreachable.
	// Maximum time to wait for a connection to be established.
	connTimeout time.Duration

	// Maximum duration to wait for a command response from the Redis server.
	// Set according to your network conditions to avoid long blocking calls.
	// Maximum time to wait for a response (per command).
	readTimeout time.Duration

	// Maximum duration to wait when sending a command to the server.
	// Ensures that write operations do not hang indefinitely.
	// Maximum time to wait for sending a command.
	writeTimeout time.Duration
}

type poolSettings struct {
	// Maximum number of connections maintained in the pool.
	// Increase for high-concurrency applications.
	poolSize int

	// Minimum number of idle connections to keep open.
	// Helps reduce latency by always having ready-to-use connections.
	// Keep at least n idle connections ready.
	minIdleConn int

	// Maximum lifetime of a connection before it is closed.
	// Set to recycle connections periodically (0 means no limit).
	// Default is 0 means connections are reused indefinitely.
	maxConnAge time.Duration

	// Maximum duration to wait for a connection if the pool is exhausted.
	// Adjust to control blocking behavior under heavy load.
	// Maximum time to wait for a connection if all are busy.
	poolTimeout time.Duration

	// Duration after which an idle connection is closed.
	// Frees up resources by closing unused connections.
	idleTimeout time.Duration

	// How often the pool checks for idle connections to close.
	// Balance between timely cleanup and the overhead of checks.
	idleCheckFrequency time.Duration
}
