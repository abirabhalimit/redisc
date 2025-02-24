package redisc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/sivaosorg/loggy"
	"github.com/sivaosorg/wrapify"
)

func NewClient(conf Settings) *Datasource {
	datasource := &Datasource{
		conf: conf,
	}
	start := time.Now()
	if !conf.IsEnabled() {
		datasource.SetWrap(wrapify.
			WrapServiceUnavailable("Redis service unavailable", nil).
			WithDebuggingKV("executed_in", time.Since(start).String()).
			WithHeader(wrapify.ServiceUnavailable).
			Reply())
		return datasource
	}
	ops := datasource.getOptions()
	c := redis.NewClient(ops)

	// Use a context with timeout to verify the connection via ping.
	err := c.Ping().Err()
	if err != nil {
		datasource.SetWrap(
			wrapify.
				WrapInternalServerError("The redis server is unreachable", nil).
				WithDebuggingKV("redis_conn_str", conf.String(true)).
				WithDebuggingKV("executed_in", time.Since(start).String()).
				WithErrSck(err).
				WithHeader(wrapify.InternalServerError).
				Reply(),
		)
		return datasource
	}

	// Set the established connection and update the wrap response to indicate success.
	datasource.SetConn(c)
	datasource.SetWrap(wrapify.New().
		WithStatusCode(http.StatusOK).
		WithDebuggingKV("redis_conn_str", conf.String(true)).
		WithDebuggingKV("executed_in", time.Since(start).String()).
		WithMessagef("Successfully connected to the redis server: '%s'", conf.String(true)).
		WithHeader(wrapify.OK).
		Reply())

	// If keepalive is enabled, initiate the background routine to monitor connection health.
	if conf.keepalive {
		datasource.keepalive()
	}
	return datasource
}

func (d *Datasource) AllKeys() wrapify.R {
	if !d.IsConnected() {
		return d.Wrap()
	}
	keys := make(map[string]string)
	var cursor uint64
	for {
		var batchKeys []string
		var err error
		batchKeys, cursor, err = d.Conn().Scan(cursor, "*", 10).Result()
		if err != nil {
			if d.conf.IsDebugging() {
				loggy.Errorf("A technical issue arose during the retrieval of all keys: %s", err.Error())
			}
			response := wrapify.
				WrapInternalServerError("A technical issue arose during the retrieval of all keys", nil).
				WithHeader(wrapify.InternalServerError).
				WithDebuggingKV("function", "all_keys").
				WithErrSck(err).Reply()
			d.notify(response)
			return response
		}
		for _, key := range batchKeys {
			keyType, err := d.Conn().Type(key).Result()
			if err != nil {
				if d.conf.IsDebugging() {
					loggy.Errorf("Failed to determine the type of key '%s': %s", key, err.Error())
				}
				response := wrapify.
					WrapInternalServerError("", nil).
					WithMessagef("Failed to determine the type of key '%s'", key).
					WithHeader(wrapify.InternalServerError).
					WithDebuggingKV("function", "all_keys").
					WithErrSck(err).Reply()
				d.notify(response)
				return response
			}
			keys[key] = keyType
		}
		if cursor == 0 {
			break
		}
	}
	return wrapify.WrapOk("Successfully retrieved all keys", keys).WithTotal(len(keys)).WithHeader(wrapify.OK).Reply()
}

func (d *Datasource) getOptions() *redis.Options {
	ops := &redis.Options{
		Network:            d.conf.conn.network,
		Addr:               d.conf.conn.connectionStrings,
		Password:           d.conf.conn.password,
		DB:                 d.conf.conn.database,
		MaxRetries:         d.conf.retry.maxRetries,
		MinRetryBackoff:    d.conf.retry.minRetryBackoff,
		MaxRetryBackoff:    d.conf.retry.maxRetryBackoff,
		DialTimeout:        d.conf.timeout.connTimeout,
		ReadTimeout:        d.conf.timeout.readTimeout,
		WriteTimeout:       d.conf.timeout.writeTimeout,
		PoolSize:           d.conf.pool.poolSize,
		MinIdleConns:       d.conf.pool.minIdleConn,
		MaxConnAge:         d.conf.pool.maxConnAge,
		PoolTimeout:        d.conf.pool.poolTimeout,
		IdleTimeout:        d.conf.pool.idleTimeout,
		IdleCheckFrequency: d.conf.pool.idleCheckFrequency,
	}
	return ops
}

// keepalive initiates a background goroutine that periodically pings the redis server
// to monitor connection health. Upon detecting a failure in the ping, it attempts to reconnect
// and subsequently invokes a callback (if set) with the updated connection status. This mechanism
// ensures that the Datasource remains current with respect to the connection state.
//
// The ping interval is determined by the configuration's PingInterval; if it is not properly set,
// a default interval is used.
func (d *Datasource) keepalive() {
	interval := d.conf.PingInterval()
	if interval <= 0 {
		interval = defaultPingInterval
	}
	var response wrapify.R
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		reconnectAttempt := 0 // Initialize reconnect attempt count
		for range ticker.C {
			ps := time.Now()
			if err := d.ping(); err != nil {
				duration := time.Since(ps)
				response = wrapify.WrapInternalServerError("The redis server is currently unreachable. Initiating reconnection process...", nil).
					WithDebuggingKV("redis_conn_str", d.conf.String(true)).
					WithDebuggingKV("ping_executed_in", duration.String()).
					WithDebuggingKV("ping_start_at", ps.Format(defaultTimeFormat)).
					WithDebuggingKV("ping_end_at", ps.Add(duration).Format(defaultTimeFormat)).
					WithErrSck(err).
					WithHeader(wrapify.InternalServerError).
					Reply()

				ps = time.Now()
				if err := d.reconnect(); err != nil {
					duration := time.Since(ps)
					reconnectAttempt++ // Increment reconnect count on failure reconnect
					response = wrapify.WrapInternalServerError("The redis server remains unreachable. The reconnection attempt has failed", nil).
						WithDebuggingKV("redis_conn_str", d.conf.String(true)).
						WithDebuggingKV("reconnect_executed_in", duration.String()).
						WithDebuggingKV("reconnect_start_at", ps.Format(defaultTimeFormat)).
						WithDebuggingKV("reconnect_end_at", ps.Add(duration).Format(defaultTimeFormat)).
						WithDebuggingKV("reconnect_attempt", reconnectAttempt).
						WithErrSck(err).
						WithHeader(wrapify.InternalServerError).
						Reply()
				} else {
					duration := time.Since(ps)
					reconnectAttempt = 0
					response = wrapify.New().
						WithStatusCode(http.StatusOK).
						WithDebuggingKV("redis_conn_str", d.conf.String(true)).
						WithDebuggingKV("reconnect_executed_in", duration.String()).
						WithDebuggingKV("reconnect_start_at", ps.Format(defaultTimeFormat)).
						WithDebuggingKV("reconnect_end_at", ps.Add(duration).Format(defaultTimeFormat)).
						WithMessagef("The connection to the redis server has been successfully re-established: '%s'", d.conf.String(true)).
						WithHeader(wrapify.OK).
						Reply()
				}
			} else {
				duration := time.Since(ps)
				reconnectAttempt = 0
				response = wrapify.New().
					WithStatusCode(http.StatusOK).
					WithDebuggingKV("redis_conn_str", d.conf.String(true)).
					WithDebuggingKV("ping_executed_in", time.Since(ps).String()).
					WithDebuggingKV("ping_start_at", ps.Format(defaultTimeFormat)).
					WithDebuggingKV("ping_end_at", ps.Add(duration).Format(defaultTimeFormat)).
					WithMessagef("The connection to the redis server has been successfully established: '%s'", d.conf.String(true)).
					WithHeader(wrapify.OK).
					Reply()
			}
			d.SetWrap(response)
			d.invoke(response)
			d.invokeReplica(response, d)
		}
	}()
}

// ping performs a health check on the current redis connection by issuing a ping
// It returns an error if the connection is nil or if the ping operation fails.
//
// Returns:
//   - nil if the connection is healthy;
//   - an error if the connection is nil or the ping fails.
func (d *Datasource) ping() error {
	d.mu.RLock()
	conn := d.conn
	d.mu.RUnlock()
	if conn == nil {
		return fmt.Errorf("the redis connection is currently unavailable")
	}
	return conn.Ping().Err()
}

// reconnect attempts to establish a new connection to the redis server using the current configuration.
// If the new connection is successfully verified via ping, it replaces the existing connection in the Datasource.
// In the event that a previous connection exists, it is closed to release associated resources.
//
// Returns:
//   - nil if reconnection is successful;
//   - an error if the reconnection fails at any stage.
func (d *Datasource) reconnect() error {
	ops := d.getOptions()
	current := redis.NewClient(ops)
	if err := current.Ping().Err(); err != nil {
		current.Close()
		return err
	}

	d.mu.Lock()
	previous := d.conn
	d.conn = current
	d.mu.Unlock()
	if previous != nil {
		previous.Close()
	}
	return nil
}

// invoke safely retrieves the registered callback function and, if one is set,
// invokes it asynchronously with the provided wrapify.R response. This ensures that
// external consumers are notified of connection status changes without blocking the
// calling goroutine.
func (d *Datasource) invoke(response wrapify.R) {
	d.mu.RLock()
	callback := d.on
	d.mu.RUnlock()
	if callback != nil {
		go callback(response)
	}
}

// invokeReplica safely retrieves the registered replica callback function and, if one is set,
// invokes it asynchronously with the provided wrapify.R response and a pointer to the replica Datasource.
// This ensures that external consumers are notified of replica-specific connection status changes,
// such as replica failovers, reconnection attempts, or health updates, without blocking the calling goroutine.
func (d *Datasource) invokeReplica(response wrapify.R, replicator *Datasource) {
	d.mu.RLock()
	callback := d.onReplica
	d.mu.RUnlock()
	if callback != nil {
		go callback(response, replicator)
	}
}

// notify safely retrieves the registered notifier callback function and, if one is set,
// invokes it asynchronously with the provided wrapify.R response. This method allows the Datasource
// to notify external components of significant events (e.g., reconnection, keepalive updates)
// without blocking the calling goroutine, ensuring that notification handling is performed concurrently.
func (d *Datasource) notify(response wrapify.R) {
	d.mu.RLock()
	callback := d.notifier
	d.mu.RUnlock()
	if callback != nil {
		go callback(response)
	}
}
