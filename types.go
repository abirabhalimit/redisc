package redisc

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/sivaosorg/wrapify"
)

type Settings struct {
	enabled   bool
	debugging bool

	// Defines the frequency at which the connection is pinged.
	// This interval is used by the keepalive mechanism to periodically check the health of the
	// database connection. If a ping fails, a reconnection attempt may be triggered.
	pingInterval time.Duration

	// Indicates whether automatic keepalive is enabled for the Redis connection.
	// When set to true, a background process will periodically ping the database and attempt
	// to reconnect if the connection is lost.
	keepalive bool

	conn *connectionSettings

	retry *retrySettings

	timeout *timeoutSettings

	pool *poolSettings
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

	// The username for authentication (Redis 6+ supports usernames).
	// When your Redis server requires username-based authentication.
	username string

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

type Datasource struct {
	// A read-write mutex that ensures safe concurrent access to the Datasource fields.
	mu sync.RWMutex
	// An instance of Settings containing all the configuration parameters for the Redis connection.
	conf Settings
	// A wrapify.R instance that holds the current connection status, error messages, and debugging information.
	wrap wrapify.R
	// A pointer to an redis.Client object representing the active connection to the Redis database.
	conn *redis.Client
	// A callback function that is invoked asynchronously when there is a change in connection status,
	//  such as when the connection is lost, re-established, or its health is updated.
	on func(response wrapify.R)
	// onReplica is a callback function that is invoked asynchronously to handle events related to replica connections.
	// When the status of a replica datasource changes (e.g., during failover, reconnection, or health updates),
	// this function is triggered with the current status (encapsulated in wrapify.R) and a pointer to the Datasource
	// representing the replica connection. This allows external components to implement replica-specific logic
	// for tasks such as load balancing, monitoring, or failover handling independently of the primary connection.
	onReplica func(response wrapify.R, replicator *Datasource)
	// notifier is an optional callback function used to propagate notifications for significant datasource events.
	// It is invoked with the current status (encapsulated in wrapify.R) whenever notable events occur,
	// such as reconnection attempts, keepalive signals, or other diagnostic updates.
	// This allows external components to receive and handle these notifications independently of the primary connection status callback.
	notifier func(response wrapify.R)
}
