package netserver

import "time"

type Option func(*server)

func Address(address string) Option {
	return func(s *server) {
		s.address = address
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.readTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.writeTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.shutdownTimeout = timeout
	}
}

func IdleTimeout(timeout time.Duration) Option {
	return func(s *server) {
		s.idleTimeout = timeout
	}
}
