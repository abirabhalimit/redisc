package redisc

import "time"

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Getter settings
//_______________________________________________________________________

// IsEnabled returns true if the configuration is enabled, indicating that
// a connection to Redis should be attempted.
func (c *settings) IsEnabled() bool {
	return c.enabled
}

// IsDebugging returns true if debugging is enabled in the configuration,
// which may allow more verbose logging.
func (c *settings) IsDebugging() bool {
	return c.debugging
}

func (c *settings) Conn() *connectionSettings {
	return c.conn
}

func (c *settings) Retry() *retrySettings {
	return c.retry
}

func (c *settings) Timeout() *timeoutSettings {
	return c.timeout
}

func (c *settings) Pool() *poolSettings {
	return c.pool
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Setter settings
//_______________________________________________________________________

// SetEnable sets the enabled flag in the configuration and returns the updated RConf,
// allowing for method chaining.
func (c *settings) SetEnable(value bool) *settings {
	c.enabled = value
	return c
}

// SetDebug sets the debugging flag in the configuration and returns the updated RConf.
func (c *settings) SetDebug(value bool) *settings {
	c.debugging = value
	return c
}

func (c *settings) SetConn(value *connectionSettings) *settings {
	c.conn = value
	return c
}

func (c *settings) SetRetry(value *retrySettings) *settings {
	c.retry = value
	return c
}

func (c *settings) SetTimeout(value *timeoutSettings) *settings {
	c.timeout = value
	return c
}

func (c *settings) SetPool(value *poolSettings) *settings {
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
