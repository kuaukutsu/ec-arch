package postgres

import "time"

type Option func(*Pgsql)

func MaxPoolSize(size int) Option {
	return func(c *Pgsql) {
		c.maxPoolSize = size
	}
}

func ConnAttempts(attempts int) Option {
	return func(c *Pgsql) {
		c.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(c *Pgsql) {
		c.connTimeout = timeout
	}
}
