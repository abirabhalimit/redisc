package redisc

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/sivaosorg/wrapify"
)

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Getter settings
//_______________________________________________________________________

// IsEnabled returns true if the configuration is enabled, indicating that
// a connection to Redis should be attempted.
func (c *Settings) IsEnabled() bool {
	return c.enabled
}

// IsDebugging returns true if debugging is enabled in the configuration,
// which may allow more verbose logging.
func (c *Settings) IsDebugging() bool {
	return c.debugging
}

// PingInterval returns the interval at which the database connection is pinged.
// This value is used by the keepalive mechanism.
func (c *Settings) PingInterval() time.Duration {
	return c.pingInterval
}

// IsPingInterval returns true if keepalive is enabled and a ping interval is specified.
func (c *Settings) IsPingInterval() bool {
	return c.keepalive && c.pingInterval != 0
}

func (c *Settings) Conn() *connectionSettings {
	return c.conn
}

func (c *Settings) Retry() *retrySettings {
	return c.retry
}

func (c *Settings) Timeout() *timeoutSettings {
	return c.timeout
}

func (c *Settings) Pool() *poolSettings {
	return c.pool
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Getter Datasource
//_______________________________________________________________________

// Conn returns the underlying redis.Client connection instance in a thread-safe manner.
func (d *Datasource) Conn() *redis.Client {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.conn
}

// Wrap returns the current wrapify.R instance, which encapsulates the connection status,
// any error messages, and debugging information in a thread-safe manner.
func (d *Datasource) Wrap() wrapify.R {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.wrap
}

// Conf returns the Settings configuration associated with the Datasource.
func (d *Datasource) Conf() Settings {
	return d.conf
}

// IsConnected returns true if the current wrap indicates a successful connection to redis,
// otherwise it returns false.
func (d *Datasource) IsConnected() bool {
	return d.Wrap().IsSuccess()
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Setter settings
//_______________________________________________________________________

// SetEnable sets the enabled flag in the configuration and returns the updated RConf,
// allowing for method chaining.
func (c *Settings) SetEnable(value bool) *Settings {
	c.enabled = value
	return c
}

// SetDebug sets the debugging flag in the configuration and returns the updated RConf.
func (c *Settings) SetDebug(value bool) *Settings {
	c.debugging = value
	return c
}

// SetPingInterval sets the interval at which the connection is pinged for keepalive
// and returns the updated Settings.
func (c *Settings) SetPingInterval(value time.Duration) *Settings {
	c.pingInterval = value
	return c
}

// SetKeepalive enables or disables the automatic keepalive mechanism and returns the updated Settings.
func (c *Settings) SetKeepalive(value bool) *Settings {
	c.keepalive = value
	return c
}

func (c *Settings) SetConn(value *connectionSettings) *Settings {
	c.conn = value
	return c
}

func (c *Settings) SetRetry(value *retrySettings) *Settings {
	c.retry = value
	return c
}

func (c *Settings) SetTimeout(value *timeoutSettings) *Settings {
	c.timeout = value
	return c
}

func (c *Settings) SetPool(value *poolSettings) *Settings {
	c.pool = value
	return c
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Setter connectionSettings
//_______________________________________________________________________

func (c *connectionSettings) SetNetwork(value string) *connectionSettings {
	c.network = value
	return c
}

func (c *connectionSettings) SetConnectionStrings(value string) *connectionSettings {
	c.connectionStrings = value
	return c
}

func (c *connectionSettings) SetPassword(value string) *connectionSettings {
	c.password = value
	return c
}

func (c *connectionSettings) SetDatabase(value int) *connectionSettings {
	c.database = value
	return c
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Setter retrySettings
//_______________________________________________________________________

func (r *retrySettings) SetMaxRetries(value int) *retrySettings {
	r.maxRetries = value
	return r
}

func (r *retrySettings) SetMinRetryBackoff(value time.Duration) *retrySettings {
	r.minRetryBackoff = value
	return r
}

func (r *retrySettings) SetMaxRetryBackoff(value time.Duration) *retrySettings {
	r.maxRetryBackoff = value
	return r
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Setter timeoutSettings
//_______________________________________________________________________

func (t *timeoutSettings) SetConnTimeout(value time.Duration) *timeoutSettings {
	t.connTimeout = value
	return t
}

func (t *timeoutSettings) SetReadTimeout(value time.Duration) *timeoutSettings {
	t.readTimeout = value
	return t
}

func (t *timeoutSettings) SetWriteTimeout(value time.Duration) *timeoutSettings {
	t.writeTimeout = value
	return t
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Setter poolSettings
//_______________________________________________________________________

func (p *poolSettings) SetPoolSize(value int) *poolSettings {
	p.poolSize = value
	return p
}

func (p *poolSettings) SetMinIdleConn(value int) *poolSettings {
	p.minIdleConn = value
	return p
}

func (p *poolSettings) SetMaxConnAge(value time.Duration) *poolSettings {
	p.maxConnAge = value
	return p
}

func (p *poolSettings) SetPoolTimeout(value time.Duration) *poolSettings {
	p.poolTimeout = value
	return p
}

func (p *poolSettings) SetIdleTimeout(value time.Duration) *poolSettings {
	p.idleTimeout = value
	return p
}

func (p *poolSettings) SetIdleCheckFrequency(value time.Duration) *poolSettings {
	p.idleCheckFrequency = value
	return p
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Setter Datasource
//_______________________________________________________________________

// SetConn safely updates the internal redis.Client connection of the Datasource and returns
// the updated Datasource for method chaining.
func (d *Datasource) SetConn(value *redis.Client) *Datasource {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.conn = value
	return d
}

// SetWrap safely updates the wrapify.R instance (which holds connection status and error info)
// of the Datasource and returns the updated Datasource.
func (d *Datasource) SetWrap(value wrapify.R) *Datasource {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.wrap = value
	return d
}

// SetOn sets the callback function that is invoked upon connection state changes (e.g., during keepalive events)
// and returns the updated Datasource for method chaining.
func (d *Datasource) SetOn(fnc func(response wrapify.R)) *Datasource {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.on = fnc
	return d
}

// SetOnReplica sets the callback function that is invoked for events specific to replica connections,
// such as replica failovers, reconnection attempts, or health status updates.
// This function accepts a callback that receives both the current status (encapsulated in wrapify.R)
// and a pointer to the Datasource representing the replica connection (replicator), allowing external
// components to implement custom logic for replica management. The updated Datasource instance is returned
// to support method chaining.
func (d *Datasource) SetOnReplica(fnc func(response wrapify.R, replicator *Datasource)) *Datasource {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.onReplica = fnc
	return d
}

// SetNotifier sets the callback function that is invoked for significant datasource events,
// such as reconnection attempts, keepalive signals, or other diagnostic updates.
// This function stores the provided notifier, which can be used to asynchronously notify
// external components of changes in the connection's status, and returns the updated Datasource instance
// to support method chaining.
func (d *Datasource) SetNotifier(fnc func(response wrapify.R)) *Datasource {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.notifier = fnc
	return d
}
